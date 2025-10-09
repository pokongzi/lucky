#!/bin/bash

# 构建Linux版本的程序
echo "开始构建Linux版本..."

# 设置环境变量进行交叉编译
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

# 构建后端程序
echo "构建后端程序..."
cd backend
go build -o ../lucky-linux-amd64 .

# 返回根目录
cd ..

# 检查构建是否成功
if [ -f "lucky-linux-amd64" ]; then
    echo "Linux版本构建成功: lucky-linux-amd64"
    
    # 显示文件信息
    ls -lh lucky-linux-amd64
    
    # 创建部署包
    echo "创建部署包..."
    tar -czf lucky-linux-amd64.tar.gz lucky-linux-amd64
    
    echo "部署包创建成功: lucky-linux-amd64.tar.gz"
else
    echo "构建失败!"
    exit 1
fi

echo "构建完成!"