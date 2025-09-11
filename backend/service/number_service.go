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

// RandomNumberResult ���������
type RandomNumberResult struct {
	RedBalls  model.NumberArray `json:"red_balls"`
	BlueBalls model.NumberArray `json:"blue_balls"`
}

// GenerateRandomNumbers �����������
func GenerateRandomNumbers(game *model.LotteryGame, count int) ([]RandomNumberResult, error) {
	results := make([]RandomNumberResult, count)

	for i := 0; i < count; i++ {
		redBalls := generateRandomBalls(1, game.RedRange, game.RedCount)
		blueBalls := generateRandomBalls(1, game.BlueRange, game.BlueCount)

		results[i] = RandomNumberResult{
			RedBalls:  redBalls,
			BlueBalls: blueBalls,
		}
	}

	return results, nil
}

// generateRandomBalls ����������
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

// ValidateNumbers ��֤����
func ValidateNumbers(game *model.LotteryGame, redBalls, blueBalls model.NumberArray) error {
	// ��֤��������
	if len(redBalls) != game.RedCount {
		return errors.New(fmt.Sprintf("��������������Ҫ%d��", game.RedCount))
	}

	// ��֤��������
	if len(blueBalls) != game.BlueCount {
		return errors.New(fmt.Sprintf("��������������Ҫ%d��", game.BlueCount))
	}

	// ��֤����Χ��Ψһ��
	usedRed := make(map[int]bool)
	for _, ball := range redBalls {
		if ball < 1 || ball > game.RedRange {
			return errors.New(fmt.Sprintf("������볬����Χ(1-%d)", game.RedRange))
		}
		if usedRed[ball] {
			return errors.New("��������ظ�")
		}
		usedRed[ball] = true
	}

	// ��֤����Χ��Ψһ��
	usedBlue := make(map[int]bool)
	for _, ball := range blueBalls {
		if ball < 1 || ball > game.BlueRange {
			return errors.New(fmt.Sprintf("������볬����Χ(1-%d)", game.BlueRange))
		}
		if usedBlue[ball] {
			return errors.New("��������ظ�")
		}
		usedBlue[ball] = true
	}

	return nil
}

// SaveUserNumbers �����û�����
func SaveUserNumbers(db *gorm.DB, userID uint64, gameID uint, redBalls, blueBalls model.NumberArray) error {
	// ������Ϸ��Ϣ
	var game model.LotteryGame
	if err := db.First(&game, gameID).Error; err != nil {
		return fmt.Errorf("��Ϸ������")
	}

	// ��֤����
	if err := ValidateNumbers(&game, redBalls, blueBalls); err != nil {
		return err
	}

	// �����û������¼
	userNumber := model.UserNumber{
		UserID:    userID,
		GameID:    uint64(gameID),
		RedBalls:  redBalls,
		BlueBalls: blueBalls,
		IsActive:  true,
	}

	return db.Create(&userNumber).Error
}

// GetUserNumbers ��ȡ�û������б�
func GetUserNumbers(db *gorm.DB, userID uint64, gameCode string, page, pageSize int) ([]model.UserNumber, int64, error) {
	var numbers []model.UserNumber
	var total int64

	query := db.Model(&model.UserNumber{}).Where("user_id = ?", userID)

	if gameCode != "" {
		// ͨ��gameCode������ѯ
		query = query.Joins("JOIN lottery_games ON user_numbers.game_id = lottery_games.id").
			Where("lottery_games.name = ?", gameCode)
	}

	// ��ȡ����
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// ��ҳ��ѯ
	offset := (page - 1) * pageSize
	err = query.Preload("Game").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&numbers).Error

	return numbers, total, err
}

