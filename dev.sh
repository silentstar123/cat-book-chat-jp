#!/bin/bash

# 设置环境变量，跳过MySQL连接
export SKIP_MYSQL=true
export PATH=$HOME/go/bin:$PATH

# 检查依赖
echo "检查Go依赖..."
go mod tidy

# 检查是否安装了air
if ! command -v air &> /dev/null; then
    echo "❌ air未安装，正在安装..."
    curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    export PATH=$HOME/go/bin:$PATH
fi

# 创建日志目录和tmp目录
mkdir -p logs
mkdir -p tmp

# 启动聊天服务器（热更新模式）
echo "💬 启动聊天服务器 (端口: 8888) - 热更新模式..."
echo "📝 修改代码后会自动重启服务"
echo "🛑 按 Ctrl+C 停止服务"
echo ""

# 使用air进行热更新
air 