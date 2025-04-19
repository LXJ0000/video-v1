package service

import (
	"context"
	"encoding/base64"
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
	"video-platform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req *model.LoginRequest) (*model.User, string, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	UpdateProfile(ctx context.Context, id string, profile *model.UserProfile) error
	GetUserProfile(ctx context.Context, id string) (*model.UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, id string, req *model.UpdateProfileRequest, avatar *multipart.FileHeader) (*model.UserProfileResponse, error)
	GetWatchHistory(ctx context.Context, id string, page, size int) (*model.WatchHistoryResponse, error)
	GetFavorites(ctx context.Context, id string, page, size int) (*model.FavoriteResponse, error)
	AddToFavorites(ctx context.Context, userID, videoID string) error
	RemoveFromFavorites(ctx context.Context, userID, videoID string) error
	RecordWatchHistory(ctx context.Context, userID, videoID string) error
	CheckFavoriteStatus(ctx context.Context, userID, videoID string) (bool, error)
}

type userService struct {
	collection string
}

// NewUserService 创建用户服务实例
func NewUserService() UserService {
	return &userService{
		collection: "users",
	}
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	// 检查用户名是否已存在
	collection := database.GetCollection(s.collection)
	var existUser model.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&existUser)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 创建新用户
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Password:  hashedPassword,
		Email:     req.Email,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, string, error) {
	collection := database.GetCollection(s.collection)
	var user model.User
	err := collection.FindOne(ctx, bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		return nil, "", errors.New("用户不存在")
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, "", errors.New("密码错误")
	}

	// 生成 JWT token
	token, err := utils.GenerateToken(user.ID.Hex(), user.Username)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// 修改GetByID方法
func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	collection := database.GetCollection(s.collection)
	var user model.User

	// 将字符串ID转换为ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("无效的ID格式: %w", err)
	}

	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("未找到ID为%s的用户", id)
		}
		return nil, err
	}

	return &user, nil
}

func (s *userService) UpdateProfile(ctx context.Context, id string, profile *model.UserProfile) error {
	// Implementation needed
	return nil
}

// GetUserProfile 获取用户详细信息
func (s *userService) GetUserProfile(ctx context.Context, id string) (*model.UserProfileResponse, error) {
	// 获取用户基本信息
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取用户统计信息
	stats, err := s.getUserStats(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取用户资料
	profile, err := s.getUserProfileData(ctx, id)
	if err != nil {
		// 如果没有资料，创建默认资料
		profile = &model.UserProfile{
			Nickname:    "", // 默认昵称
			Avatar:      "", // 默认头像
			Description: "", // 默认简介
			UpdatedAt:   time.Now(),
		}
	}

	return &model.UserProfileResponse{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		Nickname:  profile.Nickname,
		Email:     user.Email,
		Avatar:    profile.Avatar,
		Bio:       profile.Description,
		CreatedAt: user.CreatedAt,
		Stats:     *stats,
	}, nil
}

// getUserStats 获取用户统计信息
func (s *userService) getUserStats(ctx context.Context, userID string) (*model.UserStats, error) {
	// 使用集合操作
	videosCollection := database.GetCollection("videos")

	// 获取上传视频数量
	videosCount, err := videosCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}

	// 获取获得的总点赞数
	// 假设videos集合中有stats字段，包含likes
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id":        nil,
			"totalLikes": bson.M{"$sum": "$stats.likes"},
		}},
	}

	cursor, err := videosCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	totalLikes := int64(0)
	if len(result) > 0 && result[0]["totalLikes"] != nil {
		totalLikes = result[0]["totalLikes"].(int64)
	}

	// 获取总观看时长（分钟）
	// 假设watch_history集合记录了观看时长
	watchHistoryCollection := database.GetCollection("watch_history")
	pipeline = []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id":            nil,
			"totalWatchTime": bson.M{"$sum": "$progress"},
		}},
	}

	cursor, err = watchHistoryCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	result = []bson.M{}
	if err = cursor.All(ctx, &result); err != nil {
		return nil, err
	}

	totalWatchTime := int64(0)
	if len(result) > 0 && result[0]["totalWatchTime"] != nil {
		// 将观看秒数转换为分钟
		totalWatchTime = int64(result[0]["totalWatchTime"].(float64) / 60)
	}

	return &model.UserStats{
		UploadedVideos: videosCount,
		TotalLikes:     totalLikes,
		TotalWatchTime: totalWatchTime,
	}, nil
}

