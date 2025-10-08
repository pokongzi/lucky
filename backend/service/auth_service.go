package service

import (
	"errors"

	"lucky/common/jwt"
	"lucky/model"

	"gorm.io/gorm"
)

// AuthService 处理访问令牌的校验
type AuthService struct {
	db *gorm.DB
}

// NewAuthService 创建鉴权服务
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// ValidateAccessToken 校验访问令牌并返回用户
func (s *AuthService) ValidateAccessToken(token string) (*model.User, error) {
	claims, err := jwt.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// 根据 claims 中的用户信息查询数据库
	userDAO := model.NewUserDAO(s.db)
	user, err := userDAO.GetByID(int64(claims.UserID))
	if err != nil {
		return nil, err
	}

	// 基础一致性检查（可扩展：状态、tokenVersion 等）
	if user.OpenID != claims.OpenID {
		return nil, errors.New("token与用户不匹配")
	}

	return user, nil
}
