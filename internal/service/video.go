package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
	"video-platform/config"
	"video-platform/internal/model"
	"video-platform/pkg/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// VideoService 视频服务接口
type VideoService interface {
	Upload(ctx context.Context, videoFile *multipart.FileHeader, coverFile *multipart.FileHeader, info model.Video) (*model.Video, error)
	GetList(ctx context.Context, query model.VideoQuery) (*model.VideoList, error)
	GetByID(ctx context.Context, id string) (*model.Video, error)
	Update(ctx context.Context, id string, video model.Video) error
	Delete(ctx context.Context, id string) error
	BatchOperation(ctx context.Context, req model.BatchOperationRequest) (*model.BatchOperationResult, error)
	UpdateThumbnail(ctx context.Context, id string, file *multipart.FileHeader) (string, error)
	GetStats(ctx context.Context, id string) (*model.VideoStats, error)
	IncrementStats(ctx context.Context, id string, field string) error
	GetVideoList(ctx context.Context, currentUserID string, opts ListOptions) ([]model.Video, int64, error)
}

type videoService struct {
	collection string
}

// NewVideoService 创建视频服务实例
func NewVideoService() VideoService {
	return &videoService{
		collection: "videos",
	}
}

// Upload 上传视频
func (s *videoService) Upload(ctx context.Context, videoFile *multipart.FileHeader, coverFile *multipart.FileHeader, info model.Video) (*model.Video, error) {
	// 验证视频文件格式
	videoExt := filepath.Ext(videoFile.Filename)
	if !isValidVideoFormat(videoExt) {
		return nil, errors.New("不支持的视频格式")
	}

	// 如果有封面图，验证格式和大小
	var coverURL string
	if coverFile != nil {
		if coverFile.Size > 2*1024*1024 { // 2MB
			return nil, errors.New("封面图大小不能超过2MB")
		}
		coverExt := strings.ToLower(filepath.Ext(coverFile.Filename))
		if !isValidImageFormat(coverExt) {
			return nil, errors.New("不支持的图片格式")
		}

		// 处理源文件
		var coverSrc io.Reader
		if testPath, ok := ctx.Value("testCoverPath").(string); ok && testPath != "" {
			// 测试模式：直接打开测试文件
			testFile, err := os.Open(testPath)
			if err != nil {
				return nil, err
			}
			defer testFile.Close()
			coverSrc = testFile
		} else {
			// 正常模式：打开上传的文件
			srcFile, err := coverFile.Open()
			if err != nil {
				return nil, err
			}
			defer srcFile.Close()
			coverSrc = srcFile
		}

		// 保存封面图
		coverFileName := fmt.Sprintf("cover_%s%s", primitive.NewObjectID().Hex(), coverExt)
		coverPath := filepath.Join(config.GlobalConfig.Storage.UploadDir, coverFileName)
		coverDst, err := os.Create(coverPath)
		if err != nil {
			return nil, err
		}
		defer coverDst.Close()

		if _, err := io.Copy(coverDst, coverSrc); err != nil {
			return nil, err
		}
		coverURL = fmt.Sprintf("/uploads/%s", coverFileName)
	}

	// 设置默认状态为私有
	if info.Status == "" {
		info.Status = model.VideoStatusPrivate
	} else if !model.IsValidVideoStatus(info.Status) {
		return nil, errors.New("无效的视频状态")
	}

	// 生成唯一文件名
	fileName := primitive.NewObjectID().Hex() + videoExt
	filePath := filepath.Join(config.GlobalConfig.Storage.UploadDir, fileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// 处理源文件
	var src io.Reader
	if testPath, ok := ctx.Value("testFilePath").(string); ok && testPath != "" {
		// 测试模式：直接打开测试文件
		testFile, err := os.Open(testPath)
		if err != nil {
			return nil, err
		}
		defer testFile.Close()
		src = testFile
	} else {
		// 正常模式：打开上传的文件
		srcFile, err := videoFile.Open()
		if err != nil {
			return nil, err
		}
		defer srcFile.Close()
		src = srcFile
	}

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		os.Remove(filePath) // 清理失败的文件
		return nil, err
	}

	// 创建视频记录
	video := model.Video{
		ID:          primitive.NewObjectID(),
		Title:       info.Title,
		Description: info.Description,
		FileName:    fileName,
		FileSize:    videoFile.Size,
		Format:      videoExt[1:],
		Status:      info.Status,
		Duration:    info.Duration,
		CoverURL:    coverURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 保存到数据库
	_, err = database.GetCollection(s.collection).InsertOne(ctx, video)
	if err != nil {
		os.Remove(filePath) // 清理文件
		return nil, err
	}

	return &video, nil
}

// GetList 获取视频列表
func (s *videoService) GetList(ctx context.Context, query model.VideoQuery) (*model.VideoList, error) {
	// 构建查询条件
	filter := bson.M{}

	// 关键词搜索
	if query.Keyword != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": query.Keyword, "$options": "i"}},
			{"description": bson.M{"$regex": query.Keyword, "$options": "i"}},
		}
	}

	// 状态筛选
	if query.Status != "" {
		filter["status"] = query.Status
	}

	// 日期范围
	if query.StartDate != "" || query.EndDate != "" {
		dateFilter := bson.M{}
		if query.StartDate != "" {
			startDate, err := time.Parse("2006-01-02", query.StartDate)
			if err == nil {
				dateFilter["$gte"] = startDate
			}
		}
		if query.EndDate != "" {
			endDate, err := time.Parse("2006-01-02", query.EndDate)
			if err == nil {
				// 将结束日期设置为当天的最后一刻
				endDate = endDate.Add(24*time.Hour - time.Second)
				dateFilter["$lte"] = endDate
			}
		}
		if len(dateFilter) > 0 {
			filter["created_at"] = dateFilter
		}
	}

	// 标签筛选
	if query.Tags != "" {
		tags := strings.Split(query.Tags, ",")
		filter["tags"] = bson.M{"$all": tags}
	}

	// 设置排序
	sortField := "created_at"
	sortOrder := -1 // 默认按创建时间降序

	if query.SortBy != "" {
		sortField = query.SortBy
		if query.SortOrder == "asc" {
			sortOrder = 1
		}
	}

	// 设置分页
	skip := (query.Page - 1) * query.PageSize
	limit := query.PageSize

	// 查询总数
	total, err := database.GetCollection(s.collection).CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 查询数据
	cursor, err := database.GetCollection(s.collection).Find(ctx,
		filter,
		options.Find().
			SetSort(bson.D{{Key: sortField, Value: sortOrder}}).
			SetSkip(int64(skip)).
			SetLimit(int64(limit)),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// 解析结果
	videos := make([]model.Video, 0)
	if err = cursor.All(ctx, &videos); err != nil {
		return nil, err
	}

	return &model.VideoList{
		Total: total,
		Items: videos,
	}, nil
}

