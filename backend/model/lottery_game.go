package model

import "time"

// LotteryGame 彩票游戏表
type LotteryGame struct {
	ID              uint64 `gorm:"primaryKey"`
	GameCode        string `gorm:"uniqueIndex;size:32;not null"` // ssq, dlt
	GameName        string `gorm:"size:64;not null"`             // 双色球, 大乐透
	RedBallCount    int    `gorm:"not null"`                     // 红球总数
	BlueBallCount   int    `gorm:"not null"`                     // 蓝球总数
	RedSelectCount  int    `gorm:"not null"`                     // 红球选择数
	BlueSelectCount int    `gorm:"not null"`                     // 蓝球选择数
	IsActive        bool   `gorm:"default:true"`                 // 是否启用
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
