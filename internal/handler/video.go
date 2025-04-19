package handler

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"video-platform/config"
	"video-platform/internal/model"
	"video-platform/internal/service"
	"video-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type VideoHandler struct {
	videoService service.VideoService
}

func NewVideoHandler(videoService service.VideoService) *VideoHandler {
	if videoService == nil {
		videoService = service.NewVideoService()
	}
	return &VideoHandler{
		videoService: videoService,
	}
}

// Upload 上传视频
func (h *VideoHandler) Upload(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	videoFile, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "请选择要上传的视频文件")
		return
	}

	// 获取封面图文件（可选）
	var coverFile *multipart.FileHeader
	if f, err := c.FormFile("cover"); err == nil {
		coverFile = f
	}

	// 获取视频时长
	duration, err := strconv.ParseFloat(c.PostForm("duration"), 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的视频时长")
		return
	}

	info := model.Video{
		UserID:      userID.(string), // 设置用户ID
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Status:      c.PostForm("status"),
		Duration:    duration,
	}

	// 处理标签
	if tags := c.PostForm("tags"); tags != "" {
		info.Tags = strings.Split(tags, ",")
	}

	video, err := h.videoService.Upload(c.Request.Context(), videoFile, coverFile, info)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, video)
}

// GetByID 获取视频详情
func (h *VideoHandler) GetByID(c *gin.Context) {
	// 获取视频ID
	videoID := c.Param("videoId")
	if videoID == "" {
		response.Fail(c, http.StatusBadRequest, "视频ID不能为空")
		slog.Error("[GetByID] 视频ID为空")
		return
	}

	// 获取视频详情
	video, err := h.videoService.GetByID(c.Request.Context(), videoID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		slog.Error("[GetByID] 获取视频详情失败", "error", err)
		return
	}

	// 如果是私有视频，需要验证权限
	if video.Status != "public" {
		userID, exists := c.Get("userId")
		if !exists || userID.(string) != video.UserID {
			response.Fail(c, http.StatusForbidden, "无权查看该视频")
			slog.Error("[GetByID] 无权查看私有视频", "videoId", videoID)
			return
		}
	}

	// 构建响应数据
	result := gin.H{
		"video": video,
	}

	// 检查用户是否已收藏视频
	userID, exists := c.Get("userId")
	if exists && userID.(string) != "" {
		// 获取用户服务实例
		userService := service.NewUserService()
		isFavorite, err := userService.CheckFavoriteStatus(c.Request.Context(), userID.(string), videoID)
		if err == nil {
			// 只在无错误时添加是否收藏信息
			result["isFavorite"] = isFavorite
		}
	}

	response.Success(c, result)
}

// Update 更新视频信息
func (h *VideoHandler) Update(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	videoId := c.Param("videoId")

	// 检查视频是否存在且属于当前用户
	video, err := h.videoService.GetByID(c.Request.Context(), videoId)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "视频不存在")
		return
	}

	// 权限检查
	if video.UserID != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权操作此视频")
		return
	}

	var updateData model.Video
	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 设置用户ID，确保不会被修改
	updateData.UserID = userID.(string)

	if err := h.videoService.Update(c.Request.Context(), videoId, updateData); err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除视频
func (h *VideoHandler) Delete(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	videoId := c.Param("videoId")

	// 检查视频是否存在且属于当前用户
	video, err := h.videoService.GetByID(c.Request.Context(), videoId)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "视频不存在")
		return
	}

	// 权限检查
	if video.UserID != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权操作此视频")
		return
	}

	if err := h.videoService.Delete(c.Request.Context(), videoId); err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Stream 视频流式播放
func (h *VideoHandler) Stream(c *gin.Context) {
	videoId := c.Param("videoId")
	video, err := h.videoService.GetByID(c.Request.Context(), videoId)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "视频不存在")
		return
	}

	// 权限检查：非公开视频只有作者可以观看
	// if video.Status != model.VideoStatusPublic {
	// 	// 获取当前用户ID
	// 	userID, exists := c.Get("userId")
	// 	if !exists || userID.(string) != video.UserID {
	// 		response.Fail(c, http.StatusForbidden, "无权观看此视频")
	// 		return
	// 	}
	// }

	// 增加观看次数
	go h.videoService.IncrementStats(context.Background(), videoId, "views")

	filePath := filepath.Join(config.GlobalConfig.Storage.UploadDir, video.FileName)
	file, err := os.Open(filePath)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "无法打开视频文件")
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "无法获取文件信息")
		return
	}

	// 设置响应头
	c.Header("Content-Type", "video/"+video.Format)
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.Header("Accept-Ranges", "bytes")

	// 处理范围请求
	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		ranges, err := parseRange(rangeHeader, fileInfo.Size())
		if err != nil {
			response.Fail(c, http.StatusRequestedRangeNotSatisfiable, "无效的范围请求")
			return
		}

		if len(ranges) > 0 {
			start, end := ranges[0][0], ranges[0][1]
			c.Status(http.StatusPartialContent)
			c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
			c.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
			file.Seek(start, 0)
			io.CopyN(c.Writer, file, end-start+1)
			return
		}
	}

	// 非范围请求，直接返回完整文件
	io.Copy(c.Writer, file)
}

