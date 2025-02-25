package database

import (
	"context"
	"errors"
	"strings"
	"time"
	"video-platform/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client   *mongo.Client
	database string
)

// InitMongoDB 初始化MongoDB连接
func InitMongoDB(ctx context.Context, cfg config.MongoDBConfig, isTest bool) error {
	clientOptions := options.Client().ApplyURI(cfg.URI)
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 根据环境选择数据库
	if isTest {
		database = cfg.TestDatabase
	} else {
		database = cfg.Database
	}

	return client.Ping(ctx, nil)
}

// GetCollection 获取集合
func GetCollection(name string) *mongo.Collection {
	return client.Database(database).Collection(name)
}

// GetClient 获取MongoDB客户端
func GetClient() *mongo.Client {
	return client
}

// CloseMongoDB 关闭MongoDB连接
func CloseMongoDB() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		client.Disconnect(ctx)
	}
}

// CleanupTestData 清理测试数据
func CleanupTestData(ctx context.Context) error {
	if database == config.GlobalConfig.MongoDB.TestDatabase {
		// 改用删除集合而不是删除数据库
		collections := []string{"users", "marks", "annotations", "notes", "videos"}
		for _, name := range collections {
			if err := client.Database(database).Collection(name).Drop(ctx); err != nil {
				// 忽略集合不存在的错误
				if !strings.Contains(err.Error(), "ns not found") {
					return err
				}
			}
		}
		return nil
	}
	return errors.New("can only cleanup test database")
}