// CheckWinningNumbers ����н�����
func CheckWinningNumbers(db *gorm.DB, userNumberID uint, drawResultID uint) (*model.UserDraw, error) {
	// ��ȡ�û�����
	var userNumber model.UserNumber
	if err := db.Preload("Game").First(&userNumber, userNumberID).Error; err != nil {
		return nil, fmt.Errorf("�û����벻����")
	}

	// ��ȡ�������
	var drawResult model.DrawResult
	if err := db.First(&drawResult, drawResultID).Error; err != nil {
		return nil, fmt.Errorf("�������������")
	}

	// ����н����
	redMatches := countMatches(userNumber.RedBalls, drawResult.RedBalls)
	blueMatches := countMatches(userNumber.BlueBalls, drawResult.BlueBalls)

	// �жϽ���
	prizeLevel := determinePrizeLevel(userNumber.Game.Name, redMatches, blueMatches)

	// �����н���¼
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

// countMatches ����ƥ������
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

// determinePrizeLevel �жϽ���
func determinePrizeLevel(gameName string, redMatches, blueMatches int) int {
	switch gameName {
	case "˫ɫ��":
		return determineSSQPrizeLevel(redMatches, blueMatches)
	case "����͸":
		return determineDLTPrizeLevel(redMatches, blueMatches)
	default:
		return 0
	}
}

// determineSSQPrizeLevel �ж�˫ɫ�򽱼�
func determineSSQPrizeLevel(redMatches, blueMatches int) int {
	if redMatches == 6 && blueMatches == 1 {
		return 1 // һ�Ƚ�
	} else if redMatches == 6 && blueMatches == 0 {
		return 2 // ���Ƚ�
	} else if redMatches == 5 && blueMatches == 1 {
		return 3 // ���Ƚ�
	} else if (redMatches == 5 && blueMatches == 0) || (redMatches == 4 && blueMatches == 1) {
		return 4 // �ĵȽ�
	} else if (redMatches == 4 && blueMatches == 0) || (redMatches == 3 && blueMatches == 1) {
		return 5 // ��Ƚ�
	} else if (redMatches == 2 && blueMatches == 1) || (redMatches == 1 && blueMatches == 1) || (redMatches == 0 && blueMatches == 1) {
		return 6 // ���Ƚ�
	}
	return 0 // δ�н�
}

// determineDLTPrizeLevel �жϴ���͸����
func determineDLTPrizeLevel(redMatches, blueMatches int) int {
	if redMatches == 5 && blueMatches == 2 {
		return 1 // һ�Ƚ�
	} else if redMatches == 5 && blueMatches == 1 {
		return 2 // ���Ƚ�
	} else if redMatches == 5 && blueMatches == 0 {
		return 3 // ���Ƚ�
	} else if redMatches == 4 && blueMatches == 2 {
		return 4 // �ĵȽ�
	} else if (redMatches == 4 && blueMatches == 1) || (redMatches == 3 && blueMatches == 2) {
		return 5 // ��Ƚ�
	} else if (redMatches == 4 && blueMatches == 0) || (redMatches == 3 && blueMatches == 1) || (redMatches == 2 && blueMatches == 2) {
		return 6 // ���Ƚ�
	} else if (redMatches == 3 && blueMatches == 0) || (redMatches == 1 && blueMatches == 2) || (redMatches == 2 && blueMatches == 1) || (redMatches == 0 && blueMatches == 2) {
		return 7 // �ߵȽ�
	} else if (redMatches == 1 && blueMatches == 1) || (redMatches == 0 && blueMatches == 1) {
		return 8 // �˵Ƚ�
	}
	return 0 // δ�н�
}

// GetUserDraws ��ȡ�û��н���¼
func GetUserDraws(db *gorm.DB, userID uint64, page, pageSize int) ([]model.UserDraw, int64, error) {
	var draws []model.UserDraw
	var total int64

	// ��ȡ����
	err := db.Model(&model.UserDraw{}).
		Joins("JOIN user_numbers ON user_draws.user_number_id = user_numbers.id").
		Where("user_numbers.user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// ��ҳ��ѯ
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
