#!/bin/bash

# 猫书项目整合启动脚本
# 同时启动catcal主业务服务和catchat-main聊天服务

echo "🐱 猫书项目整合启动脚本"
echo "================================"

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "❌ Go环境未安装，请先安装Go"
    exit 1
fi

# 检查PostgreSQL连接
echo "📊 检查PostgreSQL数据库连接..."
pg_isready -h localhost -p 5432 -U catcal
if [ $? -ne 0 ]; then
    echo "❌ PostgreSQL连接失败，请确保数据库服务正在运行"
    echo "💡 启动命令: brew services start postgresql"
    exit 1
fi

# 创建必要的目录
mkdir -p logs
mkdir -p tmp

# 启动catcal主业务服务
echo "🚀 启动catcal主业务服务 (端口: 8082)..."
cd ../catcal
if [ ! -f "catcal" ]; then
    echo "🔨 编译catcal项目..."
    go build -o catcal main.go
fi

# 后台启动catcal
./catcal > ../catchat-main/logs/catcal.log 2>&1 &
CATCAL_PID=$!
echo "✅ catcal服务已启动 (PID: $CATCAL_PID)"

# 等待catcal启动
echo "⏳ 等待catcal服务启动..."
sleep 5

# 检查catcal是否启动成功
if curl -s http://localhost:8082/api/v1/health > /dev/null 2>&1; then
    echo "✅ catcal服务启动成功"
else
    echo "⚠️  catcal服务可能未完全启动，继续启动聊天服务..."
fi

# 启动catchat-main聊天服务
echo "💬 启动catchat-main聊天服务 (端口: 8888)..."
cd ../catchat-main

# 检查依赖
echo "📦 检查Go依赖..."
go mod tidy

# 检查是否安装了air
if ! command -v air &> /dev/null; then
    echo "❌ air未安装，正在安装..."
    curl -sSfL https://raw.githubusercontent.com/air-verse/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
    export PATH=$HOME/go/bin:$PATH
fi

echo "🎯 启动聊天服务器 (热更新模式)..."
echo "📝 修改代码后会自动重启服务"
echo "🛑 按 Ctrl+C 停止所有服务"
echo ""

# 保存PID到文件
echo $CATCAL_PID > tmp/catcal.pid

# 使用air启动聊天服务
air

# 清理函数
cleanup() {
    echo ""
    echo "🛑 正在停止所有服务..."
    
    # 停止catcal
    if [ -f "tmp/catcal.pid" ]; then
        CATCAL_PID=$(cat tmp/catcal.pid)
        kill $CATCAL_PID 2>/dev/null
        rm -f tmp/catcal.pid
        echo "✅ catcal服务已停止"
    fi
    
    # 停止air进程
    pkill -f "air" 2>/dev/null
    echo "✅ 聊天服务已停止"
    
    echo "🎉 所有服务已停止"
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 等待用户中断
wait 