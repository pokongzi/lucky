package model

import "time"

// RefreshToken 刷新Token模型
type RefreshToken struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"not null;index;comment:用户ID"`
	Token     string    `gorm:"uniqueIndex;size:255;not null;comment:刷新token"`
	ExpiresAt time.Time `gorm:"not null;index;comment:过期时间"`
	IsRevoked bool      `gorm:"not null;default:false;comment:是否已撤销"`
	UserAgent string    `gorm:"size:500;comment:用户代理"`
	ClientIP  string    `gorm:"size:45;comment:客户端IP"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// 关联用户
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
