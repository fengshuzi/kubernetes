package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful/v3"
)

// HelloResponse 响应结构体
type HelloResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Time    string `json:"time"`
}

// UserRequest 用户请求结构体
type UserRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// UserResponse 用户响应结构体
type UserResponse struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string `json:"status"`
	Service   string `json:"service"`
	Version   string `json:"version"`
	Timestamp string `json:"timestamp"`
}

// HelloWorldService 服务结构体
type HelloWorldService struct {
	users  []UserResponse
	nextID int
}

// NewHelloWorldService 创建服务实例
func NewHelloWorldService() *HelloWorldService {
	return &HelloWorldService{
		users:  make([]UserResponse, 0),
		nextID: 1,
	}
}

// GetHello 处理 GET /hello 请求
func (h *HelloWorldService) GetHello(request *restful.Request, response *restful.Response) {
	name := request.QueryParameter("name")
	if name == "" {
		name = "World"
	}

	helloResp := HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", name),
		Status:  "success",
		Time:    time.Now().Format("2006-01-02 15:04:05"),
	}

	response.WriteAsJson(helloResp)
}

// PostHello 处理 POST /hello 请求
func (h *HelloWorldService) PostHello(request *restful.Request, response *restful.Response) {
	var userReq UserRequest
	err := request.ReadEntity(&userReq)
	if err != nil {
		response.WriteErrorString(http.StatusBadRequest, "Invalid JSON format")
		return
	}

	if userReq.Name == "" {
		response.WriteErrorString(http.StatusBadRequest, "Name is required")
		return
	}

	userResp := UserResponse{
		ID:      h.nextID,
		Name:    userReq.Name,
		Age:     userReq.Age,
		Message: fmt.Sprintf("Hello, %s! (POST)", userReq.Name),
		Status:  "success",
	}

	h.users = append(h.users, userResp)
	h.nextID++

	response.WriteAsJson(userResp)
}

// GetUsers 处理 GET /users 请求
func (h *HelloWorldService) GetUsers(request *restful.Request, response *restful.Response) {
	response.WriteAsJson(map[string]interface{}{
		"users":  h.users,
		"count":  len(h.users),
		"status": "success",
	})
}

// GetUser 处理 GET /users/{id} 请求
func (h *HelloWorldService) GetUser(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")

	// 简单的 ID 查找逻辑
	for _, user := range h.users {
		if fmt.Sprintf("%d", user.ID) == id {
			response.WriteAsJson(user)
			return
		}
	}

	response.WriteErrorString(http.StatusNotFound, "User not found")
}

// GetHealth 健康检查
func (h *HelloWorldService) GetHealth(request *restful.Request, response *restful.Response) {
	healthResp := HealthResponse{
		Status:    "healthy",
		Service:   "hello-world-service",
		Version:   "1.0.0",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	response.WriteAsJson(healthResp)
}

// RegisterRoutes 注册所有路由
func (h *HelloWorldService) RegisterRoutes(container *restful.Container) {
	// 创建 WebService
	ws := new(restful.WebService)
	ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// 注册路由
	ws.Route(ws.GET("/hello").
		To(h.GetHello).
		Doc("Get hello message").
		Param(ws.QueryParameter("name", "Name to greet").DataType("string")).
		Returns(200, "OK", HelloResponse{}))

	ws.Route(ws.POST("/hello").
		To(h.PostHello).
		Doc("Post hello message").
		Reads(UserRequest{}).
		Returns(200, "OK", UserResponse{}).
		Returns(400, "Bad Request", nil))

	ws.Route(ws.GET("/users").
		To(h.GetUsers).
		Doc("Get all users").
		Returns(200, "OK", map[string]interface{}{}))

	ws.Route(ws.GET("/users/{id}").
		To(h.GetUser).
		Doc("Get user by ID").
		Param(ws.PathParameter("id", "User ID").DataType("string")).
		Returns(200, "OK", UserResponse{}).
		Returns(404, "Not Found", nil))

	ws.Route(ws.GET("/health").
		To(h.GetHealth).
		Doc("Health check").
		Returns(200, "OK", HealthResponse{}))

	// 添加到容器
	container.Add(ws)
}

func main() {
	// 创建 Go-Restful 容器
	container := restful.NewContainer()

	// 设置路由器
	container.Router(restful.CurlyRouter{})

	// 添加 CORS 支持
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedDomains: []string{"*"},
		CookiesAllowed: false,
		Container:      container}
	container.Filter(cors.Filter)

	// 添加请求日志过滤器
	container.Filter(restful.OPTIONSFilter())

	// 创建服务实例并注册路由
	helloService := NewHelloWorldService()
	helloService.RegisterRoutes(container)

	// 添加根路径处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":   "Welcome to Hello World API",
			"version":   "1.0.0",
			"endpoints": "GET /api/v1/hello, POST /api/v1/hello, GET /api/v1/users, GET /api/v1/users/{id}, GET /api/v1/health",
		})
	})

	// 启动服务器
	port := ":8080"
	fmt.Printf("🚀 Advanced Server starting on port %s\n", port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("  GET  /api/v1/hello?name=YourName")
	fmt.Println("  POST /api/v1/hello")
	fmt.Println("  GET  /api/v1/users")
	fmt.Println("  GET  /api/v1/users/{id}")
	fmt.Println("  GET  /api/v1/health")
	fmt.Println("  GET  /")
	fmt.Println("")
	fmt.Println("🔗 Test with curl:")
	fmt.Println("  curl http://localhost:8080/api/v1/hello")
	fmt.Println("  curl http://localhost:8080/api/v1/hello?name=Kubernetes")
	fmt.Println("  curl -X POST http://localhost:8080/api/v1/hello -H 'Content-Type: application/json' -d '{\"name\":\"Kubernetes\",\"age\":25}'")
	fmt.Println("  curl http://localhost:8080/api/v1/users")
	fmt.Println("  curl http://localhost:8080/api/v1/users/1")
	fmt.Println("  curl http://localhost:8080/api/v1/health")
	fmt.Println("")

	// 启动 HTTP 服务器
	log.Fatal(http.ListenAndServe(port, container))
}
