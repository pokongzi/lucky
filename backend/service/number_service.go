package service

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"lucky/model"

	"gorm.io/gorm"
)

// RandomNumberResult 随机号码结果
type RandomNumberResult struct {
	RedBalls  model.NumberArray `json:"red_balls"`
	BlueBalls model.NumberArray `json:"blue_balls"`
}

// GenerateRandomNumbers 生成随机号码
func GenerateRandomNumbers(game *model.LotteryGame, count int) ([]RandomNumberResult, error) {
	results := make([]RandomNumberResult, count)

	for i := 0; i < count; i++ {
		redBalls := generateRandomBalls(1, game.RedBallCount, game.RedSelectCount)
		blueBalls := generateRandomBalls(1, game.BlueBallCount, game.BlueSelectCount)

		results[i] = RandomNumberResult{
			RedBalls:  redBalls,
			BlueBalls: blueBalls,
		}
	}

	return results, nil
}

// generateRandomBalls 生成随机球号
func generateRandomBalls(min, max, count int) model.NumberArray {
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(1000)))

	balls := make([]int, 0, count)
	usedNumbers := make(map[int]bool)

	for len(balls) < count {
		num := rand.Intn(max-min+1) + min
		if !usedNumbers[num] {
			balls = append(balls, num)
			usedNumbers[num] = true
		}
	}

	sort.Ints(balls)
	return model.NumberArray(balls)
}

// ValidateNumbers 验证号码
func ValidateNumbers(game *model.LotteryGame, redBalls, blueBalls model.NumberArray) error {
	// 验证红球数量
	if len(redBalls) != game.RedSelectCount {
		return errors.New(fmt.Sprintf("红球数量不正确，需要%d个", game.RedSelectCount))
	}

	// 验证蓝球数量
	if len(blueBalls) != game.BlueSelectCount {
		return errors.New(fmt.Sprintf("蓝球数量不正确，需要%d个", game.BlueSelectCount))
	}

	// 验证红球范围和唯一性
	usedRed := make(map[int]bool)
	for _, ball := range redBalls {
		if ball < 1 || ball > game.RedBallCount {
			return errors.New(fmt.Sprintf("红球号码超出范围(1-%d)", game.RedBallCount))
		}
		if usedRed[ball] {
			return errors.New("红球号码重复")
		}
		usedRed[ball] = true
	}

	// 验证蓝球范围和唯一性
	usedBlue := make(map[int]bool)
	for _, ball := range blueBalls {
		if ball < 1 || ball > game.BlueBallCount {
			return errors.New(fmt.Sprintf("蓝球号码超出范围(1-%d)", game.BlueBallCount))
		}
		if usedBlue[ball] {
			return errors.New("蓝球号码重复")
		}
		usedBlue[ball] = true
	}

	return nil
}

// SaveUserNumber 保存用户号码
func SaveUserNumber(db *gorm.DB, userID uint64, gameID uint64, redBalls, blueBalls model.NumberArray, nickname, source string) (*model.UserNumber, error) {
	// 获取游戏信息
	var game model.LotteryGame
	if err := db.First(&game, gameID).Error; err != nil {
		return nil, fmt.Errorf("游戏不存在")
	}

	// 验证号码
	if err := ValidateNumbers(&game, redBalls, blueBalls); err != nil {
		return nil, err
	}

	// 创建用户号码记录
	userNumber := model.UserNumber{
		UserID:    userID,
		GameID:    uint64(gameID),
		RedBalls:  redBalls,
		BlueBalls: blueBalls,
		Nickname:  nickname,
		Source:    source,
		IsActive:  true,
	}

	if err := db.Create(&userNumber).Error; err != nil {
		return nil, err
	}

	return &userNumber, nil
}

// SaveUserNumbers 保存用户号码（批量）
func SaveUserNumbers(db *gorm.DB, userID uint64, gameID uint64, redBalls, blueBalls model.NumberArray) error {
	// 获取游戏信息
	var game model.LotteryGame
	if err := db.First(&game, gameID).Error; err != nil {
		return fmt.Errorf("游戏不存在")
	}

	// 验证号码
	if err := ValidateNumbers(&game, redBalls, blueBalls); err != nil {
		return err
	}

	// 创建用户号码记录
	userNumber := model.UserNumber{
		UserID:    userID,
		GameID:    uint64(gameID),
		RedBalls:  redBalls,
		BlueBalls: blueBalls,
		IsActive:  true,
	}

	return db.Create(&userNumber).Error
}

// UpdateUserNumber 更新用户号码
func UpdateUserNumber(db *gorm.DB, userID uint64, numberID uint64, nickname string, isActive *bool) error {
	// 查找用户号码
	var userNumber model.UserNumber
	if err := db.Where("id = ? AND user_id = ?", numberID, userID).First(&userNumber).Error; err != nil {
		return fmt.Errorf("号码不存在或不属于该用户")
	}

	// 更新字段
	updates := make(map[string]interface{})

	if nickname != "" {
		updates["nickname"] = nickname
	}

	if isActive != nil {
		updates["is_active"] = *isActive
	}

	// 如果没有需要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 执行更新
	return db.Model(&userNumber).Updates(updates).Error
}

// GetUserNumbers 获取用户号码列表
func GetUserNumbers(db *gorm.DB, userID uint64, gameCode string, page, pageSize int) ([]model.UserNumber, int64, error) {
	var numbers []model.UserNumber
	var total int64

	query := db.Model(&model.UserNumber{}).Where("user_id = ?", userID)

	if gameCode != "" {
		// 通过gameCode条件查询
		query = query.Joins("JOIN lottery_games ON user_numbers.game_id = lottery_games.id").
			Where("lottery_games.name = ?", gameCode)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Preload("Game").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&numbers).Error

	return numbers, total, err
}

// DeleteUserNumber 删除用户号码
func DeleteUserNumber(db *gorm.DB, userID uint64, numberID uint64) error {
	result := db.Where("id = ? AND user_id = ?", numberID, userID).Delete(&model.UserNumber{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("号码不存在或不属于该用户")
	}
	return nil
}

// CheckWinningNumbers 检查中奖号码
func CheckWinningNumbers(db *gorm.DB, userNumberID uint, drawResultID uint) (*model.UserDraw, error) {
	// 获取用户号码
	var userNumber model.UserNumber
	if err := db.Preload("Game").First(&userNumber, userNumberID).Error; err != nil {
		return nil, fmt.Errorf("用户号码不存在")
	}

	// 获取开奖结果
	var drawResult model.DrawResult
	if err := db.First(&drawResult, drawResultID).Error; err != nil {
		return nil, fmt.Errorf("开奖结果不存在")
	}

	// 计算中奖情况
	redMatches := countMatches(userNumber.RedBalls, drawResult.RedBalls)
	blueMatches := countMatches(userNumber.BlueBalls, drawResult.BlueBalls)

	// 判断奖级
	prizeLevel := determinePrizeLevel(userNumber.Game.GameName, redMatches, blueMatches)

	// 创建中奖记录
	userDraw := &model.UserDraw{
		UserNumberID: userNumberID,
		DrawResultID: drawResultID,
		PrizeLevel:   prizeLevel,
		IsWinning:    prizeLevel > 0,
		IsActive:     true,
	}

	if err := db.Create(userDraw).Error; err != nil {
		return nil, err
	}

	return userDraw, nil
}

// countMatches 计算匹配数量
func countMatches(userBalls, drawBalls model.NumberArray) int {
	matches := 0
	for _, userBall := range userBalls {
		for _, drawBall := range drawBalls {
			if userBall == drawBall {
				matches++
				break
			}
		}
	}
	return matches
}

// determinePrizeLevel 判断奖级
func determinePrizeLevel(gameName string, redMatches, blueMatches int) int {
	switch gameName {
	case "双色球":
		return determineSSQPrizeLevel(redMatches, blueMatches)
	case "大乐透":
		return determineDLTPrizeLevel(redMatches, blueMatches)
	default:
		return 0
	}
}

// determineSSQPrizeLevel 判断双色球奖级
func determineSSQPrizeLevel(redMatches, blueMatches int) int {
	if redMatches == 6 && blueMatches == 1 {
		return 1 // 一等奖
	} else if redMatches == 6 && blueMatches == 0 {
		return 2 // 二等奖
	} else if redMatches == 5 && blueMatches == 1 {
		return 3 // 三等奖
	} else if (redMatches == 5 && blueMatches == 0) || (redMatches == 4 && blueMatches == 1) {
		return 4 // 四等奖
	} else if (redMatches == 4 && blueMatches == 0) || (redMatches == 3 && blueMatches == 1) {
		return 5 // 五等奖
	} else if (redMatches == 2 && blueMatches == 1) || (redMatches == 1 && blueMatches == 1) || (redMatches == 0 && blueMatches == 1) {
		return 6 // 六等奖
	}
	return 0 // 未中奖
}

// determineDLTPrizeLevel 判断大乐透奖级
func determineDLTPrizeLevel(redMatches, blueMatches int) int {
	if redMatches == 5 && blueMatches == 2 {
		return 1 // 一等奖
	} else if redMatches == 5 && blueMatches == 1 {
		return 2 // 二等奖
	} else if redMatches == 5 && blueMatches == 0 {
		return 3 // 三等奖
	} else if redMatches == 4 && blueMatches == 2 {
		return 4 // 四等奖
	} else if (redMatches == 4 && blueMatches == 1) || (redMatches == 3 && blueMatches == 2) {
		return 5 // 五等奖
	} else if (redMatches == 4 && blueMatches == 0) || (redMatches == 3 && blueMatches == 1) || (redMatches == 2 && blueMatches == 2) {
		return 6 // 六等奖
	} else if (redMatches == 3 && blueMatches == 0) || (redMatches == 1 && blueMatches == 2) || (redMatches == 2 && blueMatches == 1) || (redMatches == 0 && blueMatches == 2) {
		return 7 // 七等奖
	} else if (redMatches == 1 && blueMatches == 1) || (redMatches == 0 && blueMatches == 1) {
		return 8 // 八等奖
	}
	return 0 // 未中奖
}

// GetUserDraws 获取用户中奖记录
func GetUserDraws(db *gorm.DB, userID uint64, page, pageSize int) ([]model.UserDraw, int64, error) {
	var draws []model.UserDraw
	var total int64

	// 获取总数
	err := db.Model(&model.UserDraw{}).
		Joins("JOIN user_numbers ON user_draws.user_number_id = user_numbers.id").
		Where("user_numbers.user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = db.Preload("UserNumber.Game").
		Preload("DrawResult").
		Joins("JOIN user_numbers ON user_draws.user_number_id = user_numbers.id").
		Where("user_numbers.user_id = ?", userID).
		Order("user_draws.created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&draws).Error

	return draws, total, err
}
