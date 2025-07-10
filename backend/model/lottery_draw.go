package model

import "time"

type LotteryDraw struct {
	ID         uint64    `gorm:"primaryKey"`
	Type       int       `gorm:"not null"` // 1: 双色球, 2: 大乐透
	Period     string    `gorm:"size:16;not null;index:idx_type_period,unique"`
	DrawDate   time.Time `gorm:"not null"`
	Numbers    string    `gorm:"size:64;not null"`
	PoolAmount float64   `gorm:"type:decimal(16,2)"`
	Details    string    `gorm:"type:text"`
	CreatedAt  time.Time
}
