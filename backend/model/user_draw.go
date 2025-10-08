package model

import (
	"time"

	"gorm.io/gorm"
)

// UserDraw 用户中奖记录
type UserDraw struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserNumberID uint           `json:"user_number_id" gorm:"not null;index"`
	DrawResultID uint           `json:"draw_result_id" gorm:"not null;index"`
	PrizeLevel   int            `json:"prize_level" gorm:"not null;default:0"`
	IsWinning    bool           `json:"is_winning" gorm:"not null;default:false"`
	IsActive     bool           `json:"is_active" gorm:"not null;default:true"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
