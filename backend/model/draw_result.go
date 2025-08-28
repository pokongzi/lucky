package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// NumberArray 自定义类型用于存储号码数组
type NumberArray []int

// Scan 实现 Scanner 接口
func (na *NumberArray) Scan(value interface{}) error {
	if value == nil {
		*na = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, na)
	case string:
		return json.Unmarshal([]byte(v), na)
	default:
		return fmt.Errorf("cannot scan %T into NumberArray", value)
	}
}

// Value 实现 Valuer 接口
func (na NumberArray) Value() (driver.Value, error) {
	if na == nil {
		return nil, nil
	}
	return json.Marshal(na)
}

// DrawResult 开奖结果表
type DrawResult struct {
	ID           uint64      `gorm:"primaryKey"`
	GameID       uint64      `gorm:"not null;index"`         // 游戏ID
	Period       string      `gorm:"size:32;not null;index"` // 期号，如：2023130
	DrawDate     time.Time   `gorm:"not null;index"`         // 开奖日期
	RedBalls     NumberArray `gorm:"type:json;not null"`     // 红球号码 JSON数组
	BlueBalls    NumberArray `gorm:"type:json;not null"`     // 蓝球号码 JSON数组
	SalesAmount  int64       `gorm:"default:0"`              // 销售额（分）
	PrizePool    int64       `gorm:"default:0"`              // 奖池金额（分）
	FirstPrize   int         `gorm:"default:0"`              // 一等奖注数
	FirstAmount  int64       `gorm:"default:0"`              // 一等奖单注奖金（分）
	SecondPrize  int         `gorm:"default:0"`              // 二等奖注数
	SecondAmount int64       `gorm:"default:0"`              // 二等奖单注奖金（分）
	ThirdPrize   int         `gorm:"default:0"`              // 三等奖注数
	ThirdAmount  int64       `gorm:"default:0"`              // 三等奖单注奖金（分）
	FourthPrize  int         `gorm:"default:0"`              // 四等奖注数
	FourthAmount int64       `gorm:"default:0"`              // 四等奖单注奖金（分）
	FifthPrize   int         `gorm:"default:0"`              // 五等奖注数
	FifthAmount  int64       `gorm:"default:0"`              // 五等奖单注奖金（分）
	SixthPrize   int         `gorm:"default:0"`              // 六等奖注数
	SixthAmount  int64       `gorm:"default:0"`              // 六等奖单注奖金（分）
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// 关联
	Game LotteryGame `gorm:"foreignKey:GameID"`
}

// 复合索引
func (DrawResult) TableName() string {
	return "draw_results"
}
