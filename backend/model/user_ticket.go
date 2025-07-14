package model

import "time"

type UserTicket struct {
	ID        uint64 `gorm:"primaryKey"`
	UserID    uint64 `gorm:"not null;index"`
	Type      int    `gorm:"not null"` // 1: 双色球, 2: 大乐透
	Numbers   string `gorm:"size:64;not null"`
	CreatedAt time.Time
}
