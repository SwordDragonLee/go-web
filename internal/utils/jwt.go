package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// JWT密钥（生产环境应从配置文件读取）
	jwtSecret = []byte("your-secret-key-change-in-production")
	// Token过期时间（7天）
	tokenExpireDuration = 7 * 24 * time.Hour
)

// Claims JWT Claims结构
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint, username, role string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(tokenExpireDuration)

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "go-web",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetTokenExpireDuration 获取Token过期时间（秒）
func GetTokenExpireDuration() int64 {
	return int64(tokenExpireDuration.Seconds())
}

// SetJWTSecret 设置JWT密钥（用于从配置读取）
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// SetTokenExpireDuration 设置Token过期时间
func SetTokenExpireDuration(duration time.Duration) {
	tokenExpireDuration = duration
}

