package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId               string `json:"user_id"` // 用户ID
	jwt.RegisteredClaims        // 注册声明
}

/* 生成Token */
func GenerateToken(userId string, secret string, expiration time.Duration) (string, error) {
	claims := Claims{
		UserId: userId, // 用户ID
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                 // 发行时间
			Subject:   userId,                                         // 主题
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) // 创建Token
	return token.SignedString([]byte(secret))                  // 使用密钥签名
}

/* 解析Token */
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err // 如果解析失败，则返回错误
	}
	claims, ok := token.Claims.(*Claims) // 获取声明
	// 如果声明无效，则返回错误
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid // 如果签名无效，则返回错误
	}
	return claims, nil // 返回声明
}

/* 从Token中获取用户ID */
func GetUserIdFromToken(tokenString string, secret string) (string, error) {
	claims, err := ParseToken(tokenString, secret)
	if err != nil {
		return "", err // 如果获取用户ID失败，则返回错误
	}
	return claims.UserId, nil // 返回用户ID
}
