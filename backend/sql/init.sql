-- 彩票号码生成器数据库初始化脚本
-- MySQL 8.0+

-- 创建数据库
CREATE DATABASE IF NOT EXISTS lottery_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE lottery_db;

-- 1. 用户表 (更新支持JWT)
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `open_id` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '微信OpenID',
  `nickname` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户昵称',
  `avatar_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '头像URL',
  `status` int NOT NULL DEFAULT '1' COMMENT '用户状态(1:正常 0:禁用)',
  `token_version` int NOT NULL DEFAULT '1' COMMENT 'token版本号',
  `last_login_at` datetime(3) DEFAULT NULL COMMENT '最后登录时间',
  `last_login_ip` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '最后登录IP',
  `login_count` int NOT NULL DEFAULT '0' COMMENT '登录次数',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_open_id` (`open_id`),
  KEY `idx_users_created_at` (`created_at`),
  KEY `idx_users_status_version` (`status`, `token_version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- 2. 刷新Token表
CREATE TABLE `refresh_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Token ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `token` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '刷新token',
  `expires_at` datetime(3) NOT NULL COMMENT '过期时间',
  `is_revoked` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已撤销',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户代理',
  `client_ip` varchar(45) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '客户端IP',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_refresh_tokens_token` (`token`),
  KEY `idx_refresh_tokens_user_id` (`user_id`),
  KEY `idx_refresh_tokens_expires_at` (`expires_at`),
  KEY `idx_refresh_tokens_user_expires` (`user_id`, `expires_at` DESC, `is_revoked`),
  CONSTRAINT `fk_refresh_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='刷新Token表';

-- 3. 登录日志表
CREATE TABLE `login_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `login_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录类型(wechat,refresh)',
  `client_ip` varchar(45) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户端IP',
  `user_agent` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户代理',
  `status` int NOT NULL COMMENT '登录状态(1:成功 0:失败)',
  `error_msg` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '错误信息',
  `login_at` datetime(3) NOT NULL COMMENT '登录时间',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_login_logs_user_id` (`user_id`),
  KEY `idx_login_logs_login_at` (`login_at`),
  KEY `idx_login_logs_user_time` (`user_id`, `login_at` DESC),
  CONSTRAINT `fk_login_logs_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录日志表';

-- 4. 彩票游戏表
CREATE TABLE `lottery_games` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '游戏ID',
  `game_code` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '游戏代码',
  `game_name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '游戏名称',
  `red_ball_count` int NOT NULL COMMENT '红球总数',
  `blue_ball_count` int NOT NULL COMMENT '蓝球总数',
  `red_select_count` int NOT NULL COMMENT '红球选择数',
  `blue_select_count` int NOT NULL COMMENT '蓝球选择数',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_lottery_games_game_code` (`game_code`),
  KEY `idx_lottery_games_is_active` (`is_active`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='彩票游戏表';

-- 3. 用户号码表
CREATE TABLE `user_numbers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户号码ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `game_id` bigint unsigned NOT NULL COMMENT '游戏ID',
  `red_balls` json NOT NULL COMMENT '红球号码JSON数组',
  `blue_balls` json NOT NULL COMMENT '蓝球号码JSON数组',
  `nickname` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '号码昵称',
  `source` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'manual' COMMENT '来源类型',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_numbers_user_id` (`user_id`),
  KEY `idx_user_numbers_game_id` (`game_id`),
  KEY `idx_user_numbers_is_active` (`is_active`),
  KEY `idx_user_numbers_created_at` (`created_at`),
  KEY `idx_user_numbers_user_game` (`user_id`, `game_id`, `is_active`),
  CONSTRAINT `fk_user_numbers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_numbers_game` FOREIGN KEY (`game_id`) REFERENCES `lottery_games` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户号码表';

-- 4. 开奖结果表
CREATE TABLE `draw_results` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '开奖ID',
  `game_id` bigint unsigned NOT NULL COMMENT '游戏ID',
  `period` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '期号',
  `draw_date` datetime(3) NOT NULL COMMENT '开奖日期',
  `red_balls` json NOT NULL COMMENT '红球号码JSON数组',
  `blue_balls` json NOT NULL COMMENT '蓝球号码JSON数组',
  `sales_amount` bigint NOT NULL DEFAULT '0' COMMENT '销售额(分)',
  `prize_pool` bigint NOT NULL DEFAULT '0' COMMENT '奖池金额(分)',
  `first_prize` int NOT NULL DEFAULT '0' COMMENT '一等奖注数',
  `first_amount` bigint NOT NULL DEFAULT '0' COMMENT '一等奖单注奖金(分)',
  `second_prize` int NOT NULL DEFAULT '0' COMMENT '二等奖注数',
  `second_amount` bigint NOT NULL DEFAULT '0' COMMENT '二等奖单注奖金(分)',
  `third_prize` int NOT NULL DEFAULT '0' COMMENT '三等奖注数',
  `third_amount` bigint NOT NULL DEFAULT '0' COMMENT '三等奖单注奖金(分)',
  `fourth_prize` int NOT NULL DEFAULT '0' COMMENT '四等奖注数',
  `fourth_amount` bigint NOT NULL DEFAULT '0' COMMENT '四等奖单注奖金(分)',
  `fifth_prize` int NOT NULL DEFAULT '0' COMMENT '五等奖注数',
  `fifth_amount` bigint NOT NULL DEFAULT '0' COMMENT '五等奖单注奖金(分)',
  `sixth_prize` int NOT NULL DEFAULT '0' COMMENT '六等奖注数',
  `sixth_amount` bigint NOT NULL DEFAULT '0' COMMENT '六等奖单注奖金(分)',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_draw_results_game_period` (`game_id`, `period`),
  KEY `idx_draw_results_game_id` (`game_id`),
  KEY `idx_draw_results_period` (`period`),
  KEY `idx_draw_results_draw_date` (`draw_date`),
  KEY `idx_draw_results_game_date` (`game_id`, `draw_date`),
  CONSTRAINT `fk_draw_results_game` FOREIGN KEY (`game_id`) REFERENCES `lottery_games` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='开奖结果表';

-- 初始化游戏数据
INSERT INTO `lottery_games` (`game_code`, `game_name`, `red_ball_count`, `blue_ball_count`, `red_select_count`, `blue_select_count`, `is_active`) VALUES
('ssq', '双色球', 33, 16, 6, 1, 1),
('dlt', '大乐透', 35, 12, 5, 2, 1);

-- 创建视图：用户号码详情视图
CREATE VIEW `v_user_number_details` AS
SELECT 
    un.id AS user_number_id,
    un.user_id,
    un.nickname AS number_nickname,
    un.source,
    un.is_active,
    un.created_at AS saved_at,
    lg.game_code,
    lg.game_name,
    un.red_balls,
    un.blue_balls,
    u.nickname AS user_nickname
FROM user_numbers un
INNER JOIN users u ON un.user_id = u.id
INNER JOIN lottery_games lg ON un.game_id = lg.id;

-- 创建索引优化查询性能
CREATE INDEX idx_draw_results_composite ON draw_results(game_id, period DESC, draw_date DESC);
CREATE INDEX idx_user_numbers_composite ON user_numbers(user_id, game_id, is_active, created_at DESC); 