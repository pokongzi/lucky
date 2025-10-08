package model

import (
	"time"

	"gorm.io/gorm"
)

// UserNumber 用户号码表
type UserNumber struct {
	ID        int64       `gorm:"primaryKey;column:id" json:"id"`
	UserID    int64       `gorm:"not null;index;column:user_id" json:"user_id"`           // 用户ID
	GameID    uint64      `gorm:"not null;index;column:game_id" json:"game_id"`           // 游戏ID
	Game      LotteryGame `gorm:"foreignKey:GameID" json:"game"`                          // 游戏信息
	RedBalls  NumberArray `gorm:"type:json;not null;column:red_balls" json:"red_balls"`   // 红球号码JSON数组
	BlueBalls NumberArray `gorm:"type:json;not null;column:blue_balls" json:"blue_balls"` // 蓝球号码JSON数组
	Nickname  string      `gorm:"size:128;column:nickname" json:"nickname"`               // 用户给号码起的昵称
	Source    string      `gorm:"size:32;default:'manual';column:source" json:"source"`   // 来源：manual(手动), random(机选)
	IsActive  bool        `gorm:"default:true;column:is_active" json:"is_active"`         // 是否启用
	CreatedAt time.Time   `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time   `gorm:"column:updated_at" json:"updated_at"`
}

func (UserNumber) TableName() string {
	return "user_numbers"
}

// UserNumberDAO 用户号码数据访问对象
type UserNumberDAO struct {
	db *gorm.DB
}

func NewUserNumberDAO(db *gorm.DB) *UserNumberDAO {
	return &UserNumberDAO{db: db}
}

// Create 创建用户号码
func (dao *UserNumberDAO) Create(userNumber *UserNumber) error {
	return dao.db.Create(userNumber).Error
}

// GetByID 根据ID获取用户号码
func (dao *UserNumberDAO) GetByID(id int64) (*UserNumber, error) {
	var userNumber UserNumber
	err := dao.db.First(&userNumber, id).Error
	if err != nil {
		return nil, err
	}
	return &userNumber, nil
}

// GetByUserID 根据用户ID获取用户号码列表
func (dao *UserNumberDAO) GetByUserID(userID int64, offset, limit int) ([]*UserNumber, error) {
	var userNumbers []*UserNumber
	err := dao.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&userNumbers).Error
	return userNumbers, err
}

// GetByGameID 根据游戏ID获取用户号码列表
func (dao *UserNumberDAO) GetByGameID(gameID uint64, offset, limit int) ([]*UserNumber, error) {
	var userNumbers []*UserNumber
	err := dao.db.Where("game_id = ?", gameID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&userNumbers).Error
	return userNumbers, err
}

// GetByUserIDAndGameID 根据用户ID和游戏ID获取用户号码列表
func (dao *UserNumberDAO) GetByUserIDAndGameID(userID int64, gameID uint64, offset, limit int) ([]*UserNumber, error) {
	var userNumbers []*UserNumber
	err := dao.db.Where("user_id = ? AND game_id = ?", userID, gameID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&userNumbers).Error
	return userNumbers, err
}

// GetActiveByUserID 根据用户ID获取启用的用户号码列表
func (dao *UserNumberDAO) GetActiveByUserID(userID int64) ([]*UserNumber, error) {
	var userNumbers []*UserNumber
	err := dao.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&userNumbers).Error
	return userNumbers, err
}

// Update 更新用户号码
func (dao *UserNumberDAO) Update(userNumber *UserNumber) error {
	return dao.db.Save(userNumber).Error
}

// UpdateStatus 更新用户号码状态
func (dao *UserNumberDAO) UpdateStatus(id int64, isActive bool) error {
	return dao.db.Model(&UserNumber{}).Where("id = ?", id).Update("is_active", isActive).Error
}

// UpdateNickname 更新用户号码昵称
func (dao *UserNumberDAO) UpdateNickname(id int64, nickname string) error {
	return dao.db.Model(&UserNumber{}).Where("id = ?", id).Update("nickname", nickname).Error
}

// Delete 删除用户号码
func (dao *UserNumberDAO) Delete(id int64) error {
	return dao.db.Delete(&UserNumber{}, id).Error
}

// DeleteByUserID 根据用户ID删除用户号码
func (dao *UserNumberDAO) DeleteByUserID(userID int64) error {
	return dao.db.Where("user_id = ?", userID).Delete(&UserNumber{}).Error
}

// DeleteByGameID 根据游戏ID删除用户号码
func (dao *UserNumberDAO) DeleteByGameID(gameID uint64) error {
	return dao.db.Where("game_id = ?", gameID).Delete(&UserNumber{}).Error
}

// List 获取用户号码列表
func (dao *UserNumberDAO) List(offset, limit int) ([]*UserNumber, error) {
	var userNumbers []*UserNumber
	err := dao.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&userNumbers).Error
	return userNumbers, err
}

// Count 获取用户号码总数
func (dao *UserNumberDAO) Count() (int64, error) {
	var count int64
	err := dao.db.Model(&UserNumber{}).Count(&count).Error
	return count, err
}

// CountByUserID 根据用户ID获取用户号码总数
func (dao *UserNumberDAO) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := dao.db.Model(&UserNumber{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// CountByGameID 根据游戏ID获取用户号码总数
func (dao *UserNumberDAO) CountByGameID(gameID uint64) (int64, error) {
	var count int64
	err := dao.db.Model(&UserNumber{}).Where("game_id = ?", gameID).Count(&count).Error
	return count, err
}
