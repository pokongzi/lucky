package service

import (
	"strconv"

	"lucky/model"

	"gorm.io/gorm"
)

// GetOrCreateUser 获取或创建用户
func GetOrCreateUser(db *gorm.DB, openid, nickname, avatar string) (*model.User, error) {
	var user model.User
	err := db.Where("open_id = ?", openid).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = model.User{OpenID: openid, Nickname: nickname, AvatarURL: avatar}
		err = db.Create(&user).Error
	}
	return &user, err
}

// GetUserByID 根据ID获取用户
func GetUserByID(db *gorm.DB, userIDStr string) (*model.User, error) {
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = db.First(&user, userID).Error
	return &user, err
}
