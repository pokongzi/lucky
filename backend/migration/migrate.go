package migration

import (
	"database/sql"
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

		// 兼容性处理：尝试移除历史外键以避免列类型修改失败
		// 例如：refresh_tokens.user_id -> users.id 的外键 fk_refresh_tokens_user
		if err := dropLegacyForeignKeys(db); err != nil {
			log.Printf("移除历史外键失败(可忽略继续): %v", err)
		}

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
	// 获取当前数据库名
	dbName, err := getCurrentDatabase(db)
	if err != nil {
		return err
	}

	// 定义需要的索引
	indexes := []struct {
		table     string
		indexName string
		createSQL string
	}{
		{"user_numbers", "idx_user_numbers_composite", "CREATE INDEX idx_user_numbers_composite ON user_numbers(user_id, game_id, is_active, created_at)"},
		{"draw_results", "idx_draw_results_composite", "CREATE INDEX idx_draw_results_composite ON draw_results(game_id, draw_date, period)"},
		{"refresh_tokens", "idx_refresh_tokens_user_expires", "CREATE INDEX idx_refresh_tokens_user_expires ON refresh_tokens(user_id, expires_at, is_active)"},
		{"login_logs", "idx_login_logs_user_time", "CREATE INDEX idx_login_logs_user_time ON login_logs(user_id, created_at)"},
		{"users", "idx_users_status_version", "CREATE INDEX idx_users_status_version ON users(status, token_version)"},
	}

	for _, idx := range indexes {
		exists, err := indexExists(db, dbName, idx.table, idx.indexName)
		if err != nil {
			log.Printf("检查索引是否存在失败: %s.%s %v", idx.table, idx.indexName, err)
			return err
		}
		if exists {
			continue
		}
		if err := db.Exec(idx.createSQL).Error; err != nil {
			log.Printf("创建索引失败: %s on %s, err=%v", idx.indexName, idx.table, err)
			return err
		}
	}

	log.Println("自定义索引创建完成")
	return nil
}

// dropLegacyForeignKeys 尝试移除历史外键，避免自动迁移被外键阻塞
func dropLegacyForeignKeys(db *gorm.DB) error {
	// MySQL: 如外键不存在会报错，这里忽略错误以继续后续迁移
	_ = db.Exec("ALTER TABLE refresh_tokens DROP FOREIGN KEY fk_refresh_tokens_user").Error
	_ = db.Exec("ALTER TABLE login_logs DROP FOREIGN KEY fk_login_logs_user").Error
	_ = db.Exec("ALTER TABLE user_numbers DROP FOREIGN KEY fk_user_numbers_user").Error
	// 可按需在此扩展其它可能遗留的外键名称
	return nil
}

// indexExists 检查索引是否存在
func indexExists(db *gorm.DB, dbName, table, indexName string) (bool, error) {
	const q = `SELECT COUNT(1) FROM INFORMATION_SCHEMA.STATISTICS WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?`
	var cnt int
	if err := db.Raw(q, dbName, table, indexName).Scan(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// getCurrentDatabase 获取当前数据库名
func getCurrentDatabase(db *gorm.DB) (string, error) {
	var name sql.NullString
	if err := db.Raw("SELECT DATABASE()").Scan(&name).Error; err != nil {
		return "", err
	}
	if !name.Valid {
		return "", nil
	}
	return name.String, nil
}

// ManualMigrate 手动迁移模式（用于生产环境）
func ManualMigrate(db *gorm.DB) error {
	log.Println("手动迁移模式暂未实现")
	return nil
}
