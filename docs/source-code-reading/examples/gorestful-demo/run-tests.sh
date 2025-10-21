#!/bin/bash

# Go-Restful Hello World 示例测试脚本

echo "🚀 Go-Restful Hello World 示例测试"
echo "=================================="
echo ""

# 检查服务器是否运行
check_server() {
    if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
        echo "✅ 服务器正在运行"
        return 0
    else
        echo "❌ 服务器未运行，请先启动服务器"
        echo "   运行: go run main.go 或 go run advanced-main.go"
        return 1
    fi
}

# 测试简单示例
test_simple() {
    echo "📋 测试简单示例接口"
    echo "-------------------"
    
    echo "1. 测试基本 Hello 接口"
    curl -s http://localhost:8080/api/v1/hello | jq .
    echo ""
    
    echo "2. 测试带参数的 Hello 接口"
    curl -s "http://localhost:8080/api/v1/hello?name=Kubernetes" | jq .
    echo ""
}

# 测试高级示例
test_advanced() {
    echo "📋 测试高级示例接口"
    echo "-------------------"
    
    echo "1. 测试 GET Hello 接口"
    curl -s "http://localhost:8080/api/v1/hello?name=Kubernetes" | jq .
    echo ""
    
    echo "2. 测试 POST Hello 接口"
    curl -s -X POST http://localhost:8080/api/v1/hello \
        -H 'Content-Type: application/json' \
        -d '{"name":"Kubernetes","age":25}' | jq .
    echo ""
    
    echo "3. 测试获取用户列表"
    curl -s http://localhost:8080/api/v1/users | jq .
    echo ""
    
    echo "4. 测试根据 ID 获取用户"
    curl -s http://localhost:8080/api/v1/users/1 | jq .
    echo ""
    
    echo "5. 测试健康检查"
    curl -s http://localhost:8080/api/v1/health | jq .
    echo ""
    
    echo "6. 测试欢迎页面"
    curl -s http://localhost:8080/ | jq .
    echo ""
}

# 主函数
main() {
    if ! check_server; then
        exit 1
    fi
    
    echo "请选择测试类型:"
    echo "1) 简单示例测试"
    echo "2) 高级示例测试"
    echo "3) 全部测试"
    echo ""
    read -p "请输入选择 (1-3): " choice
    
    case $choice in
        1)
            test_simple
            ;;
        2)
            test_advanced
            ;;
        3)
            test_simple
            echo ""
            test_advanced
            ;;
        *)
            echo "❌ 无效选择"
            exit 1
            ;;
    esac
    
    echo "✅ 测试完成！"
}

# 检查 jq 是否安装
if ! command -v jq &> /dev/null; then
    echo "⚠️  警告: jq 未安装，JSON 输出将不会格式化"
    echo "   安装 jq: brew install jq (macOS) 或 apt-get install jq (Ubuntu)"
    echo ""
fi

# 运行主函数
main
