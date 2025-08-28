# 彩票号码生成器 API 文档

## 概述

彩票号码生成器后端API提供了完整的彩票号码管理功能，包括用户管理、号码生成、收藏管理和开奖结果查询。

## 基础信息

- 基础URL: `http://localhost:8080`
- 返回格式: JSON
- 编码: UTF-8

## 通用响应格式

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

## 认证

目前使用简单的Header认证方式（后续可升级为JWT）：
- Header: `X-User-ID: {userID}`

## API 接口

### 1. 系统测试

#### GET /ping
测试服务是否正常运行

**响应示例：**
```json
{
  "message": "pong"
}
```

### 2. 用户管理

#### POST /api/user/login
用户登录/注册

**请求参数：**
```json
{
  "openId": "wx123456789",
  "nickname": "用户昵称",
  "avatarUrl": "https://avatar.url"
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "登录成功",
  "data": {
    "userId": 1,
    "openId": "wx123456789",
    "nickname": "用户昵称",
    "avatarUrl": "https://avatar.url"
  }
}
```

#### GET /api/user/info
获取用户信息

**请求头：**
- `X-User-ID: 1`

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "userId": 1,
    "openId": "wx123456789",
    "nickname": "用户昵称",
    "avatarUrl": "https://avatar.url"
  }
}
```

### 3. 彩票游戏

#### GET /api/games
获取游戏列表

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "id": 1,
      "gameCode": "ssq",
      "gameName": "双色球",
      "redBallCount": 33,
      "blueBallCount": 16,
      "redSelectCount": 6,
      "blueSelectCount": 1,
      "isActive": true
    },
    {
      "id": 2,
      "gameCode": "dlt", 
      "gameName": "大乐透",
      "redBallCount": 35,
      "blueBallCount": 12,
      "redSelectCount": 5,
      "blueSelectCount": 2,
      "isActive": true
    }
  ]
}
```

#### GET /api/games/:gameCode
获取游戏详情

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "gameCode": "ssq",
    "gameName": "双色球",
    "redBallCount": 33,
    "blueBallCount": 16,
    "redSelectCount": 6,
    "blueSelectCount": 1,
    "isActive": true
  }
}
```

### 4. 号码管理

#### POST /api/numbers/random
生成随机号码

**请求头：**
- `X-User-ID: 1`

**请求参数：**
```json
{
  "gameCode": "ssq",
  "count": 1
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": [
    {
      "redBalls": [1, 5, 12, 18, 25, 33],
      "blueBalls": [8]
    }
  ]
}
```

#### POST /api/numbers/save
保存用户号码

**请求头：**
- `X-User-ID: 1`

**请求参数：**
```json
{
  "gameCode": "ssq",
  "redBalls": [1, 5, 12, 18, 25, 33],
  "blueBalls": [8],
  "nickname": "我的幸运号码",
  "source": "manual"
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "保存成功",
  "data": {
    "id": 1,
    "userId": 1,
    "gameId": 1,
    "redBalls": [1, 5, 12, 18, 25, 33],
    "blueBalls": [8],
    "nickname": "我的幸运号码",
    "source": "manual",
    "isActive": true,
    "createdAt": "2023-12-01T10:00:00Z"
  }
}
```

#### GET /api/numbers/my
获取我的号码

**请求头：**
- `X-User-ID: 1`

**查询参数：**
- `gameCode`: 游戏代码（可选）
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认20）

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "userId": 1,
        "gameId": 1,
        "redBalls": [1, 5, 12, 18, 25, 33],
        "blueBalls": [8],
        "nickname": "我的幸运号码",
        "source": "manual",
        "isActive": true,
        "createdAt": "2023-12-01T10:00:00Z",
        "game": {
          "gameCode": "ssq",
          "gameName": "双色球"
        }
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  }
}
```

#### PUT /api/numbers/:id
更新号码信息

**请求头：**
- `X-User-ID: 1`

**请求参数：**
```json
{
  "nickname": "新的昵称",
  "isActive": false
}
```

**响应示例：**
```json
{
  "code": 200,
  "message": "更新成功"
}
```

#### DELETE /api/numbers/:id
删除号码

**请求头：**
- `X-User-ID: 1`

**响应示例：**
```json
{
  "code": 200,
  "message": "删除成功"
}
```

### 5. 开奖结果

#### GET /api/results/:gameCode
获取开奖结果列表

