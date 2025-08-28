package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构
type Claims struct {
	UserID   uint64 `json:"user_id"`
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
	jwt.RegisteredClaims
}

const (
	// TokenExpireDuration Token过期时间 7天
	TokenExpireDuration = time.Hour * 24 * 7
	// RefreshTokenExpireDuration 刷新Token过期时间 30天
	RefreshTokenExpireDuration = time.Hour * 24 * 30
)

var (
	// JWTSecret JWT密钥，实际项目中应该从配置文件读取
	JWTSecret = []byte("lucky-app-jwt-secret-key-2024")

	ErrTokenExpired = errors.New("token已过期")
	ErrTokenInvalid = errors.New("无效的token")
)

// GenerateToken 生成JWT Token
func GenerateToken(userID uint64, openID, nickname string) (string, error) {
	claims := Claims{
		UserID:   userID,
		OpenID:   openID,
		Nickname: nickname,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "lucky-app",
			Subject:   "user-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新Token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		return "", err
	}

	// 重新生成token
	return GenerateToken(claims.UserID, claims.OpenID, claims.Nickname)
}

// ValidateToken 验证Token有效性
func ValidateToken(tokenString string) (*Claims, error) {
	return ParseToken(tokenString)
}
