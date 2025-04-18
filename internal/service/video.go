package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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
	"go.mongodb.org/mongo-driver/mongo"
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
	if !isValidVideoFormat(strings.ToLower(videoExt)) {
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
		UserID:      info.UserID,
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

	// 使用事务确保数据一致性
	session, err := database.GetClient().StartSession()
	if err != nil {
		return fmt.Errorf("无法启动数据库会话: %w", err)
	}
	defer session.EndSession(ctx)

	// 在事务中执行所有删除操作
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. 删除相关收藏记录
		_, err := database.GetCollection("favorites").DeleteMany(
			sessCtx,
			bson.M{"video_id": id},
		)
		if err != nil {
			return nil, fmt.Errorf("删除收藏记录失败: %w", err)
		}

		// 2. 删除相关观看历史
		_, err = database.GetCollection("watch_history").DeleteMany(
			sessCtx,
			bson.M{"video_id": id},
		)
		if err != nil {
			return nil, fmt.Errorf("删除观看历史失败: %w", err)
		}

		// 3. 删除相关评论
		_, err = database.GetCollection("comments").DeleteMany(
			sessCtx,
			bson.M{"video_id": id},
		)
		if err != nil {
			return nil, fmt.Errorf("删除评论失败: %w", err)
		}

		// 4. 删除相关标记和注释
		_, err = database.GetCollection("marks").DeleteMany(
			sessCtx,
			bson.M{"video_id": id},
		)
		if err != nil {
			return nil, fmt.Errorf("删除标记失败: %w", err)
		}

		_, err = database.GetCollection("annotations").DeleteMany(
			sessCtx,
			bson.M{"video_id": id},
		)
		if err != nil {
			return nil, fmt.Errorf("删除注释失败: %w", err)
		}

		// 5. 最后删除视频记录本身
		result, err := database.GetCollection(s.collection).DeleteOne(
			sessCtx,
			bson.M{"_id": objectID},
		)
		if err != nil {
			return nil, fmt.Errorf("删除视频记录失败: %w", err)
		}

		if result.DeletedCount == 0 {
			return nil, errors.New("视频不存在")
		}

		return nil, nil
	})

	if err != nil {
		return err
	}

	// 删除视频文件(事务外执行文件系统操作)
	filePath := filepath.Join(config.GlobalConfig.Storage.UploadDir, video.FileName)
	if err := os.Remove(filePath); err != nil {
		// 文件删除失败，但数据库记录已删除，记录错误但继续执行
		log.Printf("WARNING: 视频文件删除失败(%s): %v", filePath, err)
	}

	// 删除缩略图文件(如果存在)
	if video.CoverURL != "" {
		coverFileName := filepath.Base(video.CoverURL)
		coverPath := filepath.Join(config.GlobalConfig.Storage.UploadDir, coverFileName)
		if err := os.Remove(coverPath); err != nil {
			// 缩略图删除失败，只记录日志不影响主流程
			log.Printf("WARNING: 视频缩略图删除失败(%s): %v", coverPath, err)
		}
	}

	return nil
}

// BatchOperation 批量操作视频
func (s *videoService) BatchOperation(ctx context.Context, req model.BatchOperationRequest) (*model.BatchOperationResult, error) {
	result := &model.BatchOperationResult{
		SuccessCount: 0,
		FailedCount:  0,
		FailedIDs:    []string{},
	}

	for _, id := range req.IDs {
		// 检查视频是否存在且属于当前用户
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			continue
		}

		var video model.Video
		err = database.GetCollection(s.collection).FindOne(ctx, bson.M{"_id": objectID}).Decode(&video)
		if err != nil {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			continue
		}

		// 权限检查
		if video.UserID != req.UserID {
			result.FailedCount++
			result.FailedIDs = append(result.FailedIDs, id)
			continue
		}

		// 执行操作
		switch req.Action {
		case "delete":
			if err := s.Delete(ctx, id); err != nil {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, id)
			} else {
				result.SuccessCount++
			}
		case "update_status":
			if !model.IsValidVideoStatus(req.Status) {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, id)
				continue
			}

			update := bson.M{
				"$set": bson.M{
					"status":     req.Status,
					"updated_at": time.Now(),
				},
			}

			_, err := database.GetCollection(s.collection).UpdateOne(
				ctx,
				bson.M{"_id": objectID},
				update,
			)

			if err != nil {
				result.FailedCount++
				result.FailedIDs = append(result.FailedIDs, id)
			} else {
				result.SuccessCount++
			}
		default:
			return nil, errors.New("不支持的操作类型")
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
	UserID   string   // 指定用户的视频
	Status   []string // 状态过滤
	Sort     string   // 排序方式
	Keyword  string   // 搜索关键词
}

// GetVideoList 获取视频列表
func (s *videoService) GetVideoList(ctx context.Context, currentUserID string, opts ListOptions) ([]model.Video, int64, error) {
	// 1. 构建基础过滤条件
	filter := bson.M{}

	// 2. 处理用户ID过滤
	if opts.UserID != "" {
		filter["user_id"] = opts.UserID
	}

	// 3. 处理状态过滤
	if len(opts.Status) > 0 {
		if currentUserID != "" && (opts.UserID == "" || opts.UserID == currentUserID) {
			// 已登录用户查看自己的视频：按指定状态过滤
			filter["status"] = bson.M{"$in": opts.Status}
		} else {
			// 查看其他用户的视频：只能看到公开视频
			filter["status"] = model.VideoStatusPublic
		}
	} else {
		// 未指定状态
		if currentUserID != "" && (opts.UserID == "" || opts.UserID == currentUserID) {
			// 已登录用户查看自己的视频：可以看到所有状态
		} else {
			// 查看其他用户的视频：只能看到公开视频
			filter["status"] = model.VideoStatusPublic
		}
	}

	// 4. 处理关键词搜索
	if opts.Keyword != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": opts.Keyword, "$options": "i"}},
			{"description": bson.M{"$regex": opts.Keyword, "$options": "i"}},
		}
	}

	// 5. 获取总数
	collection := database.GetCollection(s.collection)
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// 6. 处理排序
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

	// 7. 查询数据
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
