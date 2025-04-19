package handler

import (
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"video-platform/internal/model"
	"video-platform/internal/service"
	"video-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 创建一个VideoService的Mock
type MockVideoService struct {
	mock.Mock
}

func (m *MockVideoService) Upload(ctx context.Context, videoFile *multipart.FileHeader, coverFile *multipart.FileHeader, info model.Video) (*model.Video, error) {
	args := m.Called(ctx, videoFile, coverFile, info)
	return args.Get(0).(*model.Video), args.Error(1)
}

func (m *MockVideoService) GetList(ctx context.Context, query model.VideoQuery) (*model.VideoList, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*model.VideoList), args.Error(1)
}

func (m *MockVideoService) GetByID(ctx context.Context, id string) (*model.Video, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Video), args.Error(1)
}

func (m *MockVideoService) Update(ctx context.Context, id string, video model.Video) error {
	args := m.Called(ctx, id, video)
	return args.Error(0)
}

func (m *MockVideoService) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockVideoService) GetVideoList(ctx context.Context, currentUserID string, opts service.ListOptions) ([]model.Video, int64, error) {
	args := m.Called(ctx, currentUserID, opts)
	return args.Get(0).([]model.Video), args.Get(1).(int64), args.Error(2)
}

func (m *MockVideoService) GetStats(ctx context.Context, id string) (*model.VideoStats, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.VideoStats), args.Error(1)
}

func (m *MockVideoService) UpdateThumbnail(ctx context.Context, id string, file *multipart.FileHeader) (string, error) {
	args := m.Called(ctx, id, file)
	return args.String(0), args.Error(1)
}

func (m *MockVideoService) BatchOperation(ctx context.Context, req model.BatchOperationRequest) (*model.BatchOperationResult, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*model.BatchOperationResult), args.Error(1)
}

func (m *MockVideoService) IncrementStats(ctx context.Context, id string, field string) error {
	args := m.Called(ctx, id, field)
	return args.Error(0)
}

// 设置测试环境
func setupVideoTest() (*gin.Context, *httptest.ResponseRecorder, *MockVideoService, *VideoHandler) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建默认的请求对象，避免Context()方法返回nil
	c.Request = httptest.NewRequest("GET", "/", nil)

	mockService := new(MockVideoService)
	handler := NewVideoHandler(mockService)
	return c, w, mockService, handler
}

// 测试获取视频详情带收藏状态
func TestGetVideoByIDWithFavoriteStatus(t *testing.T) {
	c, w, mockService, handler := setupVideoTest()

	// 模拟参数
	videoID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	c.Params = []gin.Param{{Key: "videoId", Value: videoID}}
	c.Set("userId", userID) // 模拟用户已登录

	// 模拟视频数据
	objID, _ := primitive.ObjectIDFromHex(videoID)
	mockVideo := &model.Video{
		ID:          objID,
		UserID:      primitive.NewObjectID().Hex(), // 不同于当前用户
		Title:       "Test Video",
		Description: "Test Description",
		FileName:    "test.mp4",
		Format:      "mp4",
		FileSize:    1024000,
		Duration:    120.5,
		CoverURL:    "cover.jpg",
		Status:      "public",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Tags:        []string{"test", "video"},
		Stats: model.VideoStats{
			Views:  100,
			Likes:  50,
			Shares: 10,
		},
	}

	// 模拟服务响应
	mockService.On("GetByID", mock.Anything, videoID).Return(mockVideo, nil)

	// 创建Mock用户服务
	mockUserService := new(MockUserService)
	
	// 保存原始的用户服务获取函数
	originalGetUserService := getUserService
	
	// 替换为返回mock的函数
	getUserService = func() service.UserService {
		return mockUserService
	}
	
	// 注册测试完成后的清理函数
	t.Cleanup(func() {
		getUserService = originalGetUserService
	})
	
	// 模拟CheckFavoriteStatus调用
	mockUserService.On("CheckFavoriteStatus", mock.Anything, userID, videoID).Return(true, nil)

	// 执行测试
	handler.GetByID(c)

	// 验证响应状态码
	assert.Equal(t, http.StatusOK, w.Code)

	// 解析响应体
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证响应数据包含video和isFavorite字段
	responseData, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	_, hasVideo := responseData["video"]
	assert.True(t, hasVideo)
	isFavorite, hasFavorite := responseData["isFavorite"]
	assert.True(t, hasFavorite)
	assert.Equal(t, true, isFavorite)

	// 验证调用次数
	mockService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}

// 测试获取不存在的视频
func TestGetVideoByIDNotFound(t *testing.T) {
	c, w, mockService, handler := setupVideoTest()

	// 模拟参数
	videoID := primitive.NewObjectID().Hex()
	c.Params = []gin.Param{{Key: "videoId", Value: videoID}}

	// 模拟服务返回错误
	mockService.On("GetByID", mock.Anything, videoID).Return(&model.Video{}, errors.New("视频不存在"))

	// 执行测试
	handler.GetByID(c)

	// 验证响应状态码
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 解析响应体
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, resp.Code) // 错误码不为0

	// 验证调用次数
	mockService.AssertExpectations(t)
}

// 测试无权查看私有视频
func TestGetPrivateVideoNoPermission(t *testing.T) {
	c, w, mockService, handler := setupVideoTest()

	// 模拟参数
	videoID := primitive.NewObjectID().Hex()
	userID := primitive.NewObjectID().Hex()
	videoOwnerID := primitive.NewObjectID().Hex() // 视频所有者与当前用户不同
	c.Params = []gin.Param{{Key: "videoId", Value: videoID}}
	c.Set("userId", userID) // 模拟用户已登录

	// 模拟视频数据 - 私有视频
	objID, _ := primitive.ObjectIDFromHex(videoID)
	mockVideo := &model.Video{
		ID:        objID,
		UserID:    videoOwnerID, // 不同于当前用户
		Title:     "Private Video",
		Status:    "private", // 私有视频
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 模拟服务响应
	mockService.On("GetByID", mock.Anything, videoID).Return(mockVideo, nil)

	// 执行测试
	handler.GetByID(c)

	// 验证响应状态码 - 应该是403 Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)

	// 验证调用次数
	mockService.AssertExpectations(t)
}
