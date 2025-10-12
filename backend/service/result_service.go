package service

import (
	"sort"

	"lucky/model"

	"gorm.io/gorm"
)

// GetDrawResults 获取开奖结果列表
func GetDrawResults(db *gorm.DB, gameCode string, page, pageSize int) ([]model.DrawResult, int64, error) {
	var results []model.DrawResult
	var total int64

	// 通过gameCode关联查询
	query := db.Model(&model.DrawResult{}).
		Joins("JOIN lottery_games ON draw_results.game_id = lottery_games.id").
		Where("lottery_games.game_code = ? AND lottery_games.is_active = ?", gameCode, true)

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.Order("draw_date DESC, period DESC").
		Offset(offset).Limit(pageSize).Find(&results).Error

	return results, total, err
}

// GetDrawResultByPeriod 根据期号获取开奖结果
func GetDrawResultByPeriod(db *gorm.DB, gameCode, period string) (*model.DrawResult, error) {
	var result model.DrawResult

	err := db.
		Joins("JOIN lottery_games ON draw_results.game_id = lottery_games.id").
		Where("lottery_games.game_code = ? AND draw_results.period = ?", gameCode, period).
		First(&result).Error

	return &result, err
}

// GetLatestDrawResult 获取最新开奖结果
func GetLatestDrawResult(db *gorm.DB, gameCode string) (*model.DrawResult, error) {
	var result model.DrawResult

	err := db.
		Joins("JOIN lottery_games ON draw_results.game_id = lottery_games.id").
		Where("lottery_games.game_code = ? AND lottery_games.is_active = ?", gameCode, true).
		Order("draw_date DESC, period DESC").
		First(&result).Error

	return &result, err
}

// GetLatestDrawResults 获取最新的N期开奖结果
func GetLatestDrawResults(db *gorm.DB, gameCode string, limit int) ([]model.DrawResult, error) {
	var results []model.DrawResult

	err := db.
		Joins("JOIN lottery_games ON draw_results.game_id = lottery_games.id").
		Where("lottery_games.game_code = ? AND lottery_games.is_active = ?", gameCode, true).
		Order("draw_date DESC, period DESC").
		Limit(limit).
		Find(&results).Error

	return results, err
}

// NumberFrequency 号码频率结构
type NumberFrequency struct {
	Number    int `json:"number"`
	Frequency int `json:"frequency"`
}

// GetNumberDistribution 获取号码分布数据
func GetNumberDistribution(db *gorm.DB, gameCode string, periodCount int) (map[string][]NumberFrequency, error) {
	// 1. 获取指定期数的开奖结果
	var results []model.DrawResult
	query := db.Model(&model.DrawResult{}).
		Joins("JOIN lottery_games ON draw_results.game_id = lottery_games.id").
		Where("lottery_games.game_code = ? AND lottery_games.is_active = ?", gameCode, true).
		Order("draw_date DESC, period DESC").
		Limit(periodCount)

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	// 2. 计算号码频率
	redFreq := make(map[int]int)
	blueFreq := make(map[int]int)

	for _, result := range results {
		// 统计红球频率
		for _, num := range result.RedBalls {
			redFreq[num]++
		}
		// 统计蓝球频率
		for _, num := range result.BlueBalls {
			blueFreq[num]++
		}
	}

	// 3. 转换为返回格式
	redDistribution := []NumberFrequency{}
	blueDistribution := []NumberFrequency{}

	// 处理红球
	for num, freq := range redFreq {
		redDistribution = append(redDistribution, NumberFrequency{
			Number:    num,
			Frequency: freq,
		})
	}

	// 处理蓝球
	for num, freq := range blueFreq {
		blueDistribution = append(blueDistribution, NumberFrequency{
			Number:    num,
			Frequency: freq,
		})
	}

	// 4. 排序（按号码升序）
	sort.Slice(redDistribution, func(i, j int) bool {
		return redDistribution[i].Number < redDistribution[j].Number
	})
	sort.Slice(blueDistribution, func(i, j int) bool {
		return blueDistribution[i].Number < blueDistribution[j].Number
	})

	// 5. 补全缺失的号码（频率为0）
	var redMax, blueMax int
	if gameCode == "ssq" {
		redMax = 33
		blueMax = 16
	} else if gameCode == "dlt" {
		redMax = 35
		blueMax = 12
	}

	// 补全红球
	completeDistribution := []NumberFrequency{}
	for i := 1; i <= redMax; i++ {
		found := false
		for _, item := range redDistribution {
			if item.Number == i {
				completeDistribution = append(completeDistribution, item)
				found = true
				break
			}
		}
		if !found {
			completeDistribution = append(completeDistribution, NumberFrequency{
				Number:    i,
				Frequency: 0,
			})
		}
	}
	redDistribution = completeDistribution

	// 补全蓝球
	completeDistribution = []NumberFrequency{}
	for i := 1; i <= blueMax; i++ {
		found := false
		for _, item := range blueDistribution {
			if item.Number == i {
				completeDistribution = append(completeDistribution, item)
				found = true
				break
			}
		}
		if !found {
			completeDistribution = append(completeDistribution, NumberFrequency{
				Number:    i,
				Frequency: 0,
			})
		}
	}
	blueDistribution = completeDistribution

	return map[string][]NumberFrequency{
		"red":  redDistribution,
		"blue": blueDistribution,
	}, nil
}