// GetByID 根据ID获取视频
func (s *videoService) GetByID(ctx context.Context, id string) (*model.Video, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("无效的视频ID")
	}

	var video model.Video
	err = database.GetCollection(s.collection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&video)
	if err != nil {
		return nil, err
	}

	return &video, nil
}

// Update 更新视频信息
func (s *videoService) Update(ctx context.Context, id string, video model.Video) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("无效的视频ID")
	}

	// 构建更新字段
	updateFields := bson.M{
		"updated_at": time.Now(),
	}

	// 只更新非空字段
	if video.Title != "" {
		updateFields["title"] = video.Title
	}
	if video.Description != "" {
		updateFields["description"] = video.Description
	}
	if video.Status != "" {
		updateFields["status"] = video.Status
	}
	if len(video.Tags) > 0 {
		updateFields["tags"] = video.Tags
	}

	update := bson.M{
		"$set": updateFields,
	}

	result, err := database.GetCollection(s.collection).UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("视频不存在")
	}

	return nil
}

// Delete 删除视频
func (s *videoService) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("无效的视频ID")
	}

	// 先获取视频信息
	var video model.Video
	err = database.GetCollection(s.collection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&video)
	if err != nil {
		return err
	}

	// 删除数据库记录
	result, err := database.GetCollection(s.collection).DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("视频不存在")
	}

	// 删除文件
	filePath := filepath.Join(config.GlobalConfig.Storage.UploadDir, video.FileName)
	return os.Remove(filePath)
}

// BatchOperation 批量操作
func (s *videoService) BatchOperation(ctx context.Context, req model.BatchOperationRequest) (*model.BatchOperationResult, error) {
	result := &model.BatchOperationResult{
		FailedIDs: make([]string, 0),
	}

	for _, id := range req.IDs {
		var err error
		switch req.Action {
		case "delete":
			err = s.Delete(ctx, id)
		case "update_status":
			err = s.Update(ctx, id, model.Video{Status: req.Status})
		default:
			return nil, errors.New("不支持的操作类型")
		}

		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
		} else {
			result.SuccessCount++
		}
	}

	return result, nil
}

// UpdateThumbnail 更新视频缩略图
func (s *videoService) UpdateThumbnail(ctx context.Context, id string, file *multipart.FileHeader) (string, error) {
	// 验证文件大小
	if file.Size > 2*1024*1024 { // 2MB
		return "", errors.New("缩略图大小不能超过2MB")
	}

	// 验证文件格式
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isValidImageFormat(ext) {
		return "", errors.New("不支持的图片格式")
	}

	// 确保上传目录存在
	if err := os.MkdirAll(config.GlobalConfig.Storage.UploadDir, 0755); err != nil {
		return "", err
	}

	// 生成缩略图文件名
	fileName := fmt.Sprintf("thumb_%s%s", primitive.NewObjectID().Hex(), ext)
	filePath := filepath.Join(config.GlobalConfig.Storage.UploadDir, fileName)

	// 处理源文件
	var src io.Reader
	if testPath, ok := ctx.Value("testFilePath").(string); ok && testPath != "" {
		// 测试模式：直接打开测试文件
		testFile, err := os.Open(testPath)
		if err != nil {
			return "", err
		}
		defer testFile.Close()
		src = testFile
	} else {
		// 正常模式：打开上传的文件
		srcFile, err := file.Open()
		if err != nil {
			return "", err
		}
		defer srcFile.Close()
		src = srcFile
	}

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		os.Remove(filePath) // 清理失败的文件
		return "", err
	}

	// 更新数据库
	thumbnailURL := fmt.Sprintf("/uploads/%s", fileName)
	update := bson.M{
		"$set": bson.M{
			"thumbnail_url": thumbnailURL,
			"updated_at":    time.Now(),
		},
	}

	objectID, _ := primitive.ObjectIDFromHex(id)
	_, err = database.GetCollection(s.collection).UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	if err != nil {
		os.Remove(filePath) // 清理文件
		return "", err
	}

	return thumbnailURL, nil
}

