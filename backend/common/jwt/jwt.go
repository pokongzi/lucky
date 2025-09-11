package jwt

import (
"errors"
"time"

"lucky/common/config"

"github.com/golang-jwt/jwt/v5"
)

// Claims JWT声明结构
type Claims struct {
UserID   uint64 `json:"user_id"`
OpenID   string `json:"open_id"`
Nickname string `json:"nickname"`
jwt.RegisteredClaims
}

var (
ErrTokenExpired = errors.New("token已过期")
ErrTokenInvalid = errors.New("无效的token")
)

// JWTConfig JWT配置结构
type JWTConfig struct {
Secret                string
AccessTokenExpire     time.Duration
RefreshTokenExpire    time.Duration
Issuer                string
Subject               string
}

// getJWTConfig 从配置文件获取JWT配置
func getJWTConfig() *JWTConfig {
section := config.Config.Section("jwt")

// 解析过期时间
accessExpireStr := section.Key("access_token_expire").String()
refreshExpireStr := section.Key("refresh_token_expire").String()

accessExpire, err := time.ParseDuration(accessExpireStr)
if err != nil {
// 默认7天
accessExpire = time.Hour * 24 * 7
}

refreshExpire, err := time.ParseDuration(refreshExpireStr)
if err != nil {
// 默认30天
refreshExpire = time.Hour * 24 * 30
}

return &JWTConfig{
Secret:             section.Key("secret").String(),
AccessTokenExpire:  accessExpire,
RefreshTokenExpire: refreshExpire,
Issuer:             section.Key("issuer").String(),
Subject:            section.Key("subject").String(),
}
}

// GenerateToken 生成JWT Token
func GenerateToken(userID uint64, openID, nickname string) (string, error) {
cfg := getJWTConfig()

// 验证配置
if cfg.Secret == "" {
return "", errors.New("JWT密钥未配置")
}

claims := Claims{
UserID:   userID,
OpenID:   openID,
Nickname: nickname,
RegisteredClaims: jwt.RegisteredClaims{
ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessTokenExpire)),
IssuedAt:  jwt.NewNumericDate(time.Now()),
NotBefore: jwt.NewNumericDate(time.Now()),
Issuer:    cfg.Issuer,
Subject:   cfg.Subject,
},
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析JWT Token
func ParseToken(tokenString string) (*Claims, error) {
cfg := getJWTConfig()

// 验证配置
if cfg.Secret == "" {
return nil, errors.New("JWT密钥未配置")
}

token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
return []byte(cfg.Secret), nil
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

// GetRefreshTokenExpire 获取刷新Token过期时间
func GetRefreshTokenExpire() time.Duration {
cfg := getJWTConfig()
return cfg.RefreshTokenExpire
}

// GetAccessTokenExpire 获取访问Token过期时间
func GetAccessTokenExpire() time.Duration {
cfg := getJWTConfig()
return cfg.AccessTokenExpire
}
