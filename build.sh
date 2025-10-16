#!/bin/bash

# LNMP运维面板构建脚本 - ARMv7l架构优化

echo "开始构建LNMP运维面板..."

# 设置环境变量
export GOARCH=arm
export GOOS=linux
export GOARM=7

# 清理之前的构建
echo "清理构建缓存..."
rm -rf bin/
mkdir -p bin/

# 下载依赖
echo "下载依赖..."
go mod tidy

# 构建ARMv7l版本
echo "构建ARMv7l版本..."
go build -o bin/lnmp-panel-armv7l -ldflags="-s -w" main.go

# 构建本地测试版本（可选）
echo "构建本地测试版本..."
go build -o bin/lnmp-panel-local main.go

# 检查文件大小
echo "构建完成，文件信息："
ls -lh bin/

echo "构建成功！"
echo "ARMv7l版本: bin/lnmp-panel-armv7l"
echo "本地版本: bin/lnmp-panel-local"