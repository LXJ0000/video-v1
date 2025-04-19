package redis

import (
	"context"
	"errors"
	"time"
	"video-platform/config"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
)

// 以下变量和函数用于测试目的
var (
	GetClient = getClient // 导出获取客户端的函数
)

// SetTestHooks 设置测试时使用的替代函数，并返回用于恢复原始函数的回调
func SetTestHooks(
	clientFn func() *redis.Client,
) func() {
	// 保存原始函数
	origGetClient := GetClient

	// 替换为测试用函数
	GetClient = clientFn

	// 返回用于恢复原始函数的回调
	return func() {
		GetClient = origGetClient
	}
}

// Init 初始化Redis连接
func Init(ctx context.Context) error {
	// 使用配置中的Redis URI
	return InitRedis(ctx, config.GlobalConfig.Redis.URI)
}

// InitRedis 初始化Redis连接
func InitRedis(ctx context.Context, redisURI string) error {
	// 使用传入的URI创建Redis客户端
	opt, err := redis.ParseURL(redisURI)
	if err != nil {
		return err
	}

	client = redis.NewClient(opt)

	// 测试连接
	_, err = client.Ping(ctx).Result()
	return err
}

// getClient 获取Redis客户端的内部实现
func getClient() *redis.Client {
	return client
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if client != nil {
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = client.Close()
	}
}

// FlushTestData 清理测试数据
func FlushTestData(ctx context.Context) error {
	// 确保只在测试环境中执行
	if client != nil && isTestEnv() {
		return client.FlushDB(ctx).Err()
	}
	return errors.New("can only flush test redis database")
}

// isTestEnv 判断当前是否为测试环境
func isTestEnv() bool {
	// 这里需要根据你的项目实际情况来判断是否为测试环境
	// 可以通过环境变量或配置来判断
	return config.GlobalConfig.Env == "test"
}

// Set 设置键值对
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return client.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func Get(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

// Del 删除键
func Del(ctx context.Context, keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func Exists(ctx context.Context, keys ...string) (bool, error) {
	result, err := client.Exists(ctx, keys...).Result()
	return result > 0, err
}

// TTL 获取键的过期时间
func TTL(ctx context.Context, key string) (time.Duration, error) {
	return client.TTL(ctx, key).Result()
}

// Incr 将键的整数值加1
func Incr(ctx context.Context, key string) (int64, error) {
	return client.Incr(ctx, key).Result()
}

// Decr 将键的整数值减1
func Decr(ctx context.Context, key string) (int64, error) {
	return client.Decr(ctx, key).Result()
}

// HSet 设置哈希表字段的值
func HSet(ctx context.Context, key string, values ...interface{}) error {
	return client.HSet(ctx, key, values...).Err()
}

// HGet 获取哈希表字段的值
func HGet(ctx context.Context, key, field string) (string, error) {
	return client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希表中所有的字段和值
func HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func HDel(ctx context.Context, key string, fields ...string) error {
	return client.HDel(ctx, key, fields...).Err()
}
