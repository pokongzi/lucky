package model

import (
	"time"

	"gorm.io/gorm"
)

// LotteryGame 彩票游戏表
type LotteryGame struct {
	ID              uint64    `gorm:"primaryKey;column:id" json:"id"`
	GameCode        string    `gorm:"size:32;not null;column:game_code" json:"game_code"`         // 游戏代码
	GameName        string    `gorm:"size:64;not null;column:game_name" json:"game_name"`         // 游戏名称
	RedBallCount    int       `gorm:"not null;column:red_ball_count" json:"red_ball_count"`       // 红球总数
	BlueBallCount   int       `gorm:"not null;column:blue_ball_count" json:"blue_ball_count"`     // 蓝球总数
	RedSelectCount  int       `gorm:"not null;column:red_select_count" json:"red_select_count"`   // 红球选择数
	BlueSelectCount int       `gorm:"not null;column:blue_select_count" json:"blue_select_count"` // 蓝球选择数
	IsActive        bool      `gorm:"default:true;column:is_active" json:"is_active"`             // 是否启用
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
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
	err := dao.db.Where("game_name = ?", name).First(&game).Error
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
