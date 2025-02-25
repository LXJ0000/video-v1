package handler

import (
	"video-platform/config"
	"video-platform/internal/middleware"

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
		// 用户相关路由（无需认证）
		users := v1.Group("/users")
		{
			users.POST("/register", Register) // 用户注册
			users.POST("/login", Login)      // 用户登录
		}

		// 需要认证的路由
		auth := v1.Group("")
		auth.Use(middleware.Auth())
		{
			// 视频相关路由
			videoHandler := NewVideoHandler()
			videos := auth.Group("/videos")
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

				// 标记相关路由
				marks := auth.Group("/marks/:userId/:id")
				{
					marks.POST("", AddMark)                          // 添加标记
					marks.GET("", GetMarks)                          // 获取标记列表
					marks.PUT("/:markId", UpdateMark)                // 更新标记
					marks.DELETE("/:markId", DeleteMark)             // 删除标记
					marks.POST("/annotations/:markId", AddAnnotation) // 添加注释
					marks.GET("/annotations/:markId", GetAnnotations) // 获取注释
					marks.PUT("/annotations/:annotationId", UpdateAnnotation) // 更新注释
					marks.DELETE("/annotations/:annotationId", DeleteAnnotation) // 删除注释
				}

				// 笔记相关路由
				notes := auth.Group("/notes/:userId/:id")
				{
					notes.POST("", AddNote)           // 添加笔记
					notes.GET("", GetNotes)           // 获取笔记列表
					notes.PUT("/:noteId", UpdateNote) // 更新笔记
					notes.DELETE("/:noteId", DeleteNote) // 删除笔记
				}

				// 导出相关路由
				videos.GET("/export/:userId/:id", ExportMarks) // 导出标记、注释和笔记
			}
		}
	}
}
