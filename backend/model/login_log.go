package model

import (
	"time"

	"gorm.io/gorm"
)

// LoginLog 登录日志表
type LoginLog struct {
	ID        int64     `gorm:"primaryKey;column:id" json:"id"`
	UserID    int64     `gorm:"not null;index;column:user_id" json:"user_id"`     // 用户ID
	LoginIP   string    `gorm:"size:45;not null;column:login_ip" json:"login_ip"` // 登录IP
	UserAgent string    `gorm:"size:255;column:user_agent" json:"user_agent"`     // 用户代理
	Status    int       `gorm:"not null;default:1;column:status" json:"status"`   // 登录状态(1:成功 0:失败)
	Message   string    `gorm:"size:255;column:message" json:"message"`           // 登录消息
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (LoginLog) TableName() string {
	return "login_logs"
}

// LoginLogDAO 登录日志数据访问对象
type LoginLogDAO struct {
	db *gorm.DB
}

func NewLoginLogDAO(db *gorm.DB) *LoginLogDAO {
	return &LoginLogDAO{db: db}
}

// Create 创建登录日志
func (dao *LoginLogDAO) Create(log *LoginLog) error {
	return dao.db.Create(log).Error
}

// GetByID 根据ID获取登录日志
func (dao *LoginLogDAO) GetByID(id int64) (*LoginLog, error) {
	var log LoginLog
	err := dao.db.First(&log, id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

// GetByUserID 根据用户ID获取登录日志列表
func (dao *LoginLogDAO) GetByUserID(userID int64, offset, limit int) ([]*LoginLog, error) {
	var logs []*LoginLog
	err := dao.db.Where("user_id = ?", userID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetByIP 根据IP获取登录日志列表
func (dao *LoginLogDAO) GetByIP(ip string, offset, limit int) ([]*LoginLog, error) {
	var logs []*LoginLog
	err := dao.db.Where("login_ip = ?", ip).Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// GetByStatus 根据状态获取登录日志列表
func (dao *LoginLogDAO) GetByStatus(status int, offset, limit int) ([]*LoginLog, error) {
	var logs []*LoginLog
	err := dao.db.Where("status = ?", status).Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// List 获取登录日志列表
func (dao *LoginLogDAO) List(offset, limit int) ([]*LoginLog, error) {
	var logs []*LoginLog
	err := dao.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&logs).Error
	return logs, err
}

// Count 获取登录日志总数
func (dao *LoginLogDAO) Count() (int64, error) {
	var count int64
	err := dao.db.Model(&LoginLog{}).Count(&count).Error
	return count, err
}

// CountByUserID 根据用户ID获取登录日志总数
func (dao *LoginLogDAO) CountByUserID(userID int64) (int64, error) {
	var count int64
	err := dao.db.Model(&LoginLog{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

// Delete 删除登录日志
func (dao *LoginLogDAO) Delete(id int64) error {
	return dao.db.Delete(&LoginLog{}, id).Error
}

// DeleteByUserID 根据用户ID删除登录日志
func (dao *LoginLogDAO) DeleteByUserID(userID int64) error {
	return dao.db.Where("user_id = ?", userID).Delete(&LoginLog{}).Error
}

// DeleteOldLogs 删除指定时间之前的日志
func (dao *LoginLogDAO) DeleteOldLogs(before time.Time) error {
	return dao.db.Where("created_at < ?", before).Delete(&LoginLog{}).Error
}
