#!/bin/bash

# 彩票号码生成器 API 测试脚本
BASE_URL="http://localhost:8080"

echo "=== 彩票号码生成器 API 测试 ==="
echo ""

# 1. 测试ping接口
echo "1. 测试服务是否运行..."
curl -s "$BASE_URL/ping" | jq '.'
echo ""

# 2. 测试获取游戏列表
echo "2. 获取游戏列表..."
curl -s "$BASE_URL/api/games" | jq '.'
echo ""

# 3. 测试获取双色球游戏详情
echo "3. 获取双色球游戏详情..."
curl -s "$BASE_URL/api/games/ssq" | jq '.'
echo ""

# 4. 测试用户登录
echo "4. 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/user/login" \
  -H "Content-Type: application/json" \
  -d '{
    "openId": "test_user_001",
    "nickname": "测试用户",
    "avatarUrl": "https://avatar.example.com/test.jpg"
  }')
echo "$LOGIN_RESPONSE" | jq '.'

# 提取用户ID
USER_ID=$(echo "$LOGIN_RESPONSE" | jq -r '.data.userId')
echo "用户ID: $USER_ID"
echo ""

# 5. 测试生成随机号码
echo "5. 生成双色球随机号码..."
curl -s -X POST "$BASE_URL/api/numbers/random" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "gameCode": "ssq",
    "count": 3
  }' | jq '.'
echo ""

# 6. 测试保存号码
echo "6. 保存一注双色球号码..."
SAVE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/numbers/save" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "gameCode": "ssq",
    "redBalls": [1, 5, 12, 18, 25, 33],
    "blueBalls": [8],
    "nickname": "我的幸运号码",
    "source": "manual"
  }')
echo "$SAVE_RESPONSE" | jq '.'

# 提取号码ID
NUMBER_ID=$(echo "$SAVE_RESPONSE" | jq -r '.data.id')
echo "号码ID: $NUMBER_ID"
echo ""

# 7. 测试获取我的号码
echo "7. 获取我的号码列表..."
curl -s "$BASE_URL/api/numbers/my?gameCode=ssq" \
  -H "X-User-ID: $USER_ID" | jq '.'
echo ""

# 8. 测试更新号码昵称
echo "8. 更新号码昵称..."
curl -s -X PUT "$BASE_URL/api/numbers/$NUMBER_ID" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "nickname": "更新后的昵称"
  }' | jq '.'
echo ""

# 9. 测试获取用户信息
echo "9. 获取用户信息..."
curl -s "$BASE_URL/api/user/info" \
  -H "X-User-ID: $USER_ID" | jq '.'
echo ""

# 10. 测试大乐透随机号码生成
echo "10. 生成大乐透随机号码..."
curl -s -X POST "$BASE_URL/api/numbers/random" \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -d '{
    "gameCode": "dlt",
    "count": 2
  }' | jq '.'
echo ""

echo "=== 测试完成 ===" 