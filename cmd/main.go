package main

import (
	"context"
	"log"
	"video-platform/config"
	"video-platform/internal/handler"
	"video-platform/pkg/database"
	"video-platform/pkg/redis"

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
		log.Fatal(err)
	}

	// 连接数据库（生产环境）
	ctx := context.Background()
	if err := database.InitMongoDB(ctx, config.GlobalConfig.MongoDB, false); err != nil {
		log.Fatal(err)
	}
	defer database.CloseMongoDB()
	if err := redis.InitRedis(ctx, config.GlobalConfig.Redis.URI); err != nil {
		log.Fatal(err)
	}
	defer redis.CloseRedis()

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
