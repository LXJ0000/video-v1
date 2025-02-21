package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 视频状态常量
const (
	VideoStatusPublic  = "public"  // 公开
	VideoStatusPrivate = "private" // 私有
	VideoStatusDraft   = "draft"   // 草稿
)

// 添加状态验证函数
func IsValidVideoStatus(status string) bool {
	switch status {
	case VideoStatusPublic, VideoStatusPrivate, VideoStatusDraft:
		return true
	default:
		return false
	}
}

// Video 视频模型
type Video struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       string             `bson:"user_id" json:"userId"`             // 作者ID
	Title        string             `bson:"title" json:"title"`                // 视频标题
	Description  string             `bson:"description" json:"description"`    // 视频描述
	FileName     string             `bson:"file_name" json:"fileName"`         // 文件名
	FileSize     int64              `bson:"file_size" json:"fileSize"`         // 文件大小（字节）
	Duration     float64            `bson:"duration" json:"duration"`          // 视频时长（秒）
	Format       string             `bson:"format" json:"format"`              // 视频格式
	Status       string             `bson:"status" json:"status"`              // 视频状态
	Tags         []string           `bson:"tags" json:"tags"`                  // 视频标签
	ThumbnailURL string             `bson:"thumbnail_url" json:"thumbnailUrl"` // 缩略图URL
	CoverURL     string             `bson:"cover_url" json:"coverUrl"`         // 封面图URL
	Stats        VideoStats         `bson:"stats" json:"stats"`                // 视频统计信息
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`       // 创建时间
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`       // 更新时间
}

// VideoStats 视频统计信息
type VideoStats struct {
	Views    int64 `bson:"views" json:"views"`       // 观看次数
	Likes    int64 `bson:"likes" json:"likes"`       // 点赞数
	Comments int64 `bson:"comments" json:"comments"` // 评论数
	Shares   int64 `bson:"shares" json:"shares"`     // 分享次数
}

// VideoList 视频列表响应结构
type VideoList struct {
	Total int64   `json:"total"` // 总记录数
	Items []Video `json:"items"` // 视频列表
}

// VideoQuery 视频查询参数
type VideoQuery struct {
	Page      int    `form:"page" binding:"min=1"`                                              // 页码
	PageSize  int    `form:"pageSize" binding:"min=1,max=50"`                                   // 每页数量
	Keyword   string `form:"keyword"`                                                           // 关键词搜索（标题、描述）
	Status    string `form:"status"`                                                            // 视频状态
	StartDate string `form:"startDate"`                                                         // 开始日期
	EndDate   string `form:"endDate"`                                                           // 结束日期
	Tags      string `form:"tags"`                                                              // 标签（逗号分隔）
	SortBy    string `form:"sortBy" binding:"omitempty,oneof=created_at views likes file_size"` // 排序字段
	SortOrder string `form:"sortOrder" binding:"omitempty,oneof=asc desc"`                      // 排序方向
}

// BatchOperationRequest 批量操作请求
type BatchOperationRequest struct {
	IDs    []string `json:"ids"`              // 视频ID列表
	Action string   `json:"action"`           // 操作类型
	Status string   `json:"status,omitempty"` // 状态（可选）
}

// BatchOperationResult 批量操作结果
type BatchOperationResult struct {
	SuccessCount int      `json:"successCount"` // 成功数量
	FailedCount  int      `json:"failedCount"`  // 失败数量
	FailedIDs    []string `json:"failedIds"`    // 失败的ID列表
}

// VideoUpdateRequest 视频更新请求
type VideoUpdateRequest struct {
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Status       string   `json:"status"`
	Tags         []string `json:"tags"`
	ThumbnailURL string   `json:"thumbnail"`
}
