package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User 用户模型
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`    // 用户名
	Password  string             `bson:"password" json:"-"`           // 密码（加密存储）
	Email     string             `bson:"email" json:"email"`          // 邮箱
	Status    int                `bson:"status" json:"status"`        // 状态 1:正常 2:禁用
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"` // 创建时间
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"` // 更新时间
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=1,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

// UserProfile 用户资料
type UserProfile struct {
	Nickname    string    `json:"nickname" bson:"nickname"`
	Avatar      string    `json:"avatar" bson:"avatar"`
	Description string    `json:"description" bson:"description"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// UserProfileResponse 用户详细信息响应
type UserProfileResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Bio       string    `json:"bio"`
	CreatedAt time.Time `json:"createdAt"`
	Stats     UserStats `json:"stats"`
}

// UserStats 用户统计信息
type UserStats struct {
	UploadedVideos int64 `json:"uploadedVideos" bson:"uploaded_videos"`
	TotalWatchTime int64 `json:"totalWatchTime" bson:"total_watch_time"` // 单位：分钟
	TotalLikes     int64 `json:"totalLikes" bson:"total_likes"`
}

// UpdateProfileRequest 更新用户资料请求
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"omitempty,min=1,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Bio      string `json:"bio" binding:"omitempty,max=200"`
	Nickname string `json:"nickname" binding:"omitempty,min=1,max=32"`
	// Avatar通过multipart/form-data上传，不在JSON中
}

// WatchHistory 观看历史记录
type WatchHistory struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"user_id" json:"userId"`
	VideoID       string             `bson:"video_id" json:"videoId"`
	VideoTitle    string             `bson:"video_title" json:"videoTitle"`
	CoverURL      string             `bson:"cover_url" json:"coverUrl"`
	WatchedAt     time.Time          `bson:"watched_at" json:"watchedAt"`
	Progress      float64            `bson:"progress" json:"progress"`            // 观看进度(秒)
	VideoDuration float64            `bson:"video_duration" json:"videoDuration"` // 视频总时长(秒)
}

// WatchHistoryResponse 观看历史响应
type WatchHistoryResponse struct {
	History []WatchHistory `json:"history"`
	Total   int64          `json:"total"`
	Page    int            `json:"page"`
	Size    int            `json:"size"`
}

// Favorite 用户收藏
type Favorite struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID        string             `bson:"user_id" json:"userId"`
	VideoID       string             `bson:"video_id" json:"videoId"`
	VideoTitle    string             `bson:"video_title" json:"videoTitle"`
	CoverURL      string             `bson:"cover_url" json:"coverUrl"`
	AddedAt       time.Time          `bson:"added_at" json:"addedAt"`
	VideoDuration float64            `bson:"video_duration" json:"videoDuration"`
}

// FavoriteResponse 收藏列表响应
type FavoriteResponse struct {
	Favorites []Favorite `json:"favorites"`
	Total     int64      `json:"total"`
	Page      int        `json:"page"`
	Size      int        `json:"size"`
}
