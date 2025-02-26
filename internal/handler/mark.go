package handler

import (
	"net/http"
	"video-platform/internal/model"
	"video-platform/internal/service"
	"video-platform/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarkHandler struct {
	markService service.MarkService
}

func NewMarkHandler(markService service.MarkService) *MarkHandler {
	if markService == nil {
		markService = service.NewMarkService()
	}
	return &MarkHandler{
		markService: markService,
	}
}

// AddMark 添加标记
func (h *MarkHandler) AddMark(c *gin.Context) {
	// 从上下文获取用户ID
	userID, _ := c.Get("userId")

	var mark model.Mark
	if err := c.ShouldBindJSON(&mark); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求")
		return
	}

	// 设置用户ID
	mark.UserID = userID.(string)

	if err := h.markService.AddMark(c.Request.Context(), userID.(string), &mark); err != nil {
		response.Fail(c, http.StatusInternalServerError, "添加标记失败")
		return
	}

	response.Success(c, mark)
}

// GetMarks 获取标记列表
func (h *MarkHandler) GetMarks(c *gin.Context) {
	// 从上下文获取用户ID
	userID, _ := c.Get("userId")
	videoID := c.Query("videoId") // 从查询参数获取视频ID

	marks, err := h.markService.GetMarks(c.Request.Context(), userID.(string), videoID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取标记列表失败")
		return
	}

	response.Success(c, marks)
}

// UpdateMark 更新标记
func (h *MarkHandler) UpdateMark(c *gin.Context) {
	userID, _ := c.Get("userId")
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))

	// 获取当前登录用户ID
	currentUserID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 权限检查：只能更新自己的标记
	if userID != currentUserID.(string) {
		response.Fail(c, http.StatusForbidden, "无权操作")
		return
	}

	var mark model.Mark
	if err := c.ShouldBindJSON(&mark); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求")
		return
	}

	if err := h.markService.UpdateMark(c.Request.Context(), userID.(string), markID, &mark); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新标记失败")
		return
	}

	response.Success(c, nil)
}

// DeleteMark 删除标记
func (h *MarkHandler) DeleteMark(c *gin.Context) {
	userID, _ := c.Get("userId")
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))

	// 获取当前登录用户ID
	currentUserID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 权限检查：只能删除自己的标记
	if userID != currentUserID.(string) {
		response.Fail(c, http.StatusForbidden, "无权操作")
		return
	}

	if err := h.markService.DeleteMark(c.Request.Context(), userID.(string), markID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除标记失败")
		return
	}

	response.Success(c, nil)
}

// AddAnnotation 添加注释
func (h *MarkHandler) AddAnnotation(c *gin.Context) {
	userID, _ := c.Get("userId")
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))

	// 获取当前登录用户ID
	currentUserID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 权限检查：只能操作自己的注释
	if userID != currentUserID.(string) {
		response.Fail(c, http.StatusForbidden, "无权操作")
		return
	}

	var annotation model.Annotation
	if err := c.ShouldBindJSON(&annotation); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求")
		return
	}

	annotation.UserID = userID.(string)
	annotation.MarkID = markID

	if err := h.markService.AddAnnotation(c.Request.Context(), &annotation); err != nil {
		response.Fail(c, http.StatusInternalServerError, "添加注释失败")
		return
	}

	response.Success(c, annotation)
}

// GetAnnotations 获取注释列表
func (h *MarkHandler) GetAnnotations(c *gin.Context) {
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))

	annotations, err := h.markService.GetAnnotations(c.Request.Context(), markID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取注释列表失败")
		return
	}

	response.Success(c, annotations)
}

// UpdateAnnotation 更新注释
func (h *MarkHandler) UpdateAnnotation(c *gin.Context) {
	userID, _ := c.Get("userId")
	annotationID, _ := primitive.ObjectIDFromHex(c.Param("annotationId"))

	var annotation model.Annotation
	if err := c.ShouldBindJSON(&annotation); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求")
		return
	}

	if err := h.markService.UpdateAnnotation(c.Request.Context(), userID.(string), annotationID, &annotation); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新注释失败")
		return
	}

	response.Success(c, nil)
}

// DeleteAnnotation 删除注释
func (h *MarkHandler) DeleteAnnotation(c *gin.Context) {
	userID, _ := c.Get("userId")
	annotationID, _ := primitive.ObjectIDFromHex(c.Param("annotationId"))

	if err := h.markService.DeleteAnnotation(c.Request.Context(), userID.(string), annotationID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除注释失败")
		return
	}

	response.Success(c, nil)
}

// AddNote 添加笔记
func (h *MarkHandler) AddNote(c *gin.Context) {
	userID, _ := c.Get("userId")

	// 获取当前登录用户ID
	currentUserID, exists := c.Get("userId")
	if !exists {
		response.Fail(c, http.StatusUnauthorized, "未授权")
		return
	}

	// 权限检查：只能操作自己的笔记
	if userID != currentUserID.(string) {
		response.Fail(c, http.StatusForbidden, "无权操作")
		return
	}

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求")
		return
	}

	// 设置用户ID
	note.UserID = userID.(string)

	if err := h.markService.AddNote(c.Request.Context(), &note); err != nil {
		response.Fail(c, http.StatusInternalServerError, "添加笔记失败")
		return
	}
	response.Success(c, note)
}

// GetNotes 获取笔记列表
func (h *MarkHandler) GetNotes(c *gin.Context) {
	videoID := c.Query("videoId") // 从查询参数获取视频ID
	notes, err := h.markService.GetNotes(c.Request.Context(), videoID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取笔记失败")
		return
	}
	response.Success(c, notes)
}

// ExportMarks 导出标记、注释和笔记
func (h *MarkHandler) ExportMarks(c *gin.Context) {
	videoID := c.Query("videoId") // 从查询参数获取视频ID
	userID, _ := c.Get("userId")

	marks, err := h.markService.GetMarks(c.Request.Context(), userID.(string), videoID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取标记失败")
		return
	}

	notes, err := h.markService.GetNotes(c.Request.Context(), videoID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取笔记失败")
		return
	}

	response.Success(c, gin.H{
		"marks": marks,
		"notes": notes,
	})
}

// UpdateNote 更新笔记
func (h *MarkHandler) UpdateNote(c *gin.Context) {
	userID, _ := c.Get("userId")
	noteID, _ := primitive.ObjectIDFromHex(c.Param("noteId"))

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		response.Fail(c, http.StatusBadRequest, "无效的请求")
		return
	}

	if err := h.markService.UpdateNote(c.Request.Context(), userID.(string), noteID, &note); err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新笔记失败")
		return
	}

	response.Success(c, note)
}

// DeleteNote 删除笔记
func (h *MarkHandler) DeleteNote(c *gin.Context) {
	userID, _ := c.Get("userId")
	noteID, _ := primitive.ObjectIDFromHex(c.Param("noteId"))

	if err := h.markService.DeleteNote(c.Request.Context(), userID.(string), noteID); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除笔记失败")
		return
	}

	response.Success(c, nil)
}
