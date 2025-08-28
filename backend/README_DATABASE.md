# 彩票号码生成器 - 数据库设计文档

## 概述

本项目使用MySQL 8.0+作为数据库，采用GORM作为ORM框架。数据库设计遵循关系型数据库的设计原则，确保数据一致性和查询性能。

## 表结构设计

### 1. 用户表 (users)

存储用户基本信息，支持微信小程序登录。

```sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
  `open_id` varchar(64) NOT NULL COMMENT '微信OpenID',
  `nickname` varchar(64) NOT NULL COMMENT '用户昵称',
  `avatar_url` varchar(255) DEFAULT NULL COMMENT '头像URL',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_open_id` (`open_id`)
);
```

**字段说明：**
- `open_id`: 微信OpenID，用于用户身份识别
- `nickname`: 用户昵称，从微信获取
- `avatar_url`: 用户头像URL

### 2. 彩票游戏表 (lottery_games)

存储彩票游戏的基本配置信息，如双色球、大乐透等。

```sql
CREATE TABLE `lottery_games` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '游戏ID',
  `game_code` varchar(32) NOT NULL COMMENT '游戏代码',
  `game_name` varchar(64) NOT NULL COMMENT '游戏名称',
  `red_ball_count` int NOT NULL COMMENT '红球总数',
  `blue_ball_count` int NOT NULL COMMENT '蓝球总数',
  `red_select_count` int NOT NULL COMMENT '红球选择数',
  `blue_select_count` int NOT NULL COMMENT '蓝球选择数',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_lottery_games_game_code` (`game_code`)
);
```

**预置数据：**
- 双色球: game_code='ssq', 红球1-33选6个, 蓝球1-16选1个
- 大乐透: game_code='dlt', 红球1-35选5个, 蓝球1-12选2个

### 3. 号码组合表 (number_combinations)

存储具体的号码组合，避免重复存储相同的号码组合。

```sql
CREATE TABLE `number_combinations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '组合ID',
  `game_id` bigint unsigned NOT NULL COMMENT '游戏ID',
  `red_balls` json NOT NULL COMMENT '红球号码JSON数组',
  `blue_balls` json NOT NULL COMMENT '蓝球号码JSON数组',
  `hash` varchar(64) NOT NULL COMMENT '号码组合哈希值',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_number_combinations_hash` (`hash`),
  CONSTRAINT `fk_number_combinations_game` FOREIGN KEY (`game_id`) REFERENCES `lottery_games` (`id`)
);
```

**字段说明：**
- `red_balls`: JSON数组，如 [1,5,12,18,25,33]
- `blue_balls`: JSON数组，如 [8] (双色球) 或 [3,11] (大乐透)
- `hash`: 号码组合的哈希值，用于快速查重

### 4. 用户号码表 (user_numbers)

存储用户收藏的号码，关联用户和号码组合。

```sql
CREATE TABLE `user_numbers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '用户号码ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `combination_id` bigint unsigned NOT NULL COMMENT '号码组合ID',
  `game_id` bigint unsigned NOT NULL COMMENT '游戏ID',
  `nickname` varchar(128) DEFAULT NULL COMMENT '号码昵称',
  `source` varchar(32) NOT NULL DEFAULT 'manual' COMMENT '来源类型',
  `is_active` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_user_numbers_user_game` (`user_id`, `game_id`, `is_active`),
  CONSTRAINT `fk_user_numbers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_numbers_combination` FOREIGN KEY (`combination_id`) REFERENCES `number_combinations` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_numbers_game` FOREIGN KEY (`game_id`) REFERENCES `lottery_games` (`id`) ON DELETE CASCADE
);
```

**字段说明：**
- `nickname`: 用户给号码起的昵称
- `source`: 来源类型，'manual'(手动输入) 或 'random'(机选)

### 5. 开奖结果表 (draw_results)

存储历史开奖数据。

```sql
CREATE TABLE `draw_results` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '开奖ID',
  `game_id` bigint unsigned NOT NULL COMMENT '游戏ID',
  `period` varchar(32) NOT NULL COMMENT '期号',
  `draw_date` datetime(3) NOT NULL COMMENT '开奖日期',
  `red_balls` json NOT NULL COMMENT '红球号码JSON数组',
  `blue_balls` json NOT NULL COMMENT '蓝球号码JSON数组',
  `sales_amount` bigint NOT NULL DEFAULT '0' COMMENT '销售额(分)',
  `prize_pool` bigint NOT NULL DEFAULT '0' COMMENT '奖池金额(分)',
  `first_prize` int NOT NULL DEFAULT '0' COMMENT '一等奖注数',
  `first_amount` bigint NOT NULL DEFAULT '0' COMMENT '一等奖单注奖金(分)',
  -- ... 其他奖项字段
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_draw_results_game_period` (`game_id`, `period`),
  CONSTRAINT `fk_draw_results_game` FOREIGN KEY (`game_id`) REFERENCES `lottery_games` (`id`)
);
```