// parseRange 解析Range头部
func parseRange(rangeHeader string, size int64) ([][2]int64, error) {
	if !strings.HasPrefix(rangeHeader, "bytes=") {
		return nil, fmt.Errorf("invalid range format")
	}
	rangeHeader = strings.TrimPrefix(rangeHeader, "bytes=")
	ranges := make([][2]int64, 0)
	for _, r := range strings.Split(rangeHeader, ",") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		parts := strings.Split(r, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format")
		}
		start, end := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		var startByte, endByte int64
		if start == "" {
			// -N
			endByte, _ = strconv.ParseInt(end, 10, 64)
			startByte = size - endByte
			endByte = size - 1
		} else {
			startByte, _ = strconv.ParseInt(start, 10, 64)
			if end == "" {
				// N-
				endByte = size - 1
			} else {
				// N-M
				endByte, _ = strconv.ParseInt(end, 10, 64)
			}
		}
		if startByte < 0 || endByte >= size || startByte > endByte {
			return nil, fmt.Errorf("invalid range")
		}
		ranges = append(ranges, [2]int64{startByte, endByte})
	}
	return ranges, nil
}

// BatchOperation 批量操作视频
func (h *VideoHandler) BatchOperation(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	var req model.BatchOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求数据")
		return
	}

	// 设置用户ID，用于权限检查
	req.UserID = userID.(string)

	result, err := h.videoService.BatchOperation(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateThumbnail 更新视频缩略图
func (h *VideoHandler) UpdateThumbnail(c *gin.Context) {
	// 获取当前用户ID
	userID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	videoId := c.Param("videoId")

	// 检查视频是否存在且属于当前用户
	video, err := h.videoService.GetByID(c.Request.Context(), videoId)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "视频不存在")
		return
	}

	// 权限检查
	if video.UserID != userID.(string) {
		response.Fail(c, http.StatusForbidden, "无权操作此视频")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "请选择要上传的图片文件")
		return
	}

	thumbnailURL, err := h.videoService.UpdateThumbnail(c.Request.Context(), videoId, file)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"thumbnailUrl": thumbnailURL})
}

// GetStats 获取视频统计信息
func (h *VideoHandler) GetStats(c *gin.Context) {
	videoId := c.Param("videoId")
	stats, err := h.videoService.GetStats(c.Request.Context(), videoId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, stats)
}

// GetVideoList 获取视频列表
func (h *VideoHandler) GetVideoList(c *gin.Context) {
	// 1. 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize > 50 {
		pageSize = 50
	}

	// 获取当前用户ID（如果已登录）
	currentUserID, _ := c.Get("userId")
	userIDStr, _ := currentUserID.(string)

	// 2. 构建查询选项
	opts := service.ListOptions{
		Page:     page,
		PageSize: pageSize,
		UserID:   userIDStr,          // 可选：指定用户的视频
		Keyword:  c.Query("keyword"), // 可选：搜索关键词
		Sort:     c.Query("sortBy"),  // 修正排序字段名
	}

	// 处理排序方向
	if order := c.Query("sortOrder"); order == "desc" {
		opts.Sort = "-" + opts.Sort // 添加降序前缀
	}

	// 处理状态过滤
	if status := c.Query("status"); status != "" {
		opts.Status = []string{status} // 只使用单个状态，不再分割
	}

	// 3. 调用 service 获取视频列表
	videos, total, err := h.videoService.GetVideoList(c.Request.Context(), userIDStr, opts)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取视频列表失败")
		return
	}

	// 4. 返回结果
	response.Success(c, gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"items":    videos,
	})
}

// GetPublicVideoList 获取公开视频列表（无需认证）
func (h *VideoHandler) GetPublicVideoList(c *gin.Context) {
	// 1. 获取查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if pageSize > 50 {
		pageSize = 50
	}

	// 2. 构建查询选项
	opts := service.ListOptions{
		Page:     page,
		PageSize: pageSize,
		Status:   []string{"public"}, // 只获取公开视频
		UserID:   c.Query("userId"),  // 可选：指定用户的公开视频
		Keyword:  c.Query("keyword"), // 可选：搜索关键词
		Sort:     c.Query("sort"),    // 可选：排序方式
	}

	// 3. 调用 service 获取视频列表（传入空的 currentUserID）
	videos, total, err := h.videoService.GetVideoList(c.Request.Context(), "", opts)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取视频列表失败")
		return
	}

	// 4. 返回结果
	response.Success(c, gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"items":    videos,
	})
}