// getUserProfileData 获取用户资料数据
func (s *userService) getUserProfileData(ctx context.Context, userID string) (*model.UserProfile, error) {
	collection := database.GetCollection("user_profiles")

	var profile model.UserProfile
	err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

// UpdateUserProfile 更新用户资料
func (s *userService) UpdateUserProfile(ctx context.Context, id string, req *model.UpdateProfileRequest, avatar *multipart.FileHeader) (*model.UserProfileResponse, error) {
	// 获取用户基本信息
	user, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 更新用户基本信息（如有变更）
	updateFields := bson.M{"updated_at": time.Now()}

	if req.Username != "" && req.Username != user.Username {
		// 检查用户名是否已存在
		var existUser model.User
		err := database.GetCollection(s.collection).FindOne(ctx, bson.M{"username": req.Username}).Decode(&existUser)
		if err == nil {
			return nil, errors.New("用户名已存在")
		}

		// 更新用户名
		updateFields["username"] = req.Username
	}

	if req.Email != "" && req.Email != user.Email {
		// 检查邮箱是否已存在
		var existUser model.User
		err := database.GetCollection(s.collection).FindOne(ctx, bson.M{"email": req.Email}).Decode(&existUser)
		if err == nil {
			return nil, errors.New("邮箱已存在")
		}

		// 更新邮箱
		updateFields["email"] = req.Email
	}

	// 如果有字段需要更新
	if len(updateFields) > 1 {
		_, err = database.GetCollection(s.collection).UpdateOne(
			ctx,
			bson.M{"_id": user.ID},
			bson.M{"$set": updateFields},
		)
		if err != nil {
			return nil, err
		}
	}

	// 获取或创建用户资料
	var profile model.UserProfile
	profileCollection := database.GetCollection("user_profiles")
	err = profileCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&profile)
	if err != nil {
		// 如果用户资料不存在，创建一个新的
		profile = model.UserProfile{
			Nickname:    "", // 默认昵称
			Avatar:      "", // 默认头像
			Description: "",
			UpdatedAt:   time.Now(),
		}
	}

	// 更新用户资料
	profileUpdateFields := bson.M{"updated_at": time.Now()}

	if req.Bio != "" {
		profileUpdateFields["description"] = req.Bio
	}

	if req.Nickname != "" {
		profileUpdateFields["nickname"] = req.Nickname
	}

	// 处理头像 - 支持Base64格式
	var avatarURL string
	if req.Avatar != "" && strings.HasPrefix(req.Avatar, "data:image/") {
		// 从Base64数据中提取图片格式和数据
		dataURI := req.Avatar
		commaIndex := strings.Index(dataURI, ",")
		if commaIndex != -1 {
			// 解析MIME类型
			mimeType := ""
			if strings.HasPrefix(dataURI, "data:") && strings.Contains(dataURI[:commaIndex], ";base64") {
				mimeType = strings.TrimPrefix(dataURI[:strings.Index(dataURI, ";")], "data:")
			}

			// 根据MIME类型确定文件扩展名
			var ext string
			switch mimeType {
			case "image/jpeg", "image/jpg":
				ext = ".jpg"
			case "image/png":
				ext = ".png"
			case "image/gif":
				ext = ".gif"
			default:
				return nil, errors.New("不支持的图片格式，只支持jpg、jpeg、png、gif")
			}

			// 解码Base64数据
			base64Data := dataURI[commaIndex+1:]
			imgData, err := base64.StdEncoding.DecodeString(base64Data)
			if err != nil {
				return nil, errors.New("解码头像图片失败: " + err.Error())
			}

			// 验证文件大小
			if len(imgData) > 10*1024*1024 { // 10MB
				return nil, errors.New("头像大小不能超过10MB")
			}

			// 生成唯一文件名
			avatarFileName := fmt.Sprintf("avatar_%s%s", id, ext)
			avatarPath := filepath.Join(config.GlobalConfig.Storage.UploadDir, avatarFileName)

			// 保存文件
			err = os.WriteFile(avatarPath, imgData, 0644)
			if err != nil {
				return nil, errors.New("保存头像失败: " + err.Error())
			}

			avatarURL = fmt.Sprintf("/uploads/%s", avatarFileName)
			profileUpdateFields["avatar"] = avatarURL
		}
	} else if avatar != nil { // 保留对multipart.FileHeader的处理以兼容旧代码
		// 验证文件大小
		if avatar.Size > 10*1024*1024 { // 2MB
			return nil, errors.New("头像大小不能超过10MB")
		}

		// 验证文件类型
		ext := strings.ToLower(filepath.Ext(avatar.Filename))
		if !s.isValidAvatarFormat(ext) {
			return nil, errors.New("不支持的图片格式，只支持jpg、jpeg、png、gif")
		}

		// 生成唯一文件名
		avatarFileName := fmt.Sprintf("avatar_%s%s", id, ext)
		avatarPath := filepath.Join(config.GlobalConfig.Storage.UploadDir, avatarFileName)

		// 保存文件
		err = s.saveAvatarFile(avatar, avatarPath)
		if err != nil {
			return nil, err
		}

		avatarURL = fmt.Sprintf("/uploads/%s", avatarFileName)
		profileUpdateFields["avatar"] = avatarURL
	}

	// 更新用户资料
	_, err = database.GetCollection("user_profiles").UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": profileUpdateFields},
		options.Update().SetUpsert(true), // 如果不存在则插入
	)
	if err != nil {
		return nil, err
	}

	// 获取更新后的用户资料
	return s.GetUserProfile(ctx, id)
}

