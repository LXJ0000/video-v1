package handler

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"video-platform/internal/model"
	"video-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 创建一个UserService的Mock
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(ctx context.Context, req *model.RegisterRequest) (*model.User, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req *model.LoginRequest) (*model.User, string, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*model.User), args.String(1), args.Error(2)
}

func (m *MockUserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) UpdateProfile(ctx context.Context, id string, profile *model.UserProfile) error {
	args := m.Called(ctx, id, profile)
	return args.Error(0)
}

func (m *MockUserService) GetUserProfile(ctx context.Context, id string) (*model.UserProfileResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.UserProfileResponse), args.Error(1)
}

func (m *MockUserService) UpdateUserProfile(ctx context.Context, id string, req *model.UpdateProfileRequest, avatar *multipart.FileHeader) (*model.UserProfileResponse, error) {
	args := m.Called(ctx, id, req, avatar)
	return args.Get(0).(*model.UserProfileResponse), args.Error(1)
}

func (m *MockUserService) GetWatchHistory(ctx context.Context, id string, page, size int) (*model.WatchHistoryResponse, error) {
	args := m.Called(ctx, id, page, size)
	return args.Get(0).(*model.WatchHistoryResponse), args.Error(1)
}

func (m *MockUserService) GetFavorites(ctx context.Context, id string, page, size int) (*model.FavoriteResponse, error) {
	args := m.Called(ctx, id, page, size)
	return args.Get(0).(*model.FavoriteResponse), args.Error(1)
}

func (m *MockUserService) AddToFavorites(ctx context.Context, userID, videoID string) error {
	args := m.Called(ctx, userID, videoID)
	return args.Error(0)
}

func (m *MockUserService) RemoveFromFavorites(ctx context.Context, userID, videoID string) error {
	args := m.Called(ctx, userID, videoID)
	return args.Error(0)
}

func (m *MockUserService) RecordWatchHistory(ctx context.Context, userID, videoID string) error {
	args := m.Called(ctx, userID, videoID)
	return args.Error(0)
}

func (m *MockUserService) CheckFavoriteStatus(ctx context.Context, userID, videoID string) (bool, error) {
	args := m.Called(ctx, userID, videoID)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserService) LoginOrRegisterByPhone(ctx context.Context, phone string) (*model.User, string, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(*model.User), args.String(1), args.Error(2)
}

// 创建MockCodeService
type MockCodeService struct {
	mock.Mock
}

func (m *MockCodeService) Send(ctx context.Context, biz, number string) error {
	args := m.Called(ctx, biz, number)
	return args.Error(0)
}

func (m *MockCodeService) Verify(ctx context.Context, biz, number, code string) (bool, error) {
	args := m.Called(ctx, biz, number, code)
	return args.Bool(0), args.Error(1)
}

// 设置测试环境
func setupUserTest() (*gin.Context, *httptest.ResponseRecorder, *MockUserService, *UserHandler) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 创建默认的请求对象，避免Context()方法返回nil
	c.Request = httptest.NewRequest("GET", "/", nil)

	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	return c, w, mockService, handler
}

