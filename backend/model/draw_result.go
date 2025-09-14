package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
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
	ID           uint64      `gorm:"primaryKey;column:id" json:"id"`
	GameID       uint64      `gorm:"not null;index;column:game_id" json:"game_id"`           // 游戏ID
	Period       string      `gorm:"size:32;not null;column:period" json:"period"`           // 期号
	RedBalls     NumberArray `gorm:"type:json;not null;column:red_balls" json:"red_balls"`   // 红球号码JSON数组
	BlueBalls    NumberArray `gorm:"type:json;not null;column:blue_balls" json:"blue_balls"` // 蓝球号码JSON数组
	DrawDate     time.Time   `gorm:"not null;column:draw_date" json:"draw_date"`             // 开奖时间
	SalesAmount  int64       `gorm:"default:0;column:sales_amount" json:"sales_amount"`      // 销售额(分)
	PrizePool    int64       `gorm:"default:0;column:prize_pool" json:"prize_pool"`          // 奖池金额(分)
	FirstPrize   int         `gorm:"default:0;column:first_prize" json:"first_prize"`        // 一等奖注数
	FirstAmount  int64       `gorm:"default:0;column:first_amount" json:"first_amount"`      // 一等奖单注奖金(分)
	SecondPrize  int         `gorm:"default:0;column:second_prize" json:"second_prize"`      // 二等奖注数
	SecondAmount int64       `gorm:"default:0;column:second_amount" json:"second_amount"`    // 二等奖单注奖金(分)
	CreatedAt    time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"column:updated_at" json:"updated_at"`

	// 关联
	Game LotteryGame `gorm:"foreignKey:GameID;references:ID;constraint:OnDelete:CASCADE"`
}

func (DrawResult) TableName() string {
	return "draw_results"
}

// DrawResultDAO 开奖结果数据访问对象
type DrawResultDAO struct {
	db *gorm.DB
}

func NewDrawResultDAO(db *gorm.DB) *DrawResultDAO {
	return &DrawResultDAO{db: db}
}

// Create 创建开奖结果
func (dao *DrawResultDAO) Create(result *DrawResult) error {
	return dao.db.Create(result).Error
}

// GetByID 根据ID获取开奖结果
func (dao *DrawResultDAO) GetByID(id uint64) (*DrawResult, error) {
	var result DrawResult
	err := dao.db.First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetByPeriod 根据期号获取开奖结果
func (dao *DrawResultDAO) GetByPeriod(period string) (*DrawResult, error) {
	var result DrawResult
	err := dao.db.Where("period = ?", period).First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetByGameID 根据游戏ID获取开奖结果列表
func (dao *DrawResultDAO) GetByGameID(gameID uint64, offset, limit int) ([]*DrawResult, error) {
	var results []*DrawResult
	err := dao.db.Where("game_id = ?", gameID).Offset(offset).Limit(limit).Order("draw_time DESC").Find(&results).Error
	return results, err
}

// GetLatestByGameID 根据游戏ID获取最新开奖结果
func (dao *DrawResultDAO) GetLatestByGameID(gameID uint64) (*DrawResult, error) {
	var result DrawResult
	err := dao.db.Where("game_id = ?", gameID).Order("draw_time DESC").First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Update 更新开奖结果
func (dao *DrawResultDAO) Update(result *DrawResult) error {
	return dao.db.Save(result).Error
}

// UpdateStatus 更新开奖结果状态
func (dao *DrawResultDAO) UpdateStatus(id uint64, isActive bool) error {
	return dao.db.Model(&DrawResult{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// Delete 删除开奖结果
func (dao *DrawResultDAO) Delete(id uint64) error {
	return dao.db.Delete(&DrawResult{}, id).Error
}

// DeleteByGameID 根据游戏ID删除开奖结果
func (dao *DrawResultDAO) DeleteByGameID(gameID uint64) error {
	return dao.db.Where("game_id = ?", gameID).Delete(&DrawResult{}).Error
}

// List 获取开奖结果列表
func (dao *DrawResultDAO) List(offset, limit int) ([]*DrawResult, error) {
	var results []*DrawResult
	err := dao.db.Offset(offset).Limit(limit).Order("draw_time DESC").Find(&results).Error
	return results, err
}

// Count 获取开奖结果总数
func (dao *DrawResultDAO) Count() (int64, error) {
	var count int64
	err := dao.db.Model(&DrawResult{}).Count(&count).Error
	return count, err
}

// CountByGameID 根据游戏ID获取开奖结果总数
func (dao *DrawResultDAO) CountByGameID(gameID uint64) (int64, error) {
	var count int64
	err := dao.db.Model(&DrawResult{}).Where("game_id = ?", gameID).Count(&count).Error
	return count, err
}
