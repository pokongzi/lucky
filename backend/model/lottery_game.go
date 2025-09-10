package model

import (
	"time"

	"gorm.io/gorm"
)

// LotteryGame 彩票游戏表
type LotteryGame struct {
	ID          uint64    `gorm:"primaryKey;column:id" json:"id"`
	Name        string    `gorm:"size:64;not null;column:name" json:"name"`       // 游戏名称
	Description string    `gorm:"size:255;column:description" json:"description"` // 游戏描述
	RedCount    int       `gorm:"not null;column:red_count" json:"red_count"`     // 红球数量
	BlueCount   int       `gorm:"not null;column:blue_count" json:"blue_count"`   // 蓝球数量
	RedRange    int       `gorm:"not null;column:red_range" json:"red_range"`     // 红球范围
	BlueRange   int       `gorm:"not null;column:blue_range" json:"blue_range"`   // 蓝球范围
	IsActive    bool      `gorm:"default:true;column:is_active" json:"is_active"` // 是否启用
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (LotteryGame) TableName() string {
	return "lottery_games"
}

// LotteryGameDAO 彩票游戏数据访问对象
type LotteryGameDAO struct {
	db *gorm.DB
}

func NewLotteryGameDAO(db *gorm.DB) *LotteryGameDAO {
	return &LotteryGameDAO{db: db}
}

// Create 创建彩票游戏
func (dao *LotteryGameDAO) Create(game *LotteryGame) error {
	return dao.db.Create(game).Error
}

// GetByID 根据ID获取彩票游戏
func (dao *LotteryGameDAO) GetByID(id uint64) (*LotteryGame, error) {
	var game LotteryGame
	err := dao.db.First(&game, id).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

// GetByName 根据名称获取彩票游戏
func (dao *LotteryGameDAO) GetByName(name string) (*LotteryGame, error) {
	var game LotteryGame
	err := dao.db.Where("name = ?", name).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

// Update 更新彩票游戏
func (dao *LotteryGameDAO) Update(game *LotteryGame) error {
	return dao.db.Save(game).Error
}

// UpdateStatus 更新游戏状态
func (dao *LotteryGameDAO) UpdateStatus(id uint64, isActive bool) error {
	return dao.db.Model(&LotteryGame{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// Delete 删除彩票游戏
func (dao *LotteryGameDAO) Delete(id uint64) error {
	return dao.db.Delete(&LotteryGame{}, id).Error
}

// List 获取彩票游戏列表
func (dao *LotteryGameDAO) List(offset, limit int) ([]*LotteryGame, error) {
	var games []*LotteryGame
	err := dao.db.Offset(offset).Limit(limit).Find(&games).Error
	return games, err
}

// ListActive 获取启用的彩票游戏列表
func (dao *LotteryGameDAO) ListActive() ([]*LotteryGame, error) {
	var games []*LotteryGame
	err := dao.db.Where("is_active = ?", true).Find(&games).Error
	return games, err
}

// Count 获取彩票游戏总数
func (dao *LotteryGameDAO) Count() (int64, error) {
	var count int64
	err := dao.db.Model(&LotteryGame{}).Count(&count).Error
	return count, err
}
