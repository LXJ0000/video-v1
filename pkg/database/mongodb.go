package database

import (
	"context"
	"time"
	"video-platform/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	db     *mongo.Database
)

// InitMongoDB 初始化MongoDB连接
func InitMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 创建MongoDB客户端
	clientOptions := options.Client().ApplyURI(config.GlobalConfig.MongoDB.URI)
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// 测试连接
	err = c.Ping(ctx, nil)
	if err != nil {
		return err
	}

	client = c
	db = client.Database(config.GlobalConfig.MongoDB.Database)
	return nil
}

// GetCollection 获取集合
func GetCollection(name string) *mongo.Collection {
	return db.Collection(name)
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
