package model

import (
"time"
"gorm.io/gorm"
)

type User struct {
ID           int64      `gorm:"primaryKey;column:id" json:"id"`
OpenID       string     `gorm:"uniqueIndex;size:64;not null;column:open_id" json:"open_id"`
Nickname     string     `gorm:"size:64;not null;column:nickname" json:"nickname"`
AvatarURL    string     `gorm:"size:255;column:avatar_url" json:"avatar_url"`
Status       int        `gorm:"not null;default:1;column:status;comment:用户状态(1:正常 0:禁用)" json:"status"`
TokenVersion int        `gorm:"not null;default:1;column:token_version;comment:token版本号" json:"token_version"`
LastLoginAt  *time.Time `gorm:"column:last_login_at;comment:最后登录时间" json:"last_login_at"`
LastLoginIP  string     `gorm:"size:45;column:last_login_ip;comment:最后登录IP" json:"last_login_ip"`
LoginCount   int        `gorm:"not null;default:0;column:login_count;comment:登录次数" json:"login_count"`
CreatedAt    time.Time  `gorm:"column:created_at" json:"created_at"`
UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

func (User) TableName() string {
return "users"
}

// UserDAO 用户数据访问对象
type UserDAO struct {
db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
return &UserDAO{db: db}
}

// Create 创建用户
func (dao *UserDAO) Create(user *User) error {
return dao.db.Create(user).Error
}

// GetByID 根据ID获取用户
func (dao *UserDAO) GetByID(id int64) (*User, error) {
var user User
err := dao.db.First(&user, id).Error
if err != nil {
return nil, err
}
return &user, nil
}

// GetByOpenID 根据OpenID获取用户
func (dao *UserDAO) GetByOpenID(openID string) (*User, error) {
var user User
err := dao.db.Where("open_id = ?", openID).First(&user).Error
if err != nil {
return nil, err
}
return &user, nil
}

// Update 更新用户
func (dao *UserDAO) Update(user *User) error {
return dao.db.Save(user).Error
}

// UpdateTokenVersion 更新token版本号
func (dao *UserDAO) UpdateTokenVersion(userID int64, version int) error {
return dao.db.Model(&User{}).Where("id = ?", userID).Update("token_version", version).Error
}

// UpdateLoginInfo 更新登录信息
func (dao *UserDAO) UpdateLoginInfo(userID int64, lastLoginAt time.Time, lastLoginIP string) error {
return dao.db.Model(&User{}).Where("id = ?", userID).Updates(map[string]interface{}{
"last_login_at": lastLoginAt,
"last_login_ip": lastLoginIP,
"login_count":   gorm.Expr("login_count + 1"),
}).Error
}

// Delete 删除用户
func (dao *UserDAO) Delete(id int64) error {
return dao.db.Delete(&User{}, id).Error
}

// List 获取用户列表
func (dao *UserDAO) List(offset, limit int) ([]*User, error) {
var users []*User
err := dao.db.Offset(offset).Limit(limit).Find(&users).Error
return users, err
}

// Count 获取用户总数
func (dao *UserDAO) Count() (int64, error) {
var count int64
err := dao.db.Model(&User{}).Count(&count).Error
return count, err
}
