package service

import (
	"lucky/model"

	"gorm.io/gorm"
)

func SaveUserTicket(db *gorm.DB, userID uint64, t int, numbers string) error {
	ticket := model.UserTicket{UserID: userID, Type: t, Numbers: numbers}
	return db.Create(&ticket).Error
}

func ListUserTickets(db *gorm.DB, userID uint64, t int) ([]model.UserTicket, error) {
	var tickets []model.UserTicket
	err := db.Where("user_id = ? AND type = ?", userID, t).Order("created_at desc").Find(&tickets).Error
	return tickets, err
}

func DeleteUserTicket(db *gorm.DB, ticketID, userID uint64) error {
	return db.Where("id = ? AND user_id = ?", ticketID, userID).Delete(&model.UserTicket{}).Error
}
