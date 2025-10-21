# Go-Restful HTTP 接口示例

这个目录包含了使用 Go-Restful 框架创建 HTTP 接口的示例代码。

## 📁 文件说明

### Go-Restful 示例项目
- `gorestful-demo/` - 完整的 Go-Restful Hello World 示例项目
  - `main.go` - 简单示例
  - `advanced-main.go` - 高级示例
  - `go.mod` - 独立模块管理
  - `README.md` - 详细说明文档
  - `run-tests.sh` - 自动测试脚本

### 其他示例
- `day1-go-quickstart.go` - Go 快速开始示例
- `day2-client-go-basic.go` - Kubernetes 客户端基础示例

## 🚀 快速开始

### 1. 进入 Go-Restful 示例项目

```bash
cd gorestful-demo
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 运行简单示例

```bash
go run main.go
```

### 4. 运行高级示例

```bash
go run advanced-main.go
```

### 5. 运行测试

```bash
./run-tests.sh
```

## 📋 API 接口

### 简单示例接口

| 方法 | 路径 | 参数 | 说明 |
|------|------|------|------|
| GET | `/api/v1/hello` | `name` (可选) | 返回 Hello World 消息 |

### 完整示例接口

| 方法 | 路径 | 参数 | 说明 |
|------|------|------|------|
| GET | `/api/v1/hello` | `name` (查询参数) | 返回 Hello 消息 |
| POST | `/api/v1/hello` | `name` (请求体) | 返回 Hello 消息 |
| GET | `/api/v1/health` | 无 | 健康检查 |
| GET | `/` | 无 | 欢迎页面 |

## 🔗 测试接口

### 使用 curl 测试

```bash
# 简单示例测试
curl http://localhost:8080/api/v1/hello
curl http://localhost:8080/api/v1/hello?name=Kubernetes

# 完整示例测试
curl http://localhost:8080/api/v1/hello
curl http://localhost:8080/api/v1/hello?name=Kubernetes
curl -X POST http://localhost:8080/api/v1/hello \
  -H 'Content-Type: application/json' \
  -d '{"name":"Kubernetes"}'
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/
```

### 预期响应

#### GET /api/v1/hello
```json
{
  "message": "Hello, World!",
  "status": "success"
}
```

#### GET /api/v1/hello?name=Kubernetes
```json
{
  "message": "Hello, Kubernetes!",
  "status": "success"
}
```

#### POST /api/v1/hello
```json
{
  "message": "Hello, Kubernetes! (POST)",
  "status": "success",
  "time": "2024-12-19"
}
```

## 🏗️ 代码结构说明

### 简单示例 (`simple-gorestful-hello.go`)

```go
// 1. 创建容器
container := restful.NewContainer()

// 2. 创建 WebService
ws := new(restful.WebService)
ws.Path("/api/v1")

// 3. 定义处理器函数
helloHandler := func(request *restful.Request, response *restful.Response) {
    // 处理逻辑
}

// 4. 注册路由
ws.Route(ws.GET("/hello").To(helloHandler))

// 5. 添加到容器并启动
container.Add(ws)
http.ListenAndServe(":8080", container)
```

### 完整示例 (`gorestful-hello-world.go`)

```go
// 1. 定义服务结构体
type HelloWorldService struct {}

// 2. 实现方法
func (h *HelloWorldService) GetHelloWorld(request *restful.Request, response *restful.Response) {
    // 处理逻辑
}

// 3. 注册路由方法
func (h *HelloWorldService) RegisterRoutes(container *restful.Container) {
    // 路由注册逻辑
}

// 4. 启动服务
helloService := NewHelloWorldService()
helloService.RegisterRoutes(container)
```

## 🔧 关键特性

### Go-Restful 框架特性
- ✅ **自动路由注册** - 基于路径和方法自动路由
- ✅ **JSON 绑定** - 自动序列化/反序列化 JSON
- ✅ **参数验证** - 查询参数和路径参数验证
- ✅ **Swagger 文档** - 自动生成 API 文档
- ✅ **中间件支持** - 过滤器链支持
- ✅ **CORS 支持** - 跨域资源共享

### 示例中使用的特性
- **路径参数**: `{name}` 格式的路径参数
- **查询参数**: `?name=value` 格式的查询参数
- **请求体绑定**: JSON 请求体自动绑定到结构体
- **响应序列化**: 结构体自动序列化为 JSON
- **错误处理**: 统一的错误响应格式
- **CORS 配置**: 跨域请求支持

## 📚 学习建议

### 初学者
1. 先运行 `simple-gorestful-hello.go` 理解基本概念
2. 修改参数和响应内容进行实验
3. 查看控制台输出的路由信息

### 进阶学习
1. 运行 `gorestful-hello-world.go` 学习完整实现
2. 添加新的接口和功能
3. 实现中间件和过滤器
4. 集成 Swagger 文档生成

### 与 Kubernetes 对比
- Kubernetes API Server 使用相同的 Go-Restful 框架
- 可以参考这些示例理解 Kubernetes API 的设计模式
- 学习如何实现 RESTful API 的最佳实践

## 🐛 常见问题

### 1. 端口被占用
```bash
# 查看端口占用
lsof -i :8080
# 杀死进程
kill -9 <PID>
```

### 2. 依赖安装失败
```bash
# 清理模块缓存
go clean -modcache
# 重新下载依赖
go mod download
```

### 3. JSON 解析错误
- 确保请求头包含 `Content-Type: application/json`
- 检查 JSON 格式是否正确
- 验证结构体字段名称是否匹配

---

*这些示例展示了 Go-Restful 框架的基本用法，帮助理解 Kubernetes API Server 的 HTTP 框架实现。*