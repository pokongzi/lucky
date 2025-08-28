package model

import "time"

type User struct {
	ID           uint64     `gorm:"primaryKey"`
	OpenID       string     `gorm:"uniqueIndex;size:64;not null"`
	Nickname     string     `gorm:"size:64;not null"`
	AvatarURL    string     `gorm:"size:255"`
	Status       int        `gorm:"not null;default:1;comment:用户状态(1:正常 0:禁用)"` // 用户状态
	TokenVersion int        `gorm:"not null;default:1;comment:token版本号"`        // Token版本号，用于强制退出登录
	LastLoginAt  *time.Time `gorm:"comment:最后登录时间"`                             // 最后登录时间
	LastLoginIP  string     `gorm:"size:45;comment:最后登录IP"`                     // 最后登录IP
	LoginCount   int        `gorm:"not null;default:0;comment:登录次数"`            // 登录次数统计
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
