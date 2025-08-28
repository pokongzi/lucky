package migration

import (
	"log"
	"os"

	"lucky/model"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	log.Printf("开始数据库迁移 (环境: %s)...", env)

	// 生产环境禁用自动迁移
	if env == "production" {
		log.Println("生产环境禁用自动迁移，请使用手动迁移脚本")
		return nil
	}

	// 开发/测试环境允许自动迁移
	if env == "development" || env == "test" {
		log.Println("开发/测试环境：执行自动迁移...")

		// 自动迁移表结构
		err := db.AutoMigrate(
			&model.User{},
			&model.RefreshToken{}, // 新增
			&model.LoginLog{},     // 新增
			&model.LotteryGame{},
			&model.UserNumber{},
			&model.DrawResult{},
		)
		if err != nil {
			log.Printf("自动迁移失败: %v", err)
			return err
		}

		log.Println("表结构迁移完成")

		// 创建自定义索引
		err = createCustomIndexes(db)
		if err != nil {
			log.Printf("创建自定义索引失败: %v", err)
			return err
		}

		log.Println("数据库迁移完成")
	}

	return nil
}

// createCustomIndexes 创建自定义索引
func createCustomIndexes(db *gorm.DB) error {
	// 复合索引 - 用户号码
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_user_numbers_composite 
		ON user_numbers(user_id, game_id, is_active DESC, created_at DESC)
	`).Error; err != nil {
		log.Printf("创建用户号码复合索引失败: %v", err)
		return err
	}

	// 复合索引 - 开奖结果
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_draw_results_composite 
		ON draw_results(game_id, draw_date DESC, period DESC)
	`).Error; err != nil {
		log.Printf("创建开奖结果复合索引失败: %v", err)
		return err
	}

	// JWT相关索引 - 刷新Token
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_expires 
		ON refresh_tokens(user_id, expires_at DESC, is_revoked)
	`).Error; err != nil {
		log.Printf("创建刷新Token复合索引失败: %v", err)
		return err
	}

	// JWT相关索引 - 登录日志
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_login_logs_user_time 
		ON login_logs(user_id, login_at DESC)
	`).Error; err != nil {
		log.Printf("创建登录日志复合索引失败: %v", err)
		return err
	}

	// 用户状态索引
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_status_version 
		ON users(status, token_version)
	`).Error; err != nil {
		log.Printf("创建用户状态索引失败: %v", err)
		return err
	}

	log.Println("自定义索引创建完成")
	return nil
}

// ManualMigrate 手动迁移模式（用于生产环境）
func ManualMigrate(db *gorm.DB) error {
	log.Println("手动迁移模式暂未实现")
	return nil
}
