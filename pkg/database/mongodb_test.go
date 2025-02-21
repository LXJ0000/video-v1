package database

import (
	"testing"
	"video-platform/config"
)

func TestMongoDBConnection(t *testing.T) {
	// 初始化配置
	if err := config.Init(); err != nil {
		t.Fatalf("配置初始化失败: %v", err)
	}

	// 测试数据库连接
	err := InitMongoDB()
	if err != nil {
		t.Fatalf("MongoDB连接失败: %v", err)
	}
	defer CloseMongoDB()

	// 测试获取集合
	collection := GetCollection("videos")
	if collection == nil {
		t.Error("获取集合失败")
	}
}
