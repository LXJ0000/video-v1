package config

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	// 测试初始化配置
	err := Init()
	if err != nil {
		t.Errorf("配置初始化失败: %v", err)
	}

	// 验证配置值
	if GlobalConfig.MongoDB.URI == "" {
		t.Error("MongoDB URI 不能为空")
	}

	if GlobalConfig.MongoDB.Database == "" {
		t.Error("MongoDB Database 不能为空")
	}

	// 验证上传目录是否创建
	if _, err := os.Stat(GlobalConfig.Storage.UploadDir); os.IsNotExist(err) {
		t.Error("上传目录未创建")
	}

	// 清理测试数据
	os.RemoveAll(GlobalConfig.Storage.UploadDir)
}