// GetWatchHistory 获取用户观看历史
func (s *userService) GetWatchHistory(ctx context.Context, id string, page, size int) (*model.WatchHistoryResponse, error) {
	collection := database.GetCollection("watch_history")

	// 设置默认分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 12
	}

	// 查询条件
	filter := bson.M{"user_id": id}

	// 获取总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 获取分页数据
	skip := int64((page - 1) * size)
	limit := int64(size)

	options := options.Find().
		SetSort(bson.D{{Key: "watched_at", Value: -1}}). // 按观看时间倒序
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var history []model.WatchHistory
	if err = cursor.All(ctx, &history); err != nil {
		return nil, err
	}

	return &model.WatchHistoryResponse{
		History: history,
		Total:   total,
		Page:    page,
		Size:    size,
	}, nil
}

// GetFavorites 获取用户收藏列表
func (s *userService) GetFavorites(ctx context.Context, id string, page, size int) (*model.FavoriteResponse, error) {
	collection := database.GetCollection("favorites")

	// 设置默认分页参数
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 12
	}

	// 查询条件
	filter := bson.M{"user_id": id}

	// 获取总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, err
	}

	// 获取分页数据
	skip := int64((page - 1) * size)
	limit := int64(size)

	options := options.Find().
		SetSort(bson.D{{Key: "added_at", Value: -1}}). // 按添加时间倒序
		SetSkip(skip).
		SetLimit(limit)

	cursor, err := collection.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var favorites []model.Favorite
	if err = cursor.All(ctx, &favorites); err != nil {
		return nil, err
	}

	return &model.FavoriteResponse{
		Favorites: favorites,
		Total:     total,
		Page:      page,
		Size:      size,
	}, nil
}

