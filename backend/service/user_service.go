package service

import (
	"lucky/backend/model"

	"gorm.io/gorm"
)

func GetOrCreateUser(db *gorm.DB, openid, nickname, avatar string) (*model.User, error) {
	var user model.User
	err := db.Where("open_id = ?", openid).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		user = model.User{OpenID: openid, Nickname: nickname, AvatarURL: avatar}
		err = db.Create(&user).Error
	}
	return &user, err
}
