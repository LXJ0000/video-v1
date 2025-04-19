package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config 全局配置结构体
type Config struct {
	Env     string
	MongoDB MongoDBConfig
	Server  ServerConfig
	Storage StorageConfig
	JWT     JWTConfig
	Redis   RedisConfig
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	URI          string
	Database     string
	TestDatabase string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port          string
	AllowedOrigin []string
}

// StorageConfig 存储配置
type StorageConfig struct {
	UploadDir string
	MaxSize   int64 // 1GB
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	ExpireTime int64 // 24 hour
}

// RedisConfig Redis配置
type RedisConfig struct {
	URI string
}

var GlobalConfig Config

// 从环境变量获取字符串，如果不存在则返回默认值
func getEnvString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 从环境变量获取整数，如果不存在或解析失败则返回默认值
func getEnvInt64(key string, defaultValue int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return intValue
}

// 从环境变量获取字符串数组，以逗号分隔
func getEnvStringSlice(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return strings.Split(value, ",")
}

// Init 初始化配置
func Init() error {
	// 尝试加载.env文件，忽略文件不存在的错误
	_ = godotenv.Load()

	// 从环境变量加载配置
	GlobalConfig = Config{
		Env: getEnvString("ENV", "development"),
		MongoDB: MongoDBConfig{
			URI:          getEnvString("MONGODB_URI", "mongodb://root:9hq29bfn@test-db-mongodb.ns-bpq7yu1b.svc:27017"),
			Database:     getEnvString("MONGODB_DATABASE", "video_platform"),
			TestDatabase: getEnvString("MONGODB_TEST_DATABASE", "video_platform_test"),
		},
		Server: ServerConfig{
			Port:          getEnvString("SERVER_PORT", "8080"),
			AllowedOrigin: getEnvStringSlice("SERVER_ALLOWED_ORIGIN", []string{"*"}),
		},
		Storage: StorageConfig{
			UploadDir: getEnvString("STORAGE_UPLOAD_DIR", "./uploads"),
			MaxSize:   getEnvInt64("STORAGE_MAX_SIZE", 1024*1024*1024), // 1GB
		},
		JWT: JWTConfig{
			Secret:     getEnvString("JWT_SECRET", "your-secret-key"),
			ExpireTime: getEnvInt64("JWT_EXPIRE_TIME", 24), // 24 hours
		},
		Redis: RedisConfig{
			URI: getEnvString("REDIS_URI", "redis://localhost:6379/0"),
		},
	}

	// 确保上传目录存在
	uploadDir := filepath.Join(GlobalConfig.Storage.UploadDir)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return err
	}

	return nil
}