// GetStats 获取视频统计信息
func (s *videoService) GetStats(ctx context.Context, id string) (*model.VideoStats, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("无效的视频ID")
	}

	var video model.Video
	err = database.GetCollection(s.collection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&video)
	if err != nil {
		return nil, err
	}

	return &video.Stats, nil
}

// IncrementStats 增加视频统计信息
func (s *videoService) IncrementStats(ctx context.Context, id string, field string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("无效的视频ID")
	}

	update := bson.M{
		"$inc": bson.M{
			"stats." + field: 1,
		},
	}

	_, err = database.GetCollection(s.collection).UpdateOne(
		ctx,
		bson.M{"_id": objectID},
		update,
	)

	return err
}

// ListOptions 视频列表查询选项
type ListOptions struct {
	Page     int
	PageSize int
	UserID   string    // 指定用户的视频
	Status   []string  // 状态过滤
	Sort     string    // 排序方式
	Keyword  string    // 搜索关键词
}

// GetVideoList 获取视频列表
func (s *videoService) GetVideoList(ctx context.Context, currentUserID string, opts ListOptions) ([]model.Video, int64, error) {
	collection := database.GetCollection("videos")
	
	// 构建查询条件
	filter := bson.M{}
	
	// 1. 处理用户ID过滤
	if opts.UserID != "" {
		// 查看指定用户的视频
		filter["user_id"] = opts.UserID
		
		if currentUserID != opts.UserID {
			// 如果不是查看自己的视频，只能看到公开视频
			filter["status"] = "public"
		}
	} else {
		// 查看所有视频
		if currentUserID != "" {
			// 已登录用户：看到所有公开视频和自己的私有/草稿视频
			filter["$or"] = []bson.M{
				{"status": "public"},
				{"$and": []bson.M{
					{"user_id": currentUserID},
					{"status": bson.M{"$in": []string{"private", "draft"}}},
				}},
			}
		} else {
			// 未登录用户：只能看到公开视频
			filter["status"] = "public"
		}
	}
	
	// 2. 处理状态过滤
	if len(opts.Status) > 0 {
		validStatus := make([]string, 0)
		for _, status := range opts.Status {
			// 对于 private 和 draft 状态，需要是视频作者
			if status == "public" || 
			   (currentUserID != "" && opts.UserID == currentUserID && 
			   (status == "private" || status == "draft")) {
				validStatus = append(validStatus, status)
			}
		}
		if len(validStatus) > 0 {
			filter["status"] = bson.M{"$in": validStatus}
		}
	}
	
	// 3. 处理关键词搜索
	if opts.Keyword != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": opts.Keyword, "$options": "i"}},
			{"description": bson.M{"$regex": opts.Keyword, "$options": "i"}},
		}
	}
	
	// 4. 获取总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	
	// 5. 处理排序
	sortOpts := bson.D{}
	if opts.Sort != "" {
		if strings.HasPrefix(opts.Sort, "-") {
			sortOpts = append(sortOpts, bson.E{Key: opts.Sort[1:], Value: -1})
		} else {
			sortOpts = append(sortOpts, bson.E{Key: opts.Sort, Value: 1})
		}
	} else {
		// 默认按创建时间倒序
		sortOpts = append(sortOpts, bson.E{Key: "created_at", Value: -1})
	}
	
	// 6. 查询数据
	findOptions := options.Find().
		SetSort(sortOpts).
		SetSkip(int64((opts.Page - 1) * opts.PageSize)).
		SetLimit(int64(opts.PageSize))
	
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	
	var videos []model.Video
	if err = cursor.All(ctx, &videos); err != nil {
		return nil, 0, err
	}
	
	return videos, total, nil
}

// 其他辅助函数
func isValidVideoFormat(ext string) bool {
	validFormats := map[string]bool{
		".mp4": true,
		".mov": true,
		".avi": true,
		".wmv": true,
		".flv": true,
		".mkv": true,
	}
	return validFormats[ext]
}

func isValidImageFormat(ext string) bool {
	validFormats := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	return validFormats[ext]
}

func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
