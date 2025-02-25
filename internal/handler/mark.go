package handler

import (
	"net/http"
	"video-platform/internal/model"
	"video-platform/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var markService = service.NewMarkService()

// AddMark 添加标记
func AddMark(c *gin.Context) {
	userID := c.Param("userId")
	var mark model.Mark
	if err := c.ShouldBindJSON(&mark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "无效的请求",
			"data": nil,
		})
		return
	}
	if err := markService.AddMark(c.Request.Context(), userID, &mark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "添加标记失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": mark,
	})
}

// GetMarks 获取标记列表
func GetMarks(c *gin.Context) {
	userID := c.Param("userId")
	videoID := c.Param("id")
	marks, err := markService.GetMarks(c.Request.Context(), userID, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "获取标记失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": marks,
	})
}

// AddAnnotation 添加注释
func AddAnnotation(c *gin.Context) {
	var annotation model.Annotation
	if err := c.ShouldBindJSON(&annotation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "无效的请求",
			"data": nil,
		})
		return
	}
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))
	annotation.MarkID = markID
	if err := markService.AddAnnotation(c.Request.Context(), &annotation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "添加注释失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": annotation,
	})
}

// GetAnnotations 获取注释
func GetAnnotations(c *gin.Context) {
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))
	annotations, err := markService.GetAnnotations(c.Request.Context(), markID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "获取注释失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": annotations,
	})
}

// AddNote 添加笔记
func AddNote(c *gin.Context) {
	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "无效的请求",
			"data": nil,
		})
		return
	}
	if err := markService.AddNote(c.Request.Context(), &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "添加笔记失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": note,
	})
}

// GetNotes 获取笔记列表
func GetNotes(c *gin.Context) {
	videoID := c.Param("id")
	notes, err := markService.GetNotes(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "获取笔记失败",
			"data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": notes,
	})
}

// ExportMarks 导出标记、注释和笔记
func ExportMarks(c *gin.Context) {
	videoID := c.Param("id")
	userID := c.Param("userId")

	// 获取标记
	marks, err := markService.GetMarks(c.Request.Context(), userID, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "获取标记失败",
			"data": nil,
		})
		return
	}

	// 获取注释和笔记
	var annotations []model.Annotation
	var notes []model.Note
	for _, mark := range marks {
		annotations, _ = markService.GetAnnotations(c.Request.Context(), mark.ID)
		notes, _ = markService.GetNotes(c.Request.Context(), videoID)
	}

	// 构建导出数据
	exportData := gin.H{
		"marks":       marks,
		"annotations": annotations,
		"notes":       notes,
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": exportData,
	})
}

// UpdateMark 更新标记
func UpdateMark(c *gin.Context) {
	userID := c.Param("userId")
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))

	var mark model.Mark
	if err := c.ShouldBindJSON(&mark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "无效的请求",
			"data": nil,
		})
		return
	}

	if err := markService.UpdateMark(c.Request.Context(), userID, markID, &mark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "更新标记失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": mark,
	})
}

// DeleteMark 删除标记
func DeleteMark(c *gin.Context) {
	userID := c.Param("userId")
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))

	if err := markService.DeleteMark(c.Request.Context(), userID, markID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "删除标记失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": nil,
	})
}

// UpdateAnnotation 更新注释
func UpdateAnnotation(c *gin.Context) {
	userID := c.Param("userId")
	annotationID, _ := primitive.ObjectIDFromHex(c.Param("annotationId"))

	var annotation model.Annotation
	if err := c.ShouldBindJSON(&annotation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "无效的请求",
			"data": nil,
		})
		return
	}

	if err := markService.UpdateAnnotation(c.Request.Context(), userID, annotationID, &annotation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "更新注释失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": annotation,
	})
}

// DeleteAnnotation 删除注释
func DeleteAnnotation(c *gin.Context) {
	userID := c.Param("userId")
	annotationID, _ := primitive.ObjectIDFromHex(c.Param("annotationId"))

	if err := markService.DeleteAnnotation(c.Request.Context(), userID, annotationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "删除注释失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": nil,
	})
}

// UpdateNote 更新笔记
func UpdateNote(c *gin.Context) {
	userID := c.Param("userId")
	noteID, _ := primitive.ObjectIDFromHex(c.Param("noteId"))

	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 1,
			"msg":  "无效的请求",
			"data": nil,
		})
		return
	}

	if err := markService.UpdateNote(c.Request.Context(), userID, noteID, &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "更新笔记失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": note,
	})
}

// DeleteNote 删除笔记
func DeleteNote(c *gin.Context) {
	userID := c.Param("userId")
	noteID, _ := primitive.ObjectIDFromHex(c.Param("noteId"))

	if err := markService.DeleteNote(c.Request.Context(), userID, noteID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 1,
			"msg":  "删除笔记失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": nil,
	})
}