**字段说明：**
- `period`: 期号，如 "2023130"
- 奖金相关字段均以"分"为单位存储，避免浮点数精度问题

### 6. 中奖记录表 (winning_records)

存储用户号码的中奖情况。

```sql
CREATE TABLE `winning_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '中奖记录ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `user_number_id` bigint unsigned NOT NULL COMMENT '用户号码ID',
  `draw_result_id` bigint unsigned NOT NULL COMMENT '开奖结果ID',
  `game_id` bigint unsigned NOT NULL COMMENT '游戏ID',
  `period` varchar(32) NOT NULL COMMENT '期号',
  `prize_level` int NOT NULL COMMENT '奖级',
  `red_matches` int NOT NULL DEFAULT '0' COMMENT '红球命中数',
  `blue_matches` int NOT NULL DEFAULT '0' COMMENT '蓝球命中数',
  `prize_amount` bigint NOT NULL DEFAULT '0' COMMENT '奖金金额(分)',
  `is_verified` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已验证',
  `created_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
  `updated_at` datetime(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
  PRIMARY KEY (`id`),
  KEY `idx_winning_records_user_game` (`user_id`, `game_id`),
  CONSTRAINT `fk_winning_records_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_winning_records_user_number` FOREIGN KEY (`user_number_id`) REFERENCES `user_numbers` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_winning_records_draw_result` FOREIGN KEY (`draw_result_id`) REFERENCES `draw_results` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_winning_records_game` FOREIGN KEY (`game_id`) REFERENCES `lottery_games` (`id`) ON DELETE CASCADE
);
```

## Go模型文件

### 1. User (backend/model/user.go)
```go
type User struct {
    ID        uint64 `gorm:"primaryKey"`
    OpenID    string `gorm:"uniqueIndex;size:64;not null"`
    Nickname  string `gorm:"size:64;not null"`
    AvatarURL string `gorm:"size:255"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 2. LotteryGame (backend/model/lottery_game.go)
```go
type LotteryGame struct {
    ID              uint64 `gorm:"primaryKey"`
    GameCode        string `gorm:"uniqueIndex;size:32;not null"`
    GameName        string `gorm:"size:64;not null"`
    RedBallCount    int    `gorm:"not null"`
    BlueBallCount   int    `gorm:"not null"`
    RedSelectCount  int    `gorm:"not null"`
    BlueSelectCount int    `gorm:"not null"`
    IsActive        bool   `gorm:"default:true"`
    CreatedAt       time.Time
    UpdatedAt       time.Time
}
```

### 3. NumberCombination (backend/model/number_combination.go)
```go
type NumberCombination struct {
    ID        uint64      `gorm:"primaryKey"`
    GameID    uint64      `gorm:"not null;index"`
    RedBalls  NumberArray `gorm:"type:json;not null"`
    BlueBalls NumberArray `gorm:"type:json;not null"`
    Hash      string      `gorm:"uniqueIndex;size:64"`
    CreatedAt time.Time
    UpdatedAt time.Time
    
    Game        LotteryGame  `gorm:"foreignKey:GameID"`
    UserNumbers []UserNumber `gorm:"foreignKey:CombinationID"`
}
```

