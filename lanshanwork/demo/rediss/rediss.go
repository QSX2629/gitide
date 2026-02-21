package rediss

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client //全局仅用一个
// InitRedis 初始化
func InitRedis() error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", //默认地址
		Password: "",               //地址为空
		DB:       0,                //使用第0个数据库
	})
	//测试数据库连接(超时为3秒)
	context, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := RedisClient.Ping(context).Result()
	return err
}

// SetCache 设置缓存及其过期时间
func SetCache(ctx context.Context, key string, value interface{}, expire time.Duration) error {
	return RedisClient.Set(ctx, key, value, expire).Err()
}

// GetCache 获取缓存
func GetCache(ctx context.Context, key string) (string, error) {
	return RedisClient.Get(ctx, key).Result()
}

// DeleteCache 删除缓存
func DeleteCache(ctx context.Context, key string) error {
	return RedisClient.Del(ctx, key).Err()
}

// ExistsCache 检验缓存是否存在
func ExistsCache(ctx context.Context, key string) bool {
	_, err := RedisClient.Exists(ctx, key).Result()
	return !errors.Is(err, redis.Nil)
}
