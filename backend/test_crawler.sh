#!/bin/bash

echo "=== 彩票开奖数据抓取测试 ==="

BASE_URL="http://localhost:8080"

echo "1. 测试生成模拟数据..."
curl -X POST "$BASE_URL/api/crawler/mock/ssq?period=2025099" | jq .

echo -e "\n2. 测试抓取功能（仅抓取不保存）..."
curl -X GET "$BASE_URL/api/crawler/test/ssq" | jq .

echo -e "\n3. 测试抓取并保存..."
curl -X POST "$BASE_URL/api/crawler/crawl/ssq" | jq .

echo -e "\n4. 查看开奖结果..."
curl -X GET "$BASE_URL/api/results/ssq" | jq .

echo -e "\n=== 测试完成 ===" 