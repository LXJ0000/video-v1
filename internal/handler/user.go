package handler

import (
	"log/slog"
	"net/http"
	"video-platform/internal/model"
	"video-platform/internal/service"
	"video-platform/pkg/response"

	"mime/multipart"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	if userService == nil {
		userService = service.NewUserService()
	}
	return &UserHandler{
		userService: userService,
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求参数")
		slog.Error("[Register] 无效的请求参数", "error", err)
		return
	}

	user, err := h.userService.Register(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[Register] 注册失败", "error", err)
		return
	}

	response.Success(c, user)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求参数")
		slog.Error("[Login] 无效的请求参数", "error", err)
		return
	}

	user, token, err := h.userService.Login(c.Request.Context(), &req)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, err.Error())
		slog.Error("[Login] 登录失败", "error", err)
		return
	}

	response.Success(c, gin.H{
		"user":  user,
		"token": token,
	})
}

// GetUserProfile 获取用户详细信息
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("userId")

	// 如果是获取当前用户信息
	if userID == "me" {
		currentUserID, exists := c.Get("userId")
		if !exists {
			response.Fail(c, http.StatusUnauthorized, "用户未登录")
			slog.Error("[GetUserProfile] 用户未登录")
			return
		}
		userID = currentUserID.(string)
	}

	// 获取用户资料
	profile, err := h.userService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[GetUserProfile] 获取用户资料失败", "error", err)
		return
	}

	response.Success(c, profile)
}

// UpdateUserProfile 更新用户信息
func (h *UserHandler) UpdateUserProfile(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("userId")

	// 如果是更新当前用户信息
	if userID == "me" {
		currentUserID, exists := c.Get("userId")
		if !exists {
			response.Fail(c, http.StatusUnauthorized, "用户未登录")
			slog.Error("[UpdateUserProfile] 用户未登录")
			return
		}
		userID = currentUserID.(string)
	} else {
		// 只能更新自己的信息
		currentUserID, exists := c.Get("userId")
		if !exists || currentUserID.(string) != userID {
			response.Fail(c, http.StatusForbidden, "无权更新其他用户的信息")
			slog.Error("[UpdateUserProfile] 无权更新", "requestedId", userID)
			return
		}
	}

	// 处理表单数据
	var req model.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求参数")
		slog.Error("[UpdateUserProfile] 无效的请求参数", "error", err)
		return
	}

	// 处理头像文件
	var avatar *multipart.FileHeader
	file, err := c.FormFile("avatar")
	if err == nil {
		avatar = file
	}

	// 更新用户资料
	profile, err := h.userService.UpdateUserProfile(c.Request.Context(), userID, &req, avatar)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[UpdateUserProfile] 更新用户资料失败", "error", err)
		return
	}

	response.Success(c, profile)
}

// GetWatchHistory 获取用户观看历史
func (h *UserHandler) GetWatchHistory(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("userId")

	// 验证权限
	currentUserID, exists := c.Get("userId")
	if !exists || currentUserID.(string) != userID {
		response.Fail(c, http.StatusForbidden, "无权查看其他用户的观看历史")
		slog.Error("[GetWatchHistory] 无权查看", "requestedId", userID)
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "12"))

	// 获取观看历史
	history, err := h.userService.GetWatchHistory(c.Request.Context(), userID, page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[GetWatchHistory] 获取观看历史失败", "error", err)
		return
	}

	response.Success(c, history)
}

// GetFavorites 获取用户收藏列表
func (h *UserHandler) GetFavorites(c *gin.Context) {
	// 获取用户ID
	userID := c.Param("userId")

	// 验证权限
	currentUserID, exists := c.Get("userId")
	if !exists || currentUserID.(string) != userID {
		response.Fail(c, http.StatusForbidden, "无权查看其他用户的收藏列表")
		slog.Error("[GetFavorites] 无权查看", "requestedId", userID)
		return
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "12"))

	// 获取收藏列表
	favorites, err := h.userService.GetFavorites(c.Request.Context(), userID, page, size)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[GetFavorites] 获取收藏列表失败", "error", err)
		return
	}

	response.Success(c, favorites)
}

// AddToFavorites 添加收藏
func (h *UserHandler) AddToFavorites(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "用户未登录")
		slog.Error("[AddToFavorites] 用户未登录")
		return
	}

	// 获取视频ID
	videoID := c.Param("videoId")
	if videoID == "" {
		response.Fail(c, http.StatusBadRequest, "视频ID不能为空")
		slog.Error("[AddToFavorites] 视频ID为空")
		return
	}

	// 添加到收藏
	err := h.userService.AddToFavorites(c.Request.Context(), userID.(string), videoID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[AddToFavorites] 添加收藏失败", "error", err)
		return
	}

	response.Success(c, gin.H{"message": "添加收藏成功"})
}

// RemoveFromFavorites 取消收藏
func (h *UserHandler) RemoveFromFavorites(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "用户未登录")
		slog.Error("[RemoveFromFavorites] 用户未登录")
		return
	}

	// 获取视频ID
	videoID := c.Param("videoId")
	if videoID == "" {
		response.Fail(c, http.StatusBadRequest, "视频ID不能为空")
		slog.Error("[RemoveFromFavorites] 视频ID为空")
		return
	}

	// 从收藏中移除
	err := h.userService.RemoveFromFavorites(c.Request.Context(), userID.(string), videoID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[RemoveFromFavorites] 取消收藏失败", "error", err)
		return
	}

	response.Success(c, gin.H{"message": "取消收藏成功"})
}

// RecordWatchHistory 记录观看历史
func (h *UserHandler) RecordWatchHistory(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "用户未登录")
		slog.Error("[RecordWatchHistory] 用户未登录")
		return
	}

	// 获取视频ID
	videoID := c.Param("videoId")
	if videoID == "" {
		response.Fail(c, http.StatusBadRequest, "视频ID不能为空")
		slog.Error("[RecordWatchHistory] 视频ID为空")
		return
	}

	// 记录观看历史
	err := h.userService.RecordWatchHistory(c.Request.Context(), userID.(string), videoID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[RecordWatchHistory] 记录观看历史失败", "error", err)
		return
	}

	response.Success(c, gin.H{"message": "记录观看历史成功"})
}