// AddToFavorites 添加到收藏
func (s *userService) AddToFavorites(ctx context.Context, userID, videoID string) error {
	// 获取视频信息
	videoCollection := database.GetCollection("videos")

	var video model.Video
	objectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}
	if err := videoCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&video); err != nil {
		return errors.New("视频不存在")
	}

	// 检查是否已收藏
	collection := database.GetCollection("favorites")
	count, err := collection.CountDocuments(ctx, bson.M{
		"user_id":  userID,
		"video_id": videoID,
	})
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("已经收藏过该视频")
	}

	// 创建会话
	session, err := database.GetClient().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// 在事务中执行添加收藏和更新视频统计
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. 添加到收藏表
		favorite := model.Favorite{
			ID:            primitive.NewObjectID(),
			UserID:        userID,
			VideoID:       videoID,
			VideoTitle:    video.Title,
			CoverURL:      video.CoverURL,
			AddedAt:       time.Now(),
			VideoDuration: video.Duration,
		}

		_, err := collection.InsertOne(sessCtx, favorite)
		if err != nil {
			return nil, err
		}

		// 2. 更新视频的likes统计
		update := bson.M{
			"$inc": bson.M{
				"stats.likes": 1,
			},
		}

		_, err = videoCollection.UpdateOne(
			sessCtx,
			bson.M{"_id": objectID},
			update,
		)

		return nil, err
	})

	return err
}

// RemoveFromFavorites 从收藏中移除
func (s *userService) RemoveFromFavorites(ctx context.Context, userID, videoID string) error {
	// 创建会话
	session, err := database.GetClient().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// 准备集合
	collection := database.GetCollection("favorites")
	videoCollection := database.GetCollection("videos")

	// 解析视频ID
	objectID, err := primitive.ObjectIDFromHex(videoID)
	if err != nil {
		return fmt.Errorf("无效的ID格式: %w", err)
	}

	// 在事务中执行移除收藏和更新视频统计
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. 从收藏表中删除
		result, err := collection.DeleteOne(sessCtx, bson.M{
			"user_id":  userID,
			"video_id": videoID,
		})
		if err != nil {
			return nil, err
		}

		if result.DeletedCount == 0 {
			return nil, errors.New("收藏不存在")
		}

		// 2. 更新视频的likes统计（减1，但确保不会小于0）
		update := bson.M{
			"$inc": bson.M{
				"stats.likes": -1,
			},
		}

		// 条件更新，确保likes不会小于0
		filter := bson.M{
			"_id":         objectID,
			"stats.likes": bson.M{"$gt": 0},
		}

		_, err = videoCollection.UpdateOne(
			sessCtx,
			filter,
			update,
		)

		return nil, err
	})

	return err
}

// isValidAvatarFormat 验证头像格式
func (s *userService) isValidAvatarFormat(ext string) bool {
	validFormats := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	return validFormats[ext]
}

// saveAvatarFile 保存头像文件
func (s *userService) saveAvatarFile(file *multipart.FileHeader, dst string) error {
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

// RecordWatchHistory 记录用户观看历史
func (s *userService) RecordWatchHistory(ctx context.Context, userID, videoID string) error {
	// 获取视频信息
	videoCollection := database.GetCollection("videos")

	var video model.Video
	err := videoCollection.FindOne(ctx, bson.M{"_id": videoID}).Decode(&video)
	if err != nil {
		return errors.New("视频不存在")
	}

	// 检查是否已有观看记录
	historyCollection := database.GetCollection("watch_history")

	// 查找条件
	filter := bson.M{
		"user_id":  userID,
		"video_id": videoID,
	}

	// 更新数据
	update := bson.M{
		"$set": bson.M{
			"user_id":        userID,
			"video_id":       videoID,
			"video_title":    video.Title,
			"cover_url":      video.CoverURL,
			"watched_at":     time.Now(),
			"video_duration": video.Duration,
			// 简单实现，不记录具体进度
			"progress": video.Duration, // 假设看完了
		},
	}

	// 更新选项
	opts := options.Update().SetUpsert(true)

	// 执行更新
	_, err = historyCollection.UpdateOne(ctx, filter, update, opts)
	return err
}

// CheckFavoriteStatus 检查视频是否被用户收藏
func (s *userService) CheckFavoriteStatus(ctx context.Context, userID, videoID string) (bool, error) {
	// 如果没有userID或videoID则返回未收藏
	if userID == "" || videoID == "" {
		return false, nil
	}

	collection := database.GetCollection("favorites")
	count, err := collection.CountDocuments(ctx, bson.M{
		"user_id":  userID,
		"video_id": videoID,
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
