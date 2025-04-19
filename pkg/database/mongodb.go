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

// 以下变量和函数用于测试目的
var (
	GetCollection = getCollection // 导出获取集合的函数
	GetClient     = getClient     // 导出获取客户端的函数
)

// SetTestHooks 设置测试时使用的替代函数，并返回用于恢复原始函数的回调
func SetTestHooks(
	collectionFn func(string) *mongo.Collection,
	clientFn func() *mongo.Client,
) func() {
	// 保存原始函数
	origGetCollection := GetCollection
	origGetClient := GetClient

	// 替换为测试用函数
	GetCollection = collectionFn
	GetClient = clientFn

	// 返回用于恢复原始函数的回调
	return func() {
		GetCollection = origGetCollection
		GetClient = origGetClient
	}
}

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

// getCollection 获取集合的内部实现
func getCollection(name string) *mongo.Collection {
	return client.Database(database).Collection(name)
}

// getClient 获取MongoDB客户端的内部实现
func getClient() *mongo.Client {
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