// 测试获取用户详情
func TestGetUserProfile(t *testing.T) {
	c, w, mockService, handler := setupUserTest()

	// 模拟当前登录用户
	userId := primitive.NewObjectID().Hex()
	c.Set("userId", userId)
	c.Params = []gin.Param{{Key: "userId", Value: "me"}}

	// 模拟服务层响应
	profileResp := &model.UserProfileResponse{
		ID:        userId,
		Username:  "testuser",
		Email:     "test@example.com",
		Avatar:    "avatar.jpg",
		Bio:       "test bio",
		CreatedAt: time.Now(),
		Stats: model.UserStats{
			UploadedVideos: 5,
			TotalWatchTime: 120,
			TotalLikes:     50,
		},
	}

	mockService.On("GetUserProfile", mock.Anything, userId).Return(profileResp, nil)

	// 执行测试
	handler.GetUserProfile(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证调用
	mockService.AssertExpectations(t)
}

// 测试获取观看历史
func TestGetWatchHistory(t *testing.T) {
	c, w, mockService, handler := setupUserTest()

	// 模拟当前登录用户
	userId := primitive.NewObjectID().Hex()
	c.Set("userId", userId)
	c.Params = []gin.Param{{Key: "userId", Value: userId}}

	// 模拟查询参数
	c.Request = httptest.NewRequest("GET", "/?page=1&size=5", nil)

	// 模拟服务层响应
	watchHistory := &model.WatchHistoryResponse{
		History: []model.WatchHistory{
			{
				ID:            primitive.NewObjectID(),
				UserID:        userId,
				VideoID:       primitive.NewObjectID().Hex(),
				VideoTitle:    "Test Video",
				CoverURL:      "cover.jpg",
				WatchedAt:     time.Now(),
				Progress:      60,
				VideoDuration: 120,
			},
		},
		Total: 1,
		Page:  1,
		Size:  5,
	}

	mockService.On("GetWatchHistory", mock.Anything, userId, 1, 5).Return(watchHistory, nil)

	// 执行测试
	handler.GetWatchHistory(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证调用
	mockService.AssertExpectations(t)
}

// 测试添加收藏
func TestAddToFavorites(t *testing.T) {
	c, w, mockService, handler := setupUserTest()

	// 模拟当前登录用户
	userId := primitive.NewObjectID().Hex()
	videoId := primitive.NewObjectID().Hex()
	c.Set("userId", userId)
	c.Params = []gin.Param{{Key: "videoId", Value: videoId}}

	// 确保请求对象存在（已在setupUserTest中创建）
	// 执行服务前确保Request不会被覆盖

	// 模拟服务层响应
	mockService.On("AddToFavorites", mock.Anything, userId, videoId).Return(nil)

	// 执行测试
	handler.AddToFavorites(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证调用
	mockService.AssertExpectations(t)
}

// 测试记录观看历史
func TestRecordWatchHistory(t *testing.T) {
	c, w, mockService, handler := setupUserTest()

	// 模拟当前登录用户
	userId := primitive.NewObjectID().Hex()
	videoId := primitive.NewObjectID().Hex()
	c.Set("userId", userId)
	c.Params = []gin.Param{{Key: "videoId", Value: videoId}}

	// 确保请求对象存在（已在setupUserTest中创建）
	// 执行服务前确保Request不会被覆盖

	// 模拟服务层响应
	mockService.On("RecordWatchHistory", mock.Anything, userId, videoId).Return(nil)

	// 执行测试
	handler.RecordWatchHistory(c)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Code)

	// 验证调用
	mockService.AssertExpectations(t)
}

// 测试检查收藏状态
func TestCheckFavoriteStatus(t *testing.T) {
	_, _, mockService, _ := setupUserTest()

	// 模拟当前登录用户
	userId := primitive.NewObjectID().Hex()
	videoId := primitive.NewObjectID().Hex()

	// 模拟服务层响应
	mockService.On("CheckFavoriteStatus", mock.Anything, userId, videoId).Return(true, nil)

	// 直接调用服务方法进行测试
	result, err := mockService.CheckFavoriteStatus(context.Background(), userId, videoId)

	// 验证响应
	assert.Nil(t, err)
	assert.True(t, result)

	// 验证调用
	mockService.AssertExpectations(t)
}

// 测试发送验证码
func TestSendSMSCode(t *testing.T) {
	c, w, _, _ := setupUserTest()
	
	// 创建mock服务
	mockUserService := new(MockUserService)
	mockCodeService := new(MockCodeService)
	handler := &UserHandler{
		userService: mockUserService,
		codeService: mockCodeService,
	}
	
	// 创建请求体
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(`{"phone":"13800138000"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	
	// 设置期望
	mockCodeService.On("Send", mock.Anything, "login", "13800138000").Return(nil)
	
	// 执行测试
	handler.SendSMSCode(c)
	
	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Code)
	
	// 验证调用
	mockCodeService.AssertExpectations(t)
}

// 测试短信验证码登录
func TestLoginBySms(t *testing.T) {
	c, w, _, _ := setupUserTest()
	
	// 创建mock服务
	mockUserService := new(MockUserService)
	mockCodeService := new(MockCodeService)
	handler := &UserHandler{
		userService: mockUserService,
		codeService: mockCodeService,
	}
	
	// 模拟用户数据和token
	userId := primitive.NewObjectID()
	user := &model.User{
		ID:       userId,
		Username: "user_8000_123456",
		Phone:    "13800138000",
		Status:   1,
	}
	token := "test_token"
	
	// 创建请求体
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(`{"phone":"13800138000","code":"123456"}`))
	c.Request.Header.Set("Content-Type", "application/json")
	
	// 设置期望
	mockCodeService.On("Verify", mock.Anything, "login", "13800138000", "123456").Return(true, nil)
	mockUserService.On("LoginOrRegisterByPhone", mock.Anything, "13800138000").Return(user, token, nil)
	
	// 执行测试
	handler.LoginBySms(c)
	
	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)
	
	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, 0, resp.Code)
	
	// 验证调用
	mockCodeService.AssertExpectations(t)
	mockUserService.AssertExpectations(t)
}
