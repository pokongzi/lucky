package model

import "time"

// UserNumber 用户号码表
type UserNumber struct {
	ID        uint64      `gorm:"primaryKey"`
	UserID    uint64      `gorm:"not null;index"`           // 用户ID
	GameID    uint64      `gorm:"not null;index"`           // 游戏ID
	RedBalls  NumberArray `gorm:"type:json;not null"`       // 红球号码JSON数组
	BlueBalls NumberArray `gorm:"type:json;not null"`       // 蓝球号码JSON数组
	Nickname  string      `gorm:"size:128"`                 // 用户给号码起的昵称
	Source    string      `gorm:"size:32;default:'manual'"` // 来源：manual(手动), random(机选)
	IsActive  bool        `gorm:"default:true"`             // 是否启用
	CreatedAt time.Time
	UpdatedAt time.Time

	// 关联
	User User        `gorm:"foreignKey:UserID"`
	Game LotteryGame `gorm:"foreignKey:GameID"`
}
