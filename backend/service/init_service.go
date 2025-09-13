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
		Name:        "双色球",
		Description: "中国福利彩票双色球游戏",
		RedCount:    6,
		BlueCount:   1,
		RedRange:    33,
		BlueRange:   16,
		IsActive:    true,
	}

	// 初始化大乐透
	dlt := model.LotteryGame{
		Name:        "大乐透",
		Description: "中国体育彩票超级大乐透",
		RedCount:    5,
		BlueCount:   2,
		RedRange:    35,
		BlueRange:   12,
		IsActive:    true,
	}

	// 批量创建
	return db.Create([]*model.LotteryGame{&ssq, &dlt}).Error
}
