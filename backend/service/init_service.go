package service

import (
	"lucky/model"

	"gorm.io/gorm"
)

// InitializeData 初始化基础数据
func InitializeData(db *gorm.DB) error {
	// 初始化彩票游戏数据
	if err := initLotteryGames(db); err != nil {
		return err
	}

	return nil
}

// initLotteryGames 初始化彩票游戏数据
func initLotteryGames(db *gorm.DB) error {
	// 检查是否已存在游戏数据
	var count int64
	if err := db.Model(&model.LotteryGame{}).Count(&count).Error; err != nil {
		return err
	}

	// 如果已有数据，则不需要初始化
	if count > 0 {
		return nil
	}

	// 初始化双色球
	ssq := model.LotteryGame{
		GameCode:        "ssq",
		GameName:        "双色球",
		RedBallCount:    33,
		BlueBallCount:   16,
		RedSelectCount:  6,
		BlueSelectCount: 1,
		IsActive:        true,
	}

	// 初始化大乐透
	dlt := model.LotteryGame{
		GameCode:        "dlt",
		GameName:        "大乐透",
		RedBallCount:    35,
		BlueBallCount:   12,
		RedSelectCount:  5,
		BlueSelectCount: 2,
		IsActive:        true,
	}

	// 批量创建
	return db.Create([]*model.LotteryGame{&ssq, &dlt}).Error
}
