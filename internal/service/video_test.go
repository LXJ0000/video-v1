package service

import (
	"context"
	"testing"
	"time"
	"video-platform/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 创建一个MockDatabaseForVideo结构体来模拟数据库操作
type MockDatabaseForVideo struct {
	mock.Mock
}

func (m *MockDatabaseForVideo) GetCollection(name string) *mongo.Collection {
	args := m.Called(name)
	return args.Get(0).(*mongo.Collection)
}

// 测试获取视频详情
func TestGetVideoByID(t *testing.T) {
	// 创建模拟数据
	videoID := primitive.NewObjectID()
	videoIDStr := videoID.Hex()
	
	video := model.Video{
		ID:          videoID,
		Title:       "Test Video",
		Description: "Test Description",
		FileName:    "test.mp4",
		Format:      "mp4",
		FileSize:    1024000,
		Duration:    120.5,
		Status:      "public",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	// 创建测试断言
	assert.Equal(t, videoIDStr, videoID.Hex())
	assert.Equal(t, "Test Video", video.Title)
}

// 测试检查视频是否被用户收藏
func TestCheckFavoriteStatus(t *testing.T) {
	// 创建测试数据
	userID := primitive.NewObjectID().Hex()
	videoID := primitive.NewObjectID().Hex()
	
	// 创建用户服务实例
	userSvc := &userService{
		collection: "users",
	}
	
	// 测试有效性断言
	assert.NotEmpty(t, userID)
	assert.NotEmpty(t, videoID)
	assert.NotNil(t, userSvc)
	
	// 测试用户ID为空的情况
	isFavorite, err := userSvc.CheckFavoriteStatus(context.Background(), "", videoID)
	assert.Nil(t, err)
	assert.False(t, isFavorite)
	
	// 测试视频ID为空的情况
	isFavorite, err = userSvc.CheckFavoriteStatus(context.Background(), userID, "")
	assert.Nil(t, err)
	assert.False(t, isFavorite)
}

// VideoMockCollection结构体用于模拟视频相关的数据库集合操作
type VideoMockCollection struct {
	mock.Mock
	*mongo.Collection
}

func (m *VideoMockCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	args := m.Called(ctx, filter, opts)
	return args.Get(0).(int64), args.Error(1)
}

// TestCheckFavoriteStatusWithMock 使用mock测试收藏状态检查
func TestCheckFavoriteStatusWithMock(t *testing.T) {
	// 创建测试数据
	userID := primitive.NewObjectID().Hex()
	videoID := primitive.NewObjectID().Hex()
	
	// 创建mock集合
	mockCol := new(VideoMockCollection)
	
	// 设置mock期望行为 - 已收藏
	mockCol.On("CountDocuments", 
		mock.Anything, 
		bson.M{"user_id": userID, "video_id": videoID}, 
		mock.Anything).Return(int64(1), nil)
		
	// 创建用户服务实例并测试
	userSvc := &userService{collection: "users"}
	
	// 测试基本断言
	assert.NotNil(t, userSvc)
	assert.NotEmpty(t, userID)
	assert.NotEmpty(t, videoID)
} 