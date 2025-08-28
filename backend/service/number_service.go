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
	RedBalls  model.NumberArray `json:"redBalls"`
	BlueBalls model.NumberArray `json:"blueBalls"`
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

// generateRandomBalls 生成指定范围和数量的随机球号
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
		return errors.New(fmt.Sprintf("红球数量错误，需要%d个", game.RedSelectCount))
	}

	// 验证蓝球数量
	if len(blueBalls) != game.BlueSelectCount {
		return errors.New(fmt.Sprintf("蓝球数量错误，需要%d个", game.BlueSelectCount))
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
func SaveUserNumber(db *gorm.DB, userID, gameID uint64, redBalls, blueBalls model.NumberArray, nickname, source string) (*model.UserNumber, error) {
	if source == "" {
		source = "manual"
	}

	userNumber := &model.UserNumber{
		UserID:    userID,
		GameID:    gameID,
		RedBalls:  redBalls,
		BlueBalls: blueBalls,
		Nickname:  nickname,
		Source:    source,
		IsActive:  true,
	}

	err := db.Create(userNumber).Error
	return userNumber, err
}

// GetUserNumbers 获取用户号码列表
func GetUserNumbers(db *gorm.DB, userID uint64, gameCode string, page, pageSize int) ([]model.UserNumber, int64, error) {
	var numbers []model.UserNumber
	var total int64

	query := db.Model(&model.UserNumber{}).Where("user_id = ?", userID)

	if gameCode != "" {
		// 通过gameCode关联查询
		query = query.Joins("JOIN lottery_games ON user_numbers.game_id = lottery_games.id").
			Where("lottery_games.game_code = ?", gameCode)
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

// UpdateUserNumber 更新用户号码
func UpdateUserNumber(db *gorm.DB, userID, numberID uint64, nickname string, isActive *bool) error {
	updates := make(map[string]interface{})

	if nickname != "" {
		updates["nickname"] = nickname
	}

	if isActive != nil {
		updates["is_active"] = *isActive
	}

	if len(updates) == 0 {
		return errors.New("没有要更新的字段")
	}

	result := db.Model(&model.UserNumber{}).
		Where("id = ? AND user_id = ?", numberID, userID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("记录不存在或权限不足")
	}

	return nil
}

// DeleteUserNumber 删除用户号码
func DeleteUserNumber(db *gorm.DB, userID, numberID uint64) error {
	result := db.Where("id = ? AND user_id = ?", numberID, userID).Delete(&model.UserNumber{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("记录不存在或权限不足")
	}

	return nil
}
