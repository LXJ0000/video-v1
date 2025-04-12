package handler

import (
	"log/slog"
	"net/http"
	"video-platform/internal/model"
	"video-platform/internal/service"
	"video-platform/pkg/response"

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
