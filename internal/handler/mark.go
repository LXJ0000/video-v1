package handler

import (
	"net/http"
	"video-platform/internal/model"
	"video-platform/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"fmt"
)

var markService = service.NewMarkService()

// AddMark 添加标记
func AddMark(c *gin.Context) {
	userID := c.Param("userId") // 假设用户ID通过URL传递
	var mark model.Mark
	if err := c.ShouldBindJSON(&mark); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的请求"})
		return
	}
	if err := markService.AddMark(c.Request.Context(), userID, &mark); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加标记失败"})
		return
	}
	c.JSON(http.StatusOK, mark)
}

// GetMarks 获取标记列表
func GetMarks(c *gin.Context) {
	userID := c.Param("userId") // 假设用户ID通过URL传递
	videoID := c.Param("id")
	marks, err := markService.GetMarks(c.Request.Context(), userID, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取标记失败"})
		return
	}
	c.JSON(http.StatusOK, marks)
}

// AddAnnotation 添加注释
func AddAnnotation(c *gin.Context) {
	var annotation model.Annotation
	if err := c.ShouldBindJSON(&annotation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的请求"})
		return
	}
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))
	annotation.MarkID = markID
	if err := markService.AddAnnotation(c.Request.Context(), &annotation); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加注释失败"})
		return
	}
	c.JSON(http.StatusOK, annotation)
}

// GetAnnotations 获取注释
func GetAnnotations(c *gin.Context) {
	markID, _ := primitive.ObjectIDFromHex(c.Param("markId"))
	annotations, err := markService.GetAnnotations(c.Request.Context(), markID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取注释失败"})
		return
	}
	c.JSON(http.StatusOK, annotations)
}

// AddNote 添加笔记
func AddNote(c *gin.Context) {
	var note model.Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "无效的请求"})
		return
	}
	if err := markService.AddNote(c.Request.Context(), &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "添加笔记失败"})
		return
	}
	c.JSON(http.StatusOK, note)
}

// GetNotes 获取笔记列表
func GetNotes(c *gin.Context) {
	videoID := c.Param("id")
	notes, err := markService.GetNotes(c.Request.Context(), videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取笔记失败"})
		return
	}
	c.JSON(http.StatusOK, notes)
}

// ExportMarks 导出标记、注释和笔记
func ExportMarks(c *gin.Context) {
	videoID := c.Param("id")
	userID := c.Param("userId")

	// 获取标记
	marks, err := markService.GetMarks(c.Request.Context(), userID, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "获取标记失败"})
		return
	}

	// 获取注释和笔记
	var annotations []model.Annotation
	var notes []model.Note
	for _, mark := range marks {
		annotations, _ = markService.GetAnnotations(c.Request.Context(), mark.ID)
		notes, _ = markService.GetNotes(c.Request.Context(), videoID)
	}

	// 创建导出文件
	file, err := os.Create(fmt.Sprintf("%s_export.txt", videoID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "创建导出文件失败"})
		return
	}
	defer file.Close()

	// 写入标记、注释和笔记
	for _, mark := range marks {
		file.WriteString(fmt.Sprintf("标记: %s, 时间戳: %f, 内容: %s\n", mark.VideoID, mark.Timestamp, mark.Content))
		for _, annotation := range annotations {
			file.WriteString(fmt.Sprintf("  注释: %s\n", annotation.Content))
		}
	}
	for _, note := range notes {
		file.WriteString(fmt.Sprintf("笔记: %s, 时间戳: %f, 内容: %s\n", note.VideoID, note.Timestamp, note.Content))
	}

	c.File(fmt.Sprintf("%s_export.txt", videoID))
}
