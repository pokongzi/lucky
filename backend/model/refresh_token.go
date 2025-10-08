package model

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken 刷新令牌表
type RefreshToken struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	UserID    int64     `gorm:"not null;index;column:user_id" json:"user_id"`            // 用户ID
	Token     string    `gorm:"uniqueIndex;size:255;not null;column:token" json:"token"` // 刷新令牌
	ExpiresAt time.Time `gorm:"not null;column:expires_at" json:"expires_at"`            // 过期时间
	IsActive  bool      `gorm:"default:true;column:is_active" json:"is_active"`          // 是否有效
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// RefreshTokenDAO 刷新令牌数据访问对象
type RefreshTokenDAO struct {
	db *gorm.DB
}

func NewRefreshTokenDAO(db *gorm.DB) *RefreshTokenDAO {
	return &RefreshTokenDAO{db: db}
}

// Create 创建刷新令牌
func (dao *RefreshTokenDAO) Create(token *RefreshToken) error {
	return dao.db.Create(token).Error
}

// GetByID 根据ID获取刷新令牌
func (dao *RefreshTokenDAO) GetByID(id int64) (*RefreshToken, error) {
	var token RefreshToken
	err := dao.db.First(&token, id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetByToken 根据令牌获取刷新令牌
func (dao *RefreshTokenDAO) GetByToken(token string) (*RefreshToken, error) {
	var refreshToken RefreshToken
	err := dao.db.Where("token = ?", token).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// GetByUserID 根据用户ID获取刷新令牌列表
func (dao *RefreshTokenDAO) GetByUserID(userID int64) ([]*RefreshToken, error) {
	var tokens []*RefreshToken
	err := dao.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&tokens).Error
	return tokens, err
}

// Update 更新刷新令牌
func (dao *RefreshTokenDAO) Update(token *RefreshToken) error {
	return dao.db.Save(token).Error
}

// UpdateStatus 更新令牌状态
func (dao *RefreshTokenDAO) UpdateStatus(id int64, isActive bool) error {
	return dao.db.Model(&RefreshToken{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// InvalidateByUserID 使用户的所有令牌失效
func (dao *RefreshTokenDAO) InvalidateByUserID(userID int64) error {
	return dao.db.Model(&RefreshToken{}).Where("user_id = ?", userID).Update("is_active", false).Error
}

// Delete 删除刷新令牌
func (dao *RefreshTokenDAO) Delete(id int64) error {
	return dao.db.Delete(&RefreshToken{}, id).Error
}

// DeleteByToken 根据令牌删除刷新令牌
func (dao *RefreshTokenDAO) DeleteByToken(token string) error {
	return dao.db.Where("token = ?", token).Delete(&RefreshToken{}).Error
}

// DeleteByUserID 根据用户ID删除刷新令牌
func (dao *RefreshTokenDAO) DeleteByUserID(userID int64) error {
	return dao.db.Where("user_id = ?", userID).Delete(&RefreshToken{}).Error
}

// DeleteExpired 删除过期的令牌
func (dao *RefreshTokenDAO) DeleteExpired() error {
	return dao.db.Where("expires_at < ?", time.Now()).Delete(&RefreshToken{}).Error
}

// List 获取刷新令牌列表
func (dao *RefreshTokenDAO) List(offset, limit int) ([]*RefreshToken, error) {
	var tokens []*RefreshToken
	err := dao.db.Offset(offset).Limit(limit).Find(&tokens).Error
	return tokens, err
}

// Count 获取刷新令牌总数
func (dao *RefreshTokenDAO) Count() (int64, error) {
	var count int64
	err := dao.db.Model(&RefreshToken{}).Count(&count).Error
	return count, err
}
