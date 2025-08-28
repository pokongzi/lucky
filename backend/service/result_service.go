package service

import (
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
	err = query.Preload("Game").Order("draw_date DESC, period DESC").
		Offset(offset).Limit(pageSize).Find(&results).Error

	return results, total, err
}

// GetDrawResultByPeriod 根据期号获取开奖结果
func GetDrawResultByPeriod(db *gorm.DB, gameCode, period string) (*model.DrawResult, error) {
	var result model.DrawResult

	err := db.Preload("Game").
		Joins("JOIN lottery_games ON draw_results.game_id = lottery_games.id").
		Where("lottery_games.game_code = ? AND draw_results.period = ?", gameCode, period).
		First(&result).Error

	return &result, err
}

// GetLatestDrawResult 获取最新开奖结果
func GetLatestDrawResult(db *gorm.DB, gameCode string) (*model.DrawResult, error) {
	var result model.DrawResult

	err := db.Preload("Game").
		Joins("JOIN lottery_games ON draw_results.game_id = lottery_games.id").
		Where("lottery_games.game_code = ? AND lottery_games.is_active = ?", gameCode, true).
		Order("draw_date DESC, period DESC").
		First(&result).Error

	return &result, err
}
