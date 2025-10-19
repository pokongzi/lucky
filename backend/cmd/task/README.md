# 抓取任务服务 (Task Service)

内网混合服务，同时支持 HTTP 和 gRPC 接口。

## 服务说明

- **HTTP 端口**: 8081 - 用于直接 HTTP/1.1 调用
- **gRPC 端口**: 9091 - 用于 gocron 调度系统
- **用途**: 提供内网调用的抓取任务接口
- **功能**: 抓取并保存彩票开奖数据

## 启动服务

```bash
cd backend/cmd/task
go run crawler.go
```

或者编译后运行：

```bash
cd backend/cmd/task
go build -o task-crawler crawler.go
./task-crawler
```

## API 接口

### 1. 健康检查

**请求**:
```
GET /health
```

**响应**:
```json
{
  "code": 0,
  "message": "OK"
}
```

### 2. 抓取任务 (POST)

**请求**:
```
POST /task/crawl
Content-Type: application/json

{
  "game_code": "ssq"
}
```

**参数说明**:
- `game_code`: 游戏代码，支持 `ssq`(双色球) 或 `dlt`(大乐透)

**响应**:
```json
{
  "code": 0,
  "message": "抓取成功",
  "data": {
    "game_code": "ssq"
  }
}
```

### 3. 抓取任务 (GET)

**请求**:
```
GET /task/crawl/:gameCode
```

**示例**:
```
GET /task/crawl/ssq
GET /task/crawl/dlt
```

**响应**:
```json
{
  "code": 0,
  "message": "抓取成功",
  "data": {
    "game_code": "ssq"
  }
}
```

## 使用示例

### 使用 curl

```bash
# 健康检查
curl http://localhost:8081/health

# POST 方式抓取双色球数据
curl -X POST http://localhost:8081/task/crawl \
  -H "Content-Type: application/json" \
  -d '{"game_code":"ssq"}'

# GET 方式抓取大乐透数据
curl http://localhost:8081/task/crawl/dlt
```

### 使用 HTTP 客户端

```bash
# 抓取双色球
POST http://localhost:8081/task/crawl
{
  "game_code": "ssq"
}

# 抓取大乐透
GET http://localhost:8081/task/crawl/dlt
```

## 错误码说明

| Code | 说明 |
|------|------|
| 0    | 成功 |
| 400  | 参数错误 |
| 500  | 服务器内部错误 |

## gocron 配置说明

### 1. 添加任务节点

在 gocron 中添加任务节点：

```
节点名称: task-crawler
IP地址: 你的服务器IP
端口: 9091
备注: 彩票抓取任务节点
```

### 2. 创建任务

**双色球抓取任务：**
```
任务名称: 抓取双色球
节点: task-crawler
任务命令: ssq
定时规则: 每天 21:30 执行
```

**大乐透抓取任务：**
```
任务名称: 抓取大乐透  
节点: task-crawler
任务命令: dlt
定时规则: 每周一、三、六 21:30 执行
```

### 3. 验证连接

gocron 会自动调用 gRPC 的 Check 方法验证节点是否在线。

## HTTP 接口（保留原有功能）

HTTP 接口继续可用，适合直接调用或脚本调用：

```bash
# 使用 curl 测试
curl http://localhost:8081/task/crawl/ssq
curl http://localhost:8081/task/crawl/dlt
```

## 注意事项

1. 该服务仅供内网调用，请勿暴露到公网
2. 支持的游戏代码: `ssq` (双色球), `dlt` (大乐透)
3. 抓取任务是同步执行的，会等待抓取完成后返回结果
4. gRPC 端口 (9091) 用于 gocron，HTTP 端口 (8081) 用于直接调用
5. 两个端口都需要在防火墙中开放

## 端口说明

| 端口 | 协议 | 用途 |
|------|------|------|
| 8081 | HTTP/1.1 | 直接 HTTP 调用抓取接口 |
| 9091 | gRPC | gocron 任务调度 |

## 与主服务的区别

- **主服务 (backend/main.go)**: 端口 8080，对外提供完整的 API 服务
- **任务服务 (backend/cmd/task)**: 端口 8081(HTTP) + 9091(gRPC)，内网调用，专门用于执行数据抓取任务
- **命令行工具 (backend/cmd/command)**: 命令行界面，用于手动执行各种操作

