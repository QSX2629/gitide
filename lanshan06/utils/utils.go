package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("simple_jwt_secret")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 根据用户名和过期时间，生成一个 JWT 令牌字符串，返回给客户端
func GenerateToken(username string, expirationTime time.Time) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime), //令牌过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),     //令牌发放时间（现在的时间）
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //创建令牌（HS256是算法）
	return token.SignedString(jwtSecret)
}

// 客户端请求接口时，携带令牌，服务器解析令牌，验证有效性（是否过期、签名是否正确），并获取里面的用户名。
func ParseToken(tokenStr string) (*Claims, error) {
	//解析令牌（传入令牌字符串、载荷结构体指针、签名验证函数）
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		},
	)
	//：处理解析错误（比如令牌格式错误、过期等）
	if err != nil {
		return nil, err
	}
	//验证令牌是否有效，并转换载荷类型
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid

}
