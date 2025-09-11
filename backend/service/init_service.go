package service

import (
"log"

"lucky/common/mysql"
"lucky/model"

"gorm.io/gorm"
)

// InitDatabase 初始化数据库
func InitDatabase() error {
db := mysql.DB

// 自动迁移数据库表
if err := db.AutoMigrate(
&model.LotteryGame{},
&model.DrawResult{},
&model.User{},
&model.UserNumber{},
&model.UserDraw{},
); err != nil {
return err
}

// 初始化彩票游戏数据
if err := initLotteryGames(db); err != nil {
return err
}

return nil
}

// initLotteryGames 初始化彩票游戏数据
func initLotteryGames(db *gorm.DB) error {
games := []model.LotteryGame{
{
Name:        "双色球",
Description: "双色球",
RedCount:    6,
BlueCount:   1,
RedRange:    33,
BlueRange:   16,
IsActive:    true,
},
{
Name:        "大乐透",
Description: "大乐透",
RedCount:    5,
BlueCount:   2,
RedRange:    35,
BlueRange:   12,
IsActive:    true,
},
}

for _, game := range games {
// 检查是否已存在
var existingGame model.LotteryGame
err := db.Where("name = ?", game.Name).First(&existingGame).Error
if err == gorm.ErrRecordNotFound {
// 不存在则创建
if err := db.Create(&game).Error; err != nil {
log.Printf("创建游戏 %s 失败: %v", game.Name, err)
return err
}
log.Printf("创建游戏: %s", game.Name)
} else if err != nil {
return err
} else {
log.Printf("游戏已存在: %s", game.Name)
}
}

return nil
}
