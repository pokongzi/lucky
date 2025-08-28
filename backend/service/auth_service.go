package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net"
	"strings"
	"time"

	"lucky/common/jwt"
	"lucky/model"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrUserDisabled = errors.New("用户已被禁用")
	ErrTokenExpired = errors.New("token已过期")
	ErrTokenRevoked = errors.New("token已被撤销")
	ErrInvalidToken = errors.New("无效的token")
)

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresAt    time.Time   `json:"expires_at"`
	User         *model.User `json:"user"`
}

// AuthService JWT认证服务
type AuthService struct {
	db *gorm.DB
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Login 用户登录（微信登录）
func (s *AuthService) Login(openID, nickname, avatarURL, clientIP, userAgent string) (*LoginResult, error) {
	// 获取或创建用户
	user, err := GetOrCreateUser(s.db, openID, nickname, avatarURL)
	if err != nil {
		// 记录登录失败日志
		s.logLogin(0, "wechat", clientIP, userAgent, 0, err.Error())
		return nil, err
	}

	// 检查用户状态
	if user.Status != 1 {
		s.logLogin(user.ID, "wechat", clientIP, userAgent, 0, "用户已被禁用")
		return nil, ErrUserDisabled
	}

	// 更新用户登录信息
	now := time.Now()
	err = s.db.Model(user).Updates(map[string]interface{}{
		"last_login_at": &now,
		"last_login_ip": clientIP,
		"login_count":   gorm.Expr("login_count + 1"),
	}).Error
	if err != nil {
		return nil, err
	}

	// 生成tokens
	result, err := s.generateTokens(user, clientIP, userAgent)
	if err != nil {
		s.logLogin(user.ID, "wechat", clientIP, userAgent, 0, err.Error())
		return nil, err
	}

	// 记录登录成功日志
	s.logLogin(user.ID, "wechat", clientIP, userAgent, 1, "")

	return result, nil
}

// RefreshAccessToken 刷新访问token
func (s *AuthService) RefreshAccessToken(refreshToken, clientIP, userAgent string) (*LoginResult, error) {
	// 查找刷新token
	var tokenRecord model.RefreshToken
	err := s.db.Where("token = ? AND is_revoked = false", refreshToken).
		Preload("User").First(&tokenRecord).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	// 检查token是否过期
	if time.Now().After(tokenRecord.ExpiresAt) {
		// 撤销过期token
		s.db.Model(&tokenRecord).Update("is_revoked", true)
		s.logLogin(tokenRecord.UserID, "refresh", clientIP, userAgent, 0, "刷新token已过期")
		return nil, ErrTokenExpired
	}

	// 检查用户状态
	if tokenRecord.User.Status != 1 {
		s.logLogin(tokenRecord.UserID, "refresh", clientIP, userAgent, 0, "用户已被禁用")
		return nil, ErrUserDisabled
	}

	// 撤销旧的刷新token
	err = s.db.Model(&tokenRecord).Update("is_revoked", true).Error
	if err != nil {
		return nil, err
	}

	// 生成新的tokens
	result, err := s.generateTokens(&tokenRecord.User, clientIP, userAgent)
	if err != nil {
		s.logLogin(tokenRecord.UserID, "refresh", clientIP, userAgent, 0, err.Error())
		return nil, err
	}

	// 记录刷新成功日志
	s.logLogin(tokenRecord.UserID, "refresh", clientIP, userAgent, 1, "")

	return result, nil
}

// ValidateAccessToken 验证访问token
func (s *AuthService) ValidateAccessToken(accessToken string) (*model.User, error) {
	claims, err := jwt.ValidateToken(accessToken)
	if err != nil {
		return nil, err
	}

	// 查询用户并验证token版本
	var user model.User
	err = s.db.First(&user, claims.UserID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 检查token版本（用于强制退出登录）
	if claims.RegisteredClaims.IssuedAt.Before(time.Now().Add(-time.Duration(user.TokenVersion) * 24 * time.Hour)) {
		return nil, ErrTokenExpired
	}

	return &user, nil
}

// RevokeUserTokens 撤销用户所有token（强制退出登录）
func (s *AuthService) RevokeUserTokens(userID uint64) error {
	// 增加token版本号，使所有旧的访问token失效
	err := s.db.Model(&model.User{}).Where("id = ?", userID).
		Update("token_version", gorm.Expr("token_version + 1")).Error
	if err != nil {
		return err
	}

	// 撤销所有刷新token
	return s.db.Model(&model.RefreshToken{}).Where("user_id = ? AND is_revoked = false", userID).
		Update("is_revoked", true).Error
}

// CleanExpiredTokens 清理过期的刷新token
func (s *AuthService) CleanExpiredTokens() error {
	return s.db.Where("expires_at < ? OR is_revoked = true", time.Now().Add(-24*time.Hour)).
		Delete(&model.RefreshToken{}).Error
}

// generateTokens 生成访问token和刷新token
func (s *AuthService) generateTokens(user *model.User, clientIP, userAgent string) (*LoginResult, error) {
	// 生成访问token
	accessToken, err := jwt.GenerateToken(user.ID, user.OpenID, user.Nickname)
	if err != nil {
		return nil, err
	}

	// 生成刷新token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	// 保存刷新token到数据库
	tokenRecord := model.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(jwt.RefreshTokenExpireDuration),
		UserAgent: userAgent,
		ClientIP:  clientIP,
	}

	err = s.db.Create(&tokenRecord).Error
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(jwt.TokenExpireDuration),
		User:         user,
	}, nil
}

// generateRefreshToken 生成随机刷新token
func (s *AuthService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// logLogin 记录登录日志
func (s *AuthService) logLogin(userID uint64, loginType, clientIP, userAgent string, status int, errorMsg string) {
	// 处理IPv6地址
	if ip := net.ParseIP(clientIP); ip != nil && ip.To4() == nil {
		// IPv6地址，截取前45个字符
		if len(clientIP) > 45 {
			clientIP = clientIP[:45]
		}
	}

	// 截取UserAgent长度
	if len(userAgent) > 500 {
		userAgent = userAgent[:500]
	}

	// 截取错误信息长度
	if len(errorMsg) > 255 {
		errorMsg = errorMsg[:255]
	}

	loginLog := model.LoginLog{
		UserID:    userID,
		LoginType: loginType,
		ClientIP:  clientIP,
		UserAgent: userAgent,
		Status:    status,
		ErrorMsg:  errorMsg,
		LoginAt:   time.Now(),
	}

	// 异步记录日志，不影响主要流程
	go func() {
		s.db.Create(&loginLog)
	}()
}

// GetClientIP 获取客户端真实IP
func GetClientIP(request interface{}) string {
	// 这里需要根据具体的请求对象实现
	// 通常从 X-Forwarded-For, X-Real-IP 等 header 获取
	return "127.0.0.1" // 默认值
}

// ExtractBearerToken 从Authorization header提取Bearer token
func ExtractBearerToken(authHeader string) string {
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
