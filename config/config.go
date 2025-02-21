package config

import (
	"os"
)

// Config 全局配置结构体
type Config struct {
	MongoDB MongoDBConfig
	Server  ServerConfig
	Storage StorageConfig
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI      string
	Database string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int
}

// StorageConfig 存储配置
type StorageConfig struct {
	UploadDir string
	MaxSize   int64 // 最大文件大小（字节）
}

var GlobalConfig Config

// Init 初始化配置
func Init() error {
	// 在实际项目中应该从配置文件读取，这里为了简单直接硬编码
	GlobalConfig = Config{
		MongoDB: MongoDBConfig{
			URI:      "mongodb://root:9hq29bfn@test-db-mongodb.ns-bpq7yu1b.svc:27017",
			Database: "video_platform",
		},
		Server: ServerConfig{
			Port: 8080,
		},
		Storage: StorageConfig{
			UploadDir: "./uploads",
			MaxSize:   1024 * 1024 * 1024, // 1GB
		},
	}

	// 确保上传目录存在
	if err := os.MkdirAll(GlobalConfig.Storage.UploadDir, 0755); err != nil {
		return err
	}

	return nil
}
