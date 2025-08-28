package service

import (
	"log"

	"lucky/model"

	"gorm.io/gorm"
)

// InitializeData 初始化基础数据
func InitializeData(db *gorm.DB) error {
	log.Println("开始初始化基础数据...")

	// 初始化彩票游戏数据
	err := initLotteryGames(db)
	if err != nil {
		return err
	}

	log.Println("基础数据初始化完成")
	return nil
}

// initLotteryGames 初始化彩票游戏数据
func initLotteryGames(db *gorm.DB) error {
	games := []model.LotteryGame{
		{
			GameCode:        "ssq",
			GameName:        "双色球",
			RedBallCount:    33,
			BlueBallCount:   16,
			RedSelectCount:  6,
			BlueSelectCount: 1,
			IsActive:        true,
		},
		{
			GameCode:        "dlt",
			GameName:        "大乐透",
			RedBallCount:    35,
			BlueBallCount:   12,
			RedSelectCount:  5,
			BlueSelectCount: 2,
			IsActive:        true,
		},
	}

	for _, game := range games {
		// 检查是否已存在
		var existingGame model.LotteryGame
		err := db.Where("game_code = ?", game.GameCode).First(&existingGame).Error
		if err == gorm.ErrRecordNotFound {
			// 不存在则创建
			if err := db.Create(&game).Error; err != nil {
				log.Printf("创建游戏 %s 失败: %v", game.GameName, err)
				return err
			}
			log.Printf("创建游戏: %s", game.GameName)
		} else if err != nil {
			return err
		} else {
			log.Printf("游戏已存在: %s", game.GameName)
		}
	}

	return nil
}
