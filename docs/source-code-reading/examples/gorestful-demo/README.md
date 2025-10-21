# Go-Restful Hello World 示例

这是一个使用 Go-Restful 框架创建 HTTP 接口的简单示例，返回 "Hello World" 消息。

## 🚀 快速开始

### 1. 安装依赖

```bash
cd gorestful-demo
go mod tidy
```

### 2. 运行简单示例

```bash
go run main.go
```

### 3. 运行高级示例

```bash
go run advanced-main.go
```

### 4. 测试接口

#### 手动测试
```bash
# 基本测试
curl http://localhost:8080/api/v1/hello

# 带参数测试
curl http://localhost:8080/api/v1/hello?name=Kubernetes
```

#### 自动测试脚本
```bash
# 运行测试脚本
./run-tests.sh
```

测试脚本支持：
- 简单示例测试
- 高级示例测试  
- 全部测试
- 自动检测服务器状态

## 📋 API 接口

### 简单示例接口

| 方法 | 路径 | 参数 | 说明 |
|------|------|------|------|
| GET | `/api/v1/hello` | `name` (可选) | 返回 Hello World 消息 |

### 高级示例接口

| 方法 | 路径 | 参数 | 说明 |
|------|------|------|------|
| GET | `/api/v1/hello` | `name` (查询参数) | 返回 Hello 消息 |
| POST | `/api/v1/hello` | `name`, `age` (请求体) | 创建用户并返回 Hello 消息 |
| GET | `/api/v1/users` | 无 | 获取所有用户列表 |
| GET | `/api/v1/users/{id}` | `id` (路径参数) | 根据 ID 获取用户 |
| GET | `/api/v1/health` | 无 | 健康检查 |
| GET | `/` | 无 | 欢迎页面 |

## 🔗 测试示例

### 简单示例测试

#### 基本请求
```bash
curl http://localhost:8080/api/v1/hello
```

**响应:**
```json
{
  "message": "Hello, World!",
  "status": "success"
}
```

#### 带参数请求
```bash
curl http://localhost:8080/api/v1/hello?name=Kubernetes
```

**响应:**
```json
{
  "message": "Hello, Kubernetes!",
  "status": "success"
}
```

### 高级示例测试

#### 1. GET Hello 接口
```bash
curl http://localhost:8080/api/v1/hello?name=Kubernetes
```

**响应:**
```json
{
  "message": "Hello, Kubernetes!",
  "status": "success",
  "time": "2024-12-19 10:30:00"
}
```

#### 2. POST Hello 接口
```bash
curl -X POST http://localhost:8080/api/v1/hello \
  -H 'Content-Type: application/json' \
  -d '{"name":"Kubernetes","age":25}'
```

**响应:**
```json
{
  "id": 1,
  "name": "Kubernetes",
  "age": 25,
  "message": "Hello, Kubernetes! (POST)",
  "status": "success"
}
```

#### 3. 获取用户列表
```bash
curl http://localhost:8080/api/v1/users
```

**响应:**
```json
{
  "users": [
    {
      "id": 1,
      "name": "Kubernetes",
      "age": 25,
      "message": "Hello, Kubernetes! (POST)",
      "status": "success"
    }
  ],
  "count": 1,
  "status": "success"
}
```

#### 4. 根据 ID 获取用户
```bash
curl http://localhost:8080/api/v1/users/1
```

**响应:**
```json
{
  "id": 1,
  "name": "Kubernetes",
  "age": 25,
  "message": "Hello, Kubernetes! (POST)",
  "status": "success"
}
```

#### 5. 健康检查
```bash
curl http://localhost:8080/api/v1/health
```

**响应:**
```json
{
  "status": "healthy",
  "service": "hello-world-service",
  "version": "1.0.0",
  "timestamp": "2024-12-19T10:30:00Z"
}
```

## 🏗️ 代码结构

### 主要组件

1. **Go-Restful 容器** - 管理所有 WebService
2. **WebService** - 定义 API 路径和版本
3. **路由** - 具体的 API 端点定义
4. **处理器** - 处理请求的业务逻辑

### 关键代码

```go
// 创建容器
container := restful.NewContainer()

// 创建 WebService
ws := new(restful.WebService)
ws.Path("/api/v1")

// 定义处理器
helloHandler := func(request *restful.Request, response *restful.Response) {
    // 处理逻辑
}

// 注册路由
ws.Route(ws.GET("/hello").To(helloHandler))
```

## 🎯 与 Kubernetes 的关系

这个示例展示了 Kubernetes API Server 使用的相同框架：

- **相同的框架**: Go-Restful v3
- **相同的模式**: WebService + 路由注册
- **相同的特性**: JSON 绑定、参数验证、自动路由

通过这个简单示例，你可以理解 Kubernetes API Server 的 HTTP 框架基础。

## 🔧 扩展功能

你可以基于这个示例添加更多功能：

1. **POST 接口** - 处理 POST 请求
2. **路径参数** - 使用 `{id}` 格式的参数
3. **中间件** - 添加认证、日志等中间件
4. **错误处理** - 统一的错误响应格式
5. **Swagger 文档** - 自动生成 API 文档

---

*这个示例帮助你理解 Kubernetes API Server 的 HTTP 框架基础。*
