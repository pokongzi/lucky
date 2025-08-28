package model

import "time"

// LoginLog 登录日志模型
type LoginLog struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"not null;index;comment:用户ID"`
	LoginType string    `gorm:"size:32;not null;comment:登录类型(wechat,refresh)"`
	ClientIP  string    `gorm:"size:45;not null;comment:客户端IP"`
	UserAgent string    `gorm:"size:500;comment:用户代理"`
	Status    int       `gorm:"not null;comment:登录状态(1:成功 0:失败)"`
	ErrorMsg  string    `gorm:"size:255;comment:错误信息"`
	LoginAt   time.Time `gorm:"not null;index;comment:登录时间"`
	CreatedAt time.Time

	// 关联用户
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}
