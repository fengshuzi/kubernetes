package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emicklei/go-restful/v3"
)

// HelloResponse 响应结构体
type HelloResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	// 创建 Go-Restful 容器
	container := restful.NewContainer()

	// 创建 WebService
	ws := new(restful.WebService)
	ws.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// 定义 Hello World 处理器
	helloHandler := func(request *restful.Request, response *restful.Response) {
		// 获取查询参数 name，默认为 "World"
		name := request.QueryParameter("name")
		if name == "" {
			name = "World"
		}

		// 创建响应
		helloResp := HelloResponse{
			Message: fmt.Sprintf("Hello, %s!", name),
			Status:  "success",
		}

		// 返回 JSON 响应
		response.WriteAsJson(helloResp)
	}

	// 注册路由
	ws.Route(ws.GET("/hello").
		To(helloHandler).
		Doc("Get hello message").
		Param(ws.QueryParameter("name", "Name to greet").DataType("string")).
		Returns(200, "OK", HelloResponse{}))

	// 将 WebService 添加到容器
	container.Add(ws)

	// 启动服务器
	port := ":8080"
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("  GET /api/v1/hello")
	fmt.Println("  GET /api/v1/hello?name=YourName")
	fmt.Println("")
	fmt.Println("🔗 Test with curl:")
	fmt.Println("  curl http://localhost:8080/api/v1/hello")
	fmt.Println("  curl http://localhost:8080/api/v1/hello?name=Kubernetes")
	fmt.Println("")

	// 启动 HTTP 服务器
	log.Fatal(http.ListenAndServe(port, container))
}
