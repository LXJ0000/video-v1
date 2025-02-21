package handler

import (
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
)

type VideoHandler struct {
	videoService service.VideoService
}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{
		videoService: service.NewVideoService(),
	}
}

// Upload 上传视频
func (h *VideoHandler) Upload(c *gin.Context) {
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

// GetList 获取视频列表
func (h *VideoHandler) GetList(c *gin.Context) {
	var query model.VideoQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的查询参数")
		return
	}

	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	list, err := h.videoService.GetList(c.Request.Context(), query)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, list)
}

// GetByID 获取视频详情
func (h *VideoHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	video, err := h.videoService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "视频不存在")
		return
	}

	response.Success(c, video)
}

// Update 更新视频信息
func (h *VideoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var video model.Video
	if err := c.ShouldBindJSON(&video); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求数据")
		return
	}

	if err := h.videoService.Update(c.Request.Context(), id, video); err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Delete 删除视频
func (h *VideoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.videoService.Delete(c.Request.Context(), id); err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, nil)
}

// Stream 视频流式播放
func (h *VideoHandler) Stream(c *gin.Context) {
	id := c.Param("id")
	video, err := h.videoService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusNotFound, "视频不存在")
		return
	}

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
	var req model.BatchOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求数据")
		return
	}

	result, err := h.videoService.BatchOperation(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, result)
}

// UpdateThumbnail 更新视频缩略图
func (h *VideoHandler) UpdateThumbnail(c *gin.Context) {
	id := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, http.StatusBadRequest, "请选择要上传的图片文件")
		return
	}

	thumbnailURL, err := h.videoService.UpdateThumbnail(c.Request.Context(), id, file)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, gin.H{"thumbnailUrl": thumbnailURL})
}

// GetStats 获取视频统计信息
func (h *VideoHandler) GetStats(c *gin.Context) {
	id := c.Param("id")
	stats, err := h.videoService.GetStats(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, stats)
}
