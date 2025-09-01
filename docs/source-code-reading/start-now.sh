#!/bin/bash
# K8s专家速成 - 立即开始脚本
# 适用于10年Java开发老兵

set -e

echo "🚀 欢迎Java老兵踏上K8s专家之路!"
echo "这个脚本将帮你完成第一天的环境搭建和学习"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 创建学习目录
LEARNING_DIR="$HOME/k8s-expert-journey"
echo -e "${BLUE}📁 创建学习目录: $LEARNING_DIR${NC}"
mkdir -p "$LEARNING_DIR"
cd "$LEARNING_DIR"

# 初始化学习日志
echo "# K8s专家成长日记" > README.md
echo "开始日期: $(date)" >> README.md
echo "目标: 30天成为K8s源码专家" >> README.md
echo "背景: 10年Java开发老兵" >> README.md
echo "" >> README.md
echo "## 学习进度" >> README.md

echo -e "${GREEN}✅ 学习目录创建完成${NC}"

# 检查Go环境
echo -e "${BLUE}🔍 检查Go环境...${NC}"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo -e "${GREEN}✅ Go已安装: $GO_VERSION${NC}"
else
    echo -e "${YELLOW}⚠️ Go未安装，正在安装...${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            brew install go
        else
            echo -e "${RED}❌ 请先安装Homebrew: https://brew.sh${NC}"
            exit 1
        fi
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        # Linux
        wget -q https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
        sudo rm -rf /usr/local/go
        sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
        export PATH=$PATH:/usr/local/go/bin
        echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
        rm go1.21.5.linux-amd64.tar.gz
    else
        echo -e "${RED}❌ 不支持的操作系统，请手动安装Go: https://golang.org/dl/${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Go安装完成${NC}"
fi

# 检查kubectl
echo -e "${BLUE}🔍 检查kubectl...${NC}"
if command -v kubectl &> /dev/null; then
    echo -e "${GREEN}✅ kubectl已安装${NC}"
    kubectl cluster-info --request-timeout=5s >/dev/null 2>&1 && echo -e "${GREEN}✅ K8s集群连接正常${NC}" || echo -e "${YELLOW}⚠️ 无法连接K8s集群，请检查配置${NC}"
else
    echo -e "${YELLOW}⚠️ kubectl未安装，请安装kubectl${NC}"
    echo "安装指南: https://kubernetes.io/docs/tasks/tools/"
fi

# 创建第一个Go项目
echo -e "${BLUE}🛠️ 创建第一个Go项目...${NC}"
mkdir -p day1-go-basics
cd day1-go-basics

# 初始化Go模块
go mod init k8s-learning-day1

# 创建示例文件
cat > main.go << 'EOF'
// Day 1 - 第一个Go程序 (Java开发者版)
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("🎉 恭喜！你的第一个Go程序运行成功了！")
    fmt.Printf("当前时间: %s\n", time.Now().Format("2006-01-02 15:04:05"))
    
    // Java vs Go 语法对比
    fmt.Println("\n=== Java vs Go 快速对比 ===")
    
    // 1. 变量声明
    var javaStyle string = "Java风格声明"
    goStyle := "Go风格声明"  // 类型推断，类似Java 10的var
    fmt.Printf("Java风格: %s\n", javaStyle)
    fmt.Printf("Go风格: %s\n", goStyle)
    
    // 2. 函数调用
    result := addNumbers(10, 20)
    fmt.Printf("函数调用结果: %d\n", result)
    
    // 3. 错误处理 (Go特色)
    value, err := divide(10, 2)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
    } else {
        fmt.Printf("除法结果: %.2f\n", value)
    }
    
    fmt.Println("\n🎯 完成第一步！接下来学习client-go")
}

// 函数定义 - 类似Java方法
func addNumbers(a, b int) int {
    return a + b
}

// 多返回值 - Go特色，Java需要用对象包装
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("不能除以零")
    }
    return a / b, nil
}
EOF

# 运行第一个程序
echo -e "${BLUE}🏃 运行第一个Go程序...${NC}"
go run main.go

echo -e "${GREEN}✅ 第一个Go程序运行成功！${NC}"

# 返回主目录
cd "$LEARNING_DIR"

# 记录今天的成就
echo "" >> README.md
echo "### Day 1 ($(date +%Y-%m-%d))" >> README.md
echo "- [x] 环境搭建完成" >> README.md
echo "- [x] 第一个Go程序运行成功" >> README.md
echo "- [x] 理解Go基本语法" >> README.md
echo "- [ ] client-go入门 (明天完成)" >> README.md

# 准备client-go学习环境
echo -e "${BLUE}🛠️ 准备client-go学习环境...${NC}"
mkdir -p day2-client-go
cd day2-client-go

go mod init k8s-learning-day2
go get k8s.io/client-go@latest
go get k8s.io/api@latest
go get k8s.io/apimachinery@latest

echo -e "${GREEN}✅ client-go依赖安装完成${NC}"

# 创建明天的学习计划
cat > tomorrow-plan.md << 'EOF'
# Day 2 学习计划

## 目标
- 掌握client-go基本使用
- 理解K8s资源操作
- 学会Watch机制

## 任务清单
- [ ] 运行client-go示例程序
- [ ] 创建一个Java应用的Deployment
- [ ] 学习Watch机制监听资源变化
- [ ] 对比client-go与Java数据库操作的相似性

## 准备工作
- [x] 安装client-go依赖
- [x] 确保kubectl可以访问集群
- [ ] 下载Day2示例代码

## 预计用时
2-3小时
EOF

echo -e "${GREEN}✅ 明天的学习计划已准备好${NC}"

cd "$LEARNING_DIR"

# 初始化git仓库
echo -e "${BLUE}📝 初始化Git仓库...${NC}"
git init >/dev/null 2>&1 || true
git add . >/dev/null 2>&1 || true
git commit -m "开始K8s专家成长之旅 - Day 1" >/dev/null 2>&1 || true

# 总结今天的成就
echo ""
echo -e "${GREEN}🎉 恭喜完成Day 1的学习！${NC}"
echo ""
echo "📋 今天完成的任务:"
echo -e "  ${GREEN}✅${NC} 创建学习环境"
echo -e "  ${GREEN}✅${NC} 安装Go开发环境"
echo -e "  ${GREEN}✅${NC} 运行第一个Go程序"
echo -e "  ${GREEN}✅${NC} 理解Go基本语法"
echo -e "  ${GREEN}✅${NC} 准备client-go学习环境"
echo ""
echo "📚 学习目录结构:"
echo "  $LEARNING_DIR/"
echo "  ├── README.md (学习日记)"
echo "  ├── day1-go-basics/ (今天的代码)"
echo "  └── day2-client-go/ (明天的环境)"
echo ""
echo -e "${BLUE}🎯 明天的任务:${NC}"
echo "1. 进入目录: cd $LEARNING_DIR/day2-client-go"
echo "2. 下载示例代码到该目录"
echo "3. 运行client-go示例程序"
echo "4. 完成Day 2的学习目标"
echo ""
echo -e "${YELLOW}💡 重要提示:${NC}"
echo "- 每天只需要投入2-3小时"
echo "- 遇到问题时先查看错误信息"
echo "- 可以随时回来查看学习日记"
echo "- 不要担心进度，坚持最重要"
echo ""
echo -e "${GREEN}🚀 你已经开始了专家之路的第一步！${NC}"
echo -e "${GREEN}明天见，继续加油！💪${NC}"
