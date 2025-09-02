# Cursor 中的 Go 开发指南

## 概述

本文档介绍如何在 Cursor IDE 中高效地开发、编译和运行 Go 代码。

## 环境准备

### 1. 安装 Go
```bash
# macOS (使用 Homebrew)
brew install go

# 验证安装
go version
```

### 2. 安装必要的扩展

#### Go 扩展（必需）
```bash
cursor --install-extension golang.go
```

#### Code Runner 扩展（推荐）
```bash
cursor --install-extension formulahendry.code-runner
```

#### Delve 调试器（可选）
```bash
go install -v github.com/go-delve/delve/cmd/dlv@latest
```

## 项目设置

### 1. 创建 Go 模块
```bash
# 在项目目录中初始化模块
go mod init your-project-name

# 示例
go mod init hello
```

### 2. 基本项目结构
```
your-project/
├── go.mod
├── main.go
└── README.md
```

## 运行 Go 代码的方法

### 方法 1: 右键菜单运行（推荐）

1. **右键点击代码** → 选择 **"Run Code"**
2. 结果会在下方的 **OUTPUT** 面板显示
3. 支持快捷键：`Ctrl+Alt+N` (Windows/Linux) 或 `Cmd+Alt+N` (macOS)

### 方法 2: 使用运行按钮

1. 代码编辑器**右上角**的 ▶️ 按钮
2. 点击即可运行当前文件

### 方法 3: 快捷键运行

- **Ctrl+F5** (Windows/Linux) 或 **Cmd+F5** (macOS)：直接运行
- **F5**：调试运行

### 方法 4: CodeLens（如果可用）

在函数上方会显示 `run | debug` 链接，直接点击即可。

### 方法 5: 命令面板

1. **Ctrl+Shift+P** (Windows/Linux) 或 **Cmd+Shift+P** (macOS)
2. 输入 **"Go: Run Package"** 或 **"Go: Debug Package"**

### 方法 6: 集成终端

1. **Ctrl+`** (Windows/Linux) 或 **Cmd+`** (macOS) 打开终端
2. 运行命令：
   ```bash
   # 运行单个文件
   go run main.go
   
   # 运行整个包
   go run .
   
   # 编译后运行
   go build -o myapp && ./myapp
   ```

## 调试 Go 代码

### 1. 设置断点
- 在代码行号左侧**点击**设置断点
- 红点表示断点已设置

### 2. 开始调试
- **F5** 启动调试
- 或右键 → **"Debug Code"**

### 3. 调试控制
- **F10**：单步跳过
- **F11**：单步进入
- **Shift+F11**：单步跳出
- **F5**：继续执行

## 常见问题解决

### 问题 1: "main redeclared in this block"
**原因**：同一目录下有多个 `main` 函数

**解决方案**：
```bash
# 将文件移到单独目录
mkdir hello
mv main.go hello/
cd hello
go mod init hello
```

### 问题 2: "dlv command not available"
**解决方案**：
```bash
go install -v github.com/go-delve/delve/cmd/dlv@latest
```

### 问题 3: 右键没有 "Run Code" 选项
**解决方案**：
1. 安装 Code Runner 扩展
2. 重启 Cursor
3. 检查扩展是否启用

### 问题 4: CodeLens 不显示
**解决方案**：
1. 打开设置 (`Cmd+,`)
2. 搜索 "codelens"
3. 确保启用：
   - `go.enableCodeLens`
   - `editor.codeLens`

## 高级配置

### 1. 自定义运行配置

创建 `.vscode/launch.json`：
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}"
        }
    ]
}
```

### 2. 任务配置

创建 `.vscode/tasks.json`：
```json
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "go run",
            "type": "shell",
            "command": "go",
            "args": ["run", "main.go"],
            "group": {
                "kind": "build",
                "isDefault": true
            }
        }
    ]
}
```

## 最佳实践

1. **使用 go.mod**：始终在项目根目录创建 `go.mod` 文件
2. **代码格式化**：保存时自动格式化（`gofmt`）
3. **错误检查**：利用 Go 扩展的实时错误检查
4. **测试驱动**：编写测试文件 `*_test.go`
5. **版本控制**：使用 `.gitignore` 忽略构建产物

## 示例代码

### Hello World
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

### 带参数的程序
```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) > 1 {
        fmt.Printf("Hello, %s!\n", os.Args[1])
    } else {
        fmt.Println("Hello, World!")
    }
}
```

## 总结

Cursor 提供了多种运行 Go 代码的方式，从简单的右键运行到高级的调试功能。推荐使用 Code Runner 扩展来获得最佳的开发体验。

---
*更新时间：2024年9月*