**查询参数：**
- `page`: 页码（默认1）
- `pageSize`: 每页数量（默认20）

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "gameId": 1,
        "period": "2023130",
        "drawDate": "2023-12-01T20:00:00Z",
        "redBalls": [1, 5, 12, 18, 25, 33],
        "blueBalls": [8],
        "salesAmount": 500000000,
        "prizePool": 100000000,
        "firstPrize": 10,
        "firstAmount": 5000000,
        "game": {
          "gameCode": "ssq",
          "gameName": "双色球"
        }
      }
    ],
    "total": 1,
    "page": 1,
    "pageSize": 20
  }
}
```

#### GET /api/results/:gameCode/:period
获取指定期号开奖结果

**响应示例：**
```json
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "gameId": 1,
    "period": "2023130",
    "drawDate": "2023-12-01T20:00:00Z",
    "redBalls": [1, 5, 12, 18, 25, 33],
    "blueBalls": [8],
    "salesAmount": 500000000,
    "prizePool": 100000000,
    "firstPrize": 10,
    "firstAmount": 5000000,
    "game": {
      "gameCode": "ssq",
      "gameName": "双色球"
    }
  }
}
```

## 错误码说明

- `200`: 成功
- `400`: 请求参数错误
- `401`: 未授权
- `404`: 资源不存在
- `500`: 服务器内部错误

## 开发测试

使用以下命令启动服务：
```bash
cd backend
go run main.go
```

可以使用curl或Postman等工具测试API接口。

## 6. 数据抓取接口

### 6.1 生成模拟开奖数据

**接口**: `POST /api/crawler/mock/{gameCode}`

**描述**: 生成模拟开奖数据并保存到数据库

**请求参数**:
- `gameCode`: 游戏代码 (路径参数)
- `period`: 期号 (查询参数，可选)

**请求示例**:
```bash
curl -X POST "http://localhost:8080/api/crawler/mock/ssq?period=2025099"
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "模拟数据生成成功",
  "data": {
    "period": "2025099",
    "draw_date": "2025-08-28",
    "red_balls": [5, 8, 13, 17, 18, 29],
    "blue_balls": [2],
    "sales": 350000000,
    "pool_amount": 2500000000,
    "game_code": "ssq"
  }
}
```

### 6.2 测试抓取功能

**接口**: `GET /api/crawler/test/{gameCode}`

**描述**: 测试从外部网站抓取开奖数据（不保存到数据库）

**请求参数**:
- `gameCode`: 游戏代码 (路径参数)

**请求示例**:
```bash
curl "http://localhost:8080/api/crawler/test/ssq"
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "抓取成功",
  "data": {
    "period": "2025098",
    "draw_date": "2025-08-26",
    "red_balls": [5, 8, 13, 17, 18, 29],
    "blue_balls": [2],
    "sales": 348777356,
    "pool_amount": 2518876889,
    "game_code": "ssq"
  }
}
```

### 6.3 抓取并保存开奖数据

**接口**: `POST /api/crawler/crawl/{gameCode}`

**描述**: 从外部网站抓取最新开奖数据并保存到数据库

**请求参数**:
- `gameCode`: 游戏代码 (路径参数)

**请求示例**:
```bash
curl -X POST "http://localhost:8080/api/crawler/crawl/ssq"
```

**响应示例**:
```json
{
  "code": 0,
  "msg": "抓取成功"
}
```

### 6.4 数据源说明

系统支持多个数据源，按优先级自动切换：

1. **500彩票网** (优先级: 1)
   - 支持双色球和大乐透
   - 数据更新及时，格式相对稳定

2. **乐彩网** (优先级: 2)
   - 主要支持双色球
   - 历史数据完整

3. **新浪彩票** (优先级: 3)
   - 备用数据源
   - 待完善实现

### 6.5 命令行工具

项目提供了命令行工具用于数据抓取管理：

```bash
# 编译命令行工具
cd backend
go build -o crawler cmd/crawler.go

# 测试抓取
./crawler -action=test -game=ssq

# 抓取并保存
./crawler -action=crawl -game=ssq

# 生成模拟数据
./crawler -action=mock -game=ssq -period=2025099

# 启动定时抓取任务
./crawler -action=schedule
```

### 6.6 定时任务

系统支持定时抓取功能：
- 每30分钟检查一次新的开奖数据
- 自动抓取双色球和大乐透数据
- 支持多数据源容错机制

### 6.7 注意事项

1. **反爬措施**: 
   - 控制抓取频率，避免被目标网站封禁
   - 支持多数据源切换

2. **数据去重**: 
   - 系统会自动检查期号是否已存在
   - 避免重复保存相同期号的数据

3. **错误处理**: 
   - 网络异常时自动重试其他数据源
   - 详细的错误日志记录

4. **权限控制**: 
   - 抓取接口建议仅对管理员开放
   - 可通过中间件添加权限验证

## 总结

本API文档涵盖了彩票号码生成器的所有核心功能，包括用户管理、游戏管理、号码生成与收藏、开奖结果查询、数据抓取等。系统提供了完整的开奖数据获取方案，支持多数据源抓取、定时任务、命令行工具等功能，确保数据的及时性和可靠性。所有接口都提供了详细的请求示例和响应格式，便于前端开发和第三方集成。 