package handler

import (
	"video-platform/config"
	"video-platform/internal/middleware"
	"video-platform/internal/service"

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
		// 创建服务实例
		userService := service.NewUserService()
		markService := service.NewMarkService()
		videoService := service.NewVideoService()
		// codeService := service.NewCodeSerivce(nil)

		// 创建 handler 实例
		userHandler := NewUserHandler(userService)
		markHandler := NewMarkHandler(markService)
		videoHandler := NewVideoHandler(videoService)

		// 用户相关路由（无需认证）
		users := v1.Group("/users")
		{
			users.POST("/register", userHandler.Register) // 用户注册
			users.POST("/login", userHandler.Login)       // 用户登录
			users.GET("/:userId/profile", userHandler.GetUserProfile)
			users.PUT("/:userId/profile", middleware.Auth(), userHandler.UpdateUserProfile)
			users.GET("/:userId/watch-history", middleware.Auth(), userHandler.GetWatchHistory)
			users.GET("/:userId/favorites", middleware.Auth(), userHandler.GetFavorites)
		}

		// 公开接口（无需认证）
		videos := v1.Group("/videos")
		{
			videos.GET("/public", videoHandler.GetPublicVideoList)                                  // 获取公开视频列表
			videos.GET("/:videoId/stream", videoHandler.Stream)                                     // 视频流式播放
			videos.GET("/:videoId", middleware.SetUserId(), videoHandler.GetByID)                   // 获取视频详情
			videos.POST("/:videoId/favorite", middleware.Auth(), userHandler.AddToFavorites)        // 添加收藏
			videos.DELETE("/:videoId/favorite", middleware.Auth(), userHandler.RemoveFromFavorites) // 取消收藏
			videos.POST("/:videoId/watch", middleware.Auth(), userHandler.RecordWatchHistory)       // 记录观看历史
		}

		// 需要认证的路由
		auth := v1.Group("")
		auth.Use(middleware.Auth())
		{
			// 视频相关路由
			authVideos := auth.Group("/videos")
			{
				authVideos.GET("", videoHandler.GetVideoList) // 获取视频列表
				authVideos.POST("", videoHandler.Upload)      // 上传视频
				// authVideos.GET("/:videoId", videoHandler.GetByID)   // 获取视频详情
				authVideos.PUT("/:videoId", videoHandler.Update)    // 更新视频信息
				authVideos.DELETE("/:videoId", videoHandler.Delete) // 删除视频
				// authVideos.GET("/:videoId/stream", videoHandler.Stream)              // 视频流式播放
				authVideos.POST("/batch", videoHandler.BatchOperation)               // 批量操作
				authVideos.POST("/:videoId/thumbnail", videoHandler.UpdateThumbnail) // 更新缩略图
				authVideos.GET("/:videoId/stats", videoHandler.GetStats)             // 获取统计信息
			}

			// 标记相关路由
			marks := auth.Group("/marks")
			{
				marks.POST("", markHandler.AddMark)                                      // 添加标记
				marks.GET("", markHandler.GetMarks)                                      // 获取标记列表
				marks.PUT("/:markId", markHandler.UpdateMark)                            // 更新标记
				marks.DELETE("/:markId", markHandler.DeleteMark)                         // 删除标记
				marks.POST("/:markId/annotations", markHandler.AddAnnotation)            // 添加注释
				marks.GET("/:markId/annotations", markHandler.GetAnnotations)            // 获取注释列表
				marks.PUT("/annotations/:annotationId", markHandler.UpdateAnnotation)    // 更新注释
				marks.DELETE("/annotations/:annotationId", markHandler.DeleteAnnotation) // 删除注释
			}

			// 笔记相关路由
			notes := auth.Group("/notes")
			{
				notes.POST("", markHandler.AddNote)              // 添加笔记
				notes.GET("", markHandler.GetNotes)              // 获取笔记列表
				notes.PUT("/:noteId", markHandler.UpdateNote)    // 更新笔记
				notes.DELETE("/:noteId", markHandler.DeleteNote) // 删除笔记
			}

			// 导出相关路由
			authVideos.GET("/export", markHandler.ExportMarks) // 导出标记、注释和笔记
		}
	}
}
