package main

import (
	"log"
	"video-platform/config"
	"video-platform/internal/handler"
	"video-platform/pkg/database"

	"github.com/gin-gonic/gin"
)

// @title 视频管理平台 API
// @version 1.0
// @description 视频管理平台后端服务API文档
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 初始化数据库连接
	if err := database.InitMongoDB(); err != nil {
		log.Fatalf("MongoDB连接失败: %v", err)
	}
	defer database.CloseMongoDB()

	// 创建Gin引擎
	r := gin.Default()

	// 配置中间件
	r.Use(gin.Recovery())
	r.Use(handler.CORSMiddleware())
	r.Use(handler.LoggerMiddleware())

	// 初始化路由
	handler.InitRoutes(r)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
