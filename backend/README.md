# 彩票号码生成器后端

## 项目概述

这是一个基于Go语言开发的彩票号码生成器后端服务，支持双色球、大乐透等主流彩票游戏的号码生成、收藏管理等功能。

## 功能特性

### 🎯 核心功能
- **多彩票支持**: 支持双色球、大乐透等主流彩票游戏
- **智能号码生成**: 提供真随机算法生成彩票号码
- **号码收藏管理**: 用户可以保存、管理自己的号码
- **开奖结果查询**: 提供历史开奖数据查询

### 🔧 技术特性
- **RESTful API**: 标准的REST API设计
- **数据库设计**: 基于MySQL的优化数据库结构
- **JSON支持**: 使用JSON存储号码数据，便于扩展
- **模块化架构**: 清晰的MVC架构设计

## 技术栈

- **编程语言**: Go 1.21+
- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: MySQL 8.0+
- **配置管理**: INI配置文件

## 项目结构

```
backend/
├── main.go                 # 主程序入口
├── go.mod                  # Go模块文件
├── go.sum                  # Go依赖锁定文件
├── API_DOCS.md            # API文档
├── README.md              # 项目说明
├── test_api.sh            # API测试脚本
├── api/                   # API控制器层
│   ├── routes.go          # 路由注册
│   ├── user.go            # 用户相关API
│   ├── game.go            # 游戏相关API
│   ├── number.go          # 号码相关API
│   └── result.go          # 开奖结果API
├── service/               # 业务逻辑层
│   ├── user_service.go    # 用户服务
│   ├── game_service.go    # 游戏服务
│   ├── number_service.go  # 号码服务
│   ├── result_service.go  # 开奖结果服务
│   └── init_service.go    # 初始化服务
├── model/                 # 数据模型层
│   ├── user.go            # 用户模型
│   ├── lottery_game.go    # 彩票游戏模型
│   ├── user_number.go     # 用户号码模型
│   └── draw_result.go     # 开奖结果模型
├── migration/             # 数据库迁移
│   └── migrate.go         # 迁移脚本
├── common/                # 公共组件
│   ├── config/            # 配置管理
│   ├── mysql/             # MySQL连接
│   ├── redis/             # Redis连接
│   └── util/              # 工具函数
└── sql/                   # SQL脚本
    └── init.sql           # 数据库初始化脚本
```

## 安装与运行

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Git

### 1. 克隆项目
```bash
git clone <repository-url>
cd lucky/backend
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置数据库
创建MySQL数据库并配置`common/config/config.ini`:
```ini
[mysql]
user = root
password = your_password
host = localhost
port = 3306
db = lottery_db
```

### 4. 启动服务
```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

### 5. 测试API
```bash
# 使用测试脚本
chmod +x test_api.sh
./test_api.sh

# 或手动测试
curl http://localhost:8080/ping
```

## API接口

详细的API文档请参考 [API_DOCS.md](./API_DOCS.md)

### 主要接口概览

| 接口 | 方法 | 描述 |
|------|------|------|
| `/ping` | GET | 健康检查 |
| `/api/user/login` | POST | 用户登录 |
| `/api/games` | GET | 获取游戏列表 |
| `/api/numbers/random` | POST | 生成随机号码 |
| `/api/numbers/save` | POST | 保存号码 |
| `/api/numbers/my` | GET | 获取我的号码 |
| `/api/results/:gameCode` | GET | 获取开奖结果 |

## 数据库设计

### 核心表结构

1. **users**: 用户表
   - 存储用户基本信息
   - 支持微信小程序登录

2. **lottery_games**: 彩票游戏表
   - 游戏配置信息
   - 红球/蓝球数量配置

3. **user_numbers**: 用户号码表
   - 用户收藏的号码
   - JSON存储具体号码

4. **draw_results**: 开奖结果表
   - 历史开奖数据
   - 奖项信息

## 开发指南

### 添加新彩票游戏

1. 在 `service/init_service.go` 中添加游戏配置
2. 根据需要调整号码验证逻辑
3. 更新API文档

### 自定义号码算法

在 `service/number_service.go` 中的 `generateRandomBalls` 函数可以自定义随机算法。

### 扩展API功能

1. 在对应的 `api/*.go` 文件中添加新的handler
2. 在 `api/routes.go` 中注册路由
3. 在 `service/*.go` 中实现业务逻辑

## 部署

### Docker部署（推荐）
```bash
# 构建镜像
docker build -t lottery-backend .

# 运行容器
docker run -d -p 8080:8080 lottery-backend
```

### 传统部署
```bash
# 编译
go build -o lottery-backend main.go

# 运行
./lottery-backend
```

## 性能优化

- 数据库索引优化
- 连接池配置
- Redis缓存（可选）
- 分页查询优化

## 安全考虑

- 输入验证
- SQL注入防护
- 数据权限控制
- 敏感信息加密

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交更改
4. 发起Pull Request

## 许可证

MIT License

## 联系方式

如有问题，请创建Issue或联系开发者。