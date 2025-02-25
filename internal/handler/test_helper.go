package handler

import (
	"context"
	"testing"
	"time"
	"video-platform/config"
	"video-platform/pkg/database"
	"video-platform/pkg/utils"

	"github.com/gin-gonic/gin"
)

// 添加测试用户信息
const (
	TestUserID   = "test_user_id"
	TestUsername = "test_user"
)

// waitForDBOperation 等待数据库操作完成
func waitForDBOperation() {
	time.Sleep(100 * time.Millisecond)
}

func setupTestEnvironment(t *testing.T) func() {
	// 初始化配置
	if err := config.Init(); err != nil {
		t.Fatal(err)
	}

	// 连接测试数据库
	ctx := context.Background()
	if err := database.InitMongoDB(ctx, config.GlobalConfig.MongoDB, true); err != nil {
		t.Fatal(err)
	}

	// 等待连接完全建立
	waitForDBOperation()

	// 设置Gin测试模式
	gin.SetMode(gin.TestMode)

	// 返回清理函数
	return func() {
		ctx := context.Background()
		if err := database.CleanupTestData(ctx); err != nil {
			t.Error(err)
		}
		// 等待清理完成
		waitForDBOperation()
	}
}

// 生成测试用的认证token
func getTestToken(t *testing.T) string {
	token, err := utils.GenerateToken(TestUserID, TestUsername)
	if err != nil {
		t.Fatalf("生成测试token失败: %v", err)
	}
	return token
}
