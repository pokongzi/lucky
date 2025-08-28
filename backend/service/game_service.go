package service

import (
	"lucky/model"

	"gorm.io/gorm"
)

// GetActiveGames 获取所有启用的游戏
func GetActiveGames(db *gorm.DB) ([]model.LotteryGame, error) {
	var games []model.LotteryGame
	err := db.Where("is_active = ?", true).Find(&games).Error
	return games, err
}

// GetGameByCode 根据游戏代码获取游戏
func GetGameByCode(db *gorm.DB, gameCode string) (*model.LotteryGame, error) {
	var game model.LotteryGame
	err := db.Where("game_code = ? AND is_active = ?", gameCode, true).First(&game).Error
	return &game, err
}

// GetGameByID 根据ID获取游戏
func GetGameByID(db *gorm.DB, gameID uint64) (*model.LotteryGame, error) {
	var game model.LotteryGame
	err := db.First(&game, gameID).Error
	return &game, err
}
