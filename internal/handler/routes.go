package handler

import (
	"video-platform/config"

	"github.com/gin-gonic/gin"
)

// InitRoutes 初始化路由
func InitRoutes(r *gin.Engine) {
	// 添加静态文件服务
	r.Static("/uploads", config.GlobalConfig.Storage.UploadDir)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})
	// API v1 分组
	v1 := r.Group("/api/v1")
	{
		// 视频相关路由
		videoHandler := NewVideoHandler()
		videos := v1.Group("/videos")
		{
			videos.POST("", videoHandler.Upload)                        // 上传视频
			videos.GET("", videoHandler.GetList)                        // 获取视频列表
			videos.GET("/:id", videoHandler.GetByID)                    // 获取视频详情
			videos.PUT("/:id", videoHandler.Update)                     // 更新视频信息
			videos.DELETE("/:id", videoHandler.Delete)                  // 删除视频
			videos.GET("/:id/stream", videoHandler.Stream)              // 视频流式播放
			videos.POST("/batch", videoHandler.BatchOperation)          // 批量操作
			videos.POST("/:id/thumbnail", videoHandler.UpdateThumbnail) // 更新缩略图
			videos.GET("/:id/stats", videoHandler.GetStats)             // 获取统计信息
		}
	}
}
