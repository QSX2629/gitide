package utils_lock

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 全局Redis
var RedisClient *redis.Client

// ************************* 分布式锁配置 ******************************//
// 分布式锁前缀
const lockPrefix = "dist_lock:"

// 生成唯一的锁值
func generateLockValue() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// ******************* 分布式锁 Lua 脚本（注意`	`!） ****************************//
// lockLuaScript: 获取分布式锁的Lua脚本（实现SET NX EX的原子操作）
const lockLuaScript = `
SET NX PX是仅当key不存在时设置值，同时设置过期时间
return redis.call('SET', KEYS[1], ARGV[1], 'NX', 'PX', ARGV[2])
`

// unlockLuaScript: 解锁的Lua脚本（实现“检查值-删除key”的原子操作）
const unlockLuaScript = `
验证锁的持有者是否是当前客户端
if redis.call('GET', KEYS[1]) ~= ARGV[1] then
    return 0
end
删除锁
return redis.call('DEL', KEYS[1])
`

// 预编译Lua脚本，提升执行效率
var (
	lockScript   = redis.NewScript(lockLuaScript)
	unlockScript = redis.NewScript(unlockLuaScript)
)

// Lock 尝试获取分布式锁（使用Lua脚本保证原子性）
func Lock(ctx context.Context, key string, expire time.Duration) (string, bool, error) {
	// 拼接完整的锁key
	lockKey := fmt.Sprintf("%s%s", lockPrefix, key)
	// 生成唯一的锁值（用于验证持有者）
	lockValue := generateLockValue()
	// 过期时间转换为毫秒（Lua脚本中PX参数接受毫秒）
	expireMs := int64(expire / time.Millisecond)

	// 执行Lua脚本
	result, err := lockScript.Run(ctx, RedisClient, []string{lockKey}, lockValue, expireMs).Result()
	if err != nil {
		return "", false, fmt.Errorf("获取锁失败：%v", err)
	}

	// 脚本返回nil表示设置失败（锁已存在），返回"OK"表示成功
	ok := result == "OK"
	return lockValue, ok, nil
}

// Unlock 解锁（使用Lua脚本保证“检查-删除”的原子性）
func Unlock(ctx context.Context, key string, value string) error {
	lockKey := fmt.Sprintf("%s%s", lockPrefix, key)

	// 执行Lua脚本
	result, err := unlockScript.Run(ctx, RedisClient, []string{lockKey}, value).Result()
	if err != nil {
		return fmt.Errorf("解锁失败：%v", err)
	}

	// 解析脚本返回结果
	switch res := result.(type) {
	case int64:
		// 返回0：锁不是当前客户端持有；返回1：解锁成功
		if res == 0 {
			return errors.New("不是当前客户端持有的锁或锁已过期")
		}
	case nil:
		return errors.New("锁已过期")
	default:
		return fmt.Errorf("解锁结果异常：%v", result)
	}

	return nil
}

// *************************分布式限流配置 ***************************//
// 限流锁的前缀
const rateLimitLockPrefix = "rate_limit_lock:"

// ***************************分布式限流 Lua 脚本 ************************//
// rateLimitLuaScript: 令牌桶限流的Lua脚本（将所有操作打包为原子操作）
const rateLimitLuaScript = `
令牌桶key和限流锁key
local limitKey = KEYS[1]
local lockKey = KEYS[2]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local interval = tonumber(ARGV[4])
local lockExpireMs = tonumber(ARGV[5])

 1. 获取限流锁（SET NX PX）
local lockOk = redis.call('SET', lockKey, '1', 'NX', 'PX', lockExpireMs)
if not lockOk then
    未获取到锁，返回{false, 0}
    return {0, 0}
end

 2. 获取令牌桶的last_time和tokens
local lastTime = tonumber(redis.call('HGET', limitKey, 'last_time')) or now
local tokens = tonumber(redis.call('HGET', limitKey, 'tokens')) or capacity

 3. 计算生成的令牌数
local generateTokens = math.floor((now - lastTime) / interval)
if generateTokens > 0 then
    tokens = math.min(tokens + generateTokens, capacity)
    lastTime = now
end

 4. 判断是否允许访问（消耗令牌）
local allow = 0
if tokens > 0 then
    tokens = tokens - 1
    allow = 1
end

 5. 更新令牌桶数据，并设置过期时间
redis.call('HSET', limitKey, 'last_time', lastTime, 'tokens', tokens)
redis.call('EXPIRE', limitKey, 60) -- 60秒过期

 6. 释放限流锁
redis.call('DEL', lockKey)

 返回结果：{是否允许(1/0), 剩余令牌数}
return {allow, tokens}
`

// 预编译限流Lua脚本
var rateLimitScript = redis.NewScript(rateLimitLuaScript)

// RateLimit 分布式限流（使用Lua脚本保证原子性）
func RateLimit(ctx context.Context, key string, capacity int, rate float64) (bool, int, error) {
	// 令牌桶的key和限流锁的key
	limitKey := fmt.Sprintf("rate_limit:%s", key)
	lockKey := fmt.Sprintf("%s%s", rateLimitLockPrefix, key)

	// 当前时间戳（毫秒）
	now := time.Now().UnixMilli()
	// 每个令牌的生成间隔（毫秒），防止rate为0导致除零错误
	if rate <= 0 {
		return false, 0, errors.New("令牌生成速率不能为0或负数")
	}
	interval := int64(1000 / rate)
	if interval == 0 {
		interval = 1 // 避免间隔为0（比如rate>1000时）
	}

	// 限流锁过期时间（1秒=1000毫秒）
	lockExpireMs := int64(1000)

	// 执行Lua脚本
	result, err := rateLimitScript.Run(
		ctx,
		RedisClient,
		[]string{limitKey, lockKey},                 // KEYS：[limitKey, lockKey]
		capacity, rate, now, interval, lockExpireMs, // ARGV：入参列表
	).Result()
	if err != nil {
		return false, 0, fmt.Errorf("执行限流脚本失败：%v", err)
	}

	// 解析脚本返回的结果（Lua返回的table对应Go的[]interface{}）
	resultSlice, ok := result.([]interface{})
	if !ok {
		return false, 0, errors.New("限流脚本返回结果格式错误")
	}
	if len(resultSlice) != 2 {
		return false, 0, errors.New("限流脚本返回结果长度错误")
	}

	// 解析是否允许访问（1=允许，0=拒绝）
	allowInt, ok := resultSlice[0].(int64)
	if !ok {
		return false, 0, errors.New("限流脚本返回allow格式错误")
	}
	allow := allowInt == 1

	// 解析剩余令牌数
	tokensInt, ok := resultSlice[1].(int64)
	if !ok {
		return false, 0, errors.New("限流脚本返回tokens格式错误")
	}
	tokens := int(tokensInt)

	return allow, tokens, nil
}