### 4. UserNumber (backend/model/user_number.go)
```go
type UserNumber struct {
    ID            uint64 `gorm:"primaryKey"`
    UserID        uint64 `gorm:"not null;index"`
    CombinationID uint64 `gorm:"not null;index"`
    GameID        uint64 `gorm:"not null;index"`
    Nickname      string `gorm:"size:128"`
    Source        string `gorm:"size:32;default:'manual'"`
    IsActive      bool   `gorm:"default:true"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
    
    User        User               `gorm:"foreignKey:UserID"`
    Combination NumberCombination  `gorm:"foreignKey:CombinationID"`
    Game        LotteryGame        `gorm:"foreignKey:GameID"`
    WinRecords  []WinningRecord    `gorm:"foreignKey:UserNumberID"`
}
```

### 5. DrawResult (backend/model/draw_result.go)
```go
type DrawResult struct {
    ID           uint64      `gorm:"primaryKey"`
    GameID       uint64      `gorm:"not null;index"`
    Period       string      `gorm:"size:32;not null;index"`
    DrawDate     time.Time   `gorm:"not null;index"`
    RedBalls     NumberArray `gorm:"type:json;not null"`
    BlueBalls    NumberArray `gorm:"type:json;not null"`
    SalesAmount  int64       `gorm:"default:0"`
    PrizePool    int64       `gorm:"default:0"`
    // ... 奖项字段
    CreatedAt    time.Time
    UpdatedAt    time.Time
    
    Game LotteryGame `gorm:"foreignKey:GameID"`
}
```

### 6. WinningRecord (backend/model/winning_record.go)
```go
type WinningRecord struct {
    ID           uint64 `gorm:"primaryKey"`
    UserID       uint64 `gorm:"not null;index"`
    UserNumberID uint64 `gorm:"not null;index"`
    DrawResultID uint64 `gorm:"not null;index"`
    GameID       uint64 `gorm:"not null;index"`
    Period       string `gorm:"size:32;not null;index"`
    PrizeLevel   int    `gorm:"not null"`
    RedMatches   int    `gorm:"not null;default:0"`
    BlueMatches  int    `gorm:"not null;default:0"`
    PrizeAmount  int64  `gorm:"not null;default:0"`
    IsVerified   bool   `gorm:"default:false"`
    CreatedAt    time.Time
    UpdatedAt    time.Time
    
    User       User           `gorm:"foreignKey:UserID"`
    UserNumber UserNumber     `gorm:"foreignKey:UserNumberID"`
    DrawResult DrawResult     `gorm:"foreignKey:DrawResultID"`
    Game       LotteryGame    `gorm:"foreignKey:GameID"`
}
```

## 自定义类型

### NumberArray
```go
type NumberArray []int

func (na *NumberArray) Scan(value interface{}) error {
    // 实现数据库到Go类型的转换
}

func (na NumberArray) Value() (driver.Value, error) {
    // 实现Go类型到数据库的转换
}
```

## 数据库初始化

1. **使用SQL文件初始化**：
   ```bash
   mysql -u root -p < backend/sql/init.sql
   ```

2. **使用GORM自动迁移**：
   ```go
   import "your-project/backend/migration"
   
   err := migration.AutoMigrate(db)
   if err != nil {
       log.Fatal("数据库迁移失败:", err)
   }
   ```

## 中奖规则

### 双色球中奖规则
- 一等奖: 6红+1蓝
- 二等奖: 6红+0蓝
- 三等奖: 5红+1蓝
- 四等奖: 5红+0蓝 或 4红+1蓝
- 五等奖: 4红+0蓝 或 3红+1蓝
- 六等奖: 2红+1蓝 或 1红+1蓝 或 0红+1蓝

### 大乐透中奖规则
- 一等奖: 5红+2蓝
- 二等奖: 5红+1蓝
- 三等奖: 5红+0蓝 或 4红+2蓝
- 四等奖: 4红+1蓝 或 3红+2蓝
- 五等奖: 4红+0蓝 或 3红+1蓝 或 2红+2蓝
- 六等奖: 3红+0蓝 或 2红+1蓝 或 1红+2蓝 或 0红+2蓝

## 性能优化

1. **索引策略**：
   - 所有外键字段都有索引
   - 经常查询的字段组合建立复合索引
   - 时间字段建立索引用于排序

2. **查询优化**：
   - 用户号码查询使用复合索引 (user_id, game_id, is_active)
   - 开奖结果查询使用复合索引 (game_id, period DESC, draw_date DESC)
   - 中奖记录查询使用复合索引 (user_id, game_id, period DESC)

3. **数据类型优化**：
   - 金额字段使用 bigint 存储分单位，避免浮点数精度问题
   - JSON字段存储数组类型数据，减少关联表
   - 使用合适的字符串长度限制

## 数据迁移注意事项

1. **版本控制**：每次表结构变更都应该有对应的迁移文件
2. **数据备份**：生产环境迁移前必须备份数据
3. **兼容性**：新字段应该有默认值，避免影响现有数据
4. **索引管理**：大表添加索引应该在低峰期进行

## 安全考虑

1. **数据权限**：用户只能访问自己的号码和中奖记录
2. **SQL注入防护**：使用GORM的参数化查询
3. **数据加密**：敏感信息（如有）应该加密存储
4. **访问控制**：数据库账号使用最小权限原则 