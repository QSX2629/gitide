package utils

import (
	"time"

	// 只保留一个JWT库（推荐使用较新的v5版本，删除重复的dgrijalva库）
	"github.com/golang-jwt/jwt/v5"
)

// 1. 定义JWT密钥（修复“2个用法”的警告，确保只在本包内使用）
var jwtSecret = []byte("simple-jwt-secret")

// 2. 定义令牌载荷结构体（修复“可导出的字段”警告，首字母大写）
type Claims struct {
	Username             string `json:"username"` // 用户名（首字母大写，支持JSON序列化）
	jwt.RegisteredClaims        // 替换StandardClaims（v5版本推荐用法）
}

// 3. 生成JWT令牌（支持自定义过期时间，修复“未使用”警告）
// 参数：
//   - username：用户名
//   - expirationTime：过期时间（time.Time类型）
//
// 返回：令牌字符串和错误信息
func GenerateToken(username string, expirationTime time.Time) (string, error) {
	// 创建载荷
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), // 过期时间（v5版本需用NumericDate包装）
			IssuedAt:  jwt.NewNumericDate(time.Now()),     // 签发时间
		},
	}

	// 生成令牌（使用HS256算法签名）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥签名并返回字符串
	return token.SignedString(jwtSecret)
}

// 4. 补充解析令牌函数（与GenerateToken配套，用于中间件验证）
func ParseToken(tokenStr string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil // 使用相同密钥验证签名
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌有效性并返回载荷
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid // 令牌无效
}
