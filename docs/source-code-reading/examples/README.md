# Kubernetes 代码示例集合

## 🎯 示例目标

提供实用的代码示例，帮助您：
- 快速上手 Kubernetes 相关开发
- 理解关键 API 的使用方法
- 掌握常见的开发模式
- 建立从理论到实践的桥梁

## 📁 示例分类

### 🔌 Client-Go 示例

#### 基础操作 (`client-go-examples/basic-operations/`)
```go
// basic-client.go - 基本客户端操作
package main

import (
    "context"
    "fmt"
    
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // 从 kubeconfig 创建客户端
    config, err := clientcmd.BuildConfigFromFlags("", "/path/to/kubeconfig")
    if err != nil {
        panic(err)
    }
    
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }
    
    // 列出 Pod
    pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        panic(err)
    }
    
    for _, pod := range pods.Items {
        fmt.Printf("Pod: %s\n", pod.Name)
    }
}
```

**包含示例**:
- [x] `basic-client.go` - 基本客户端创建和使用
- [x] `crud-operations.go` - CRUD 操作示例
- [x] `watch-events.go` - 监听资源变化
- [x] `field-selectors.go` - 字段选择器使用
- [x] `label-selectors.go` - 标签选择器使用

#### Informer 和 Controller (`client-go-examples/informer-controller/`)
```go
// informer-example.go - Informer 使用示例
package main

import (
    "time"
    
    v1 "k8s.io/api/core/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
)

func main() {
    clientset := getClientset() // 假设已有客户端
    
    // 创建 Informer Factory
    factory := informers.NewSharedInformerFactory(clientset, 30*time.Second)
    podInformer := factory.Core().V1().Pods()
    
    // 添加事件处理器
    podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            pod := obj.(*v1.Pod)
            fmt.Printf("Pod %s was created\n", pod.Name)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            pod := newObj.(*v1.Pod)
            fmt.Printf("Pod %s was updated\n", pod.Name)
        },
        DeleteFunc: func(obj interface{}) {
            pod := obj.(*v1.Pod)
            fmt.Printf("Pod %s was deleted\n", pod.Name)
        },
    })
    
    // 启动 Informer
    stopCh := make(chan struct{})
    factory.Start(stopCh)
    factory.WaitForCacheSync(stopCh)
    
    // 保持运行
    <-stopCh
}
```

**包含示例**:
- [x] `informer-example.go` - 基础 Informer 使用
- [x] `shared-informer.go` - 共享 Informer 使用
- [x] `workqueue-controller.go` - 工作队列控制器
- [x] `custom-controller.go` - 自定义控制器框架

#### 高级模式 (`client-go-examples/advanced-patterns/`)
**包含示例**:
- [x] `dynamic-client.go` - 动态客户端使用
- [x] `discovery-client.go` - 发现客户端
- [x] `rest-client.go` - REST 客户端直接使用
- [x] `leader-election.go` - 领导者选举
- [x] `rate-limiting.go` - 速率限制

### 🎛️ 控制器示例

#### 简单控制器 (`controller-examples/simple-controller/`)
```go
// simple-controller.go - 简单控制器示例
package main

import (
    "context"
    "fmt"
    "time"
    
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/fields"
    "k8s.io/apimachinery/pkg/util/runtime"
    "k8s.io/apimachinery/pkg/util/wait"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/util/workqueue"
)

type Controller struct {
    clientset      kubernetes.Interface
    deploymentLister cache.Indexer
    deploymentInformer cache.Controller
    queue          workqueue.RateLimitingInterface
}

func NewController(clientset kubernetes.Interface) *Controller {
    // 创建 Deployment 监听器
    deploymentListWatcher := cache.NewListWatchFromClient(
        clientset.AppsV1().RESTClient(),
        "deployments",
        corev1.NamespaceAll,
        fields.Everything(),
    )
    
    queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
    
    indexer, informer := cache.NewIndexerInformer(deploymentListWatcher, &appsv1.Deployment{}, 0, cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            key, err := cache.MetaNamespaceKeyFunc(obj)
            if err == nil {
                queue.Add(key)
            }
        },
        UpdateFunc: func(old interface{}, new interface{}) {
            key, err := cache.MetaNamespaceKeyFunc(new)
            if err == nil {
                queue.Add(key)
            }
        },
        DeleteFunc: func(obj interface{}) {
            key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
            if err == nil {
                queue.Add(key)
            }
        },
    }, cache.Indexers{})
    
    return &Controller{
        clientset:          clientset,
        deploymentLister:   indexer,
        deploymentInformer: informer,
        queue:              queue,
    }
}

func (c *Controller) Run(stopCh chan struct{}) {
    defer runtime.HandleCrash()
    defer c.queue.ShutDown()
    
    go c.deploymentInformer.Run(stopCh)
    
    if !cache.WaitForCacheSync(stopCh, c.deploymentInformer.HasSynced) {
        runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
        return
    }
    
    wait.Until(c.runWorker, time.Second, stopCh)
}

func (c *Controller) runWorker() {
    for c.processNextItem() {
    }
}

func (c *Controller) processNextItem() bool {
    key, quit := c.queue.Get()
    if quit {
        return false
    }
    defer c.queue.Done(key)
    
    err := c.syncToStdout(key.(string))
    c.handleErr(err, key)
    return true
}

func (c *Controller) syncToStdout(key string) error {
    obj, exists, err := c.deploymentLister.GetByKey(key)
    if err != nil {
        return err
    }
    
    if !exists {
        fmt.Printf("Deployment %s does not exist anymore\n", key)
    } else {
        deployment := obj.(*appsv1.Deployment)
        fmt.Printf("Sync/Add/Update for Deployment %s\n", deployment.GetName())
    }
    
    return nil
}

func (c *Controller) handleErr(err error, key interface{}) {
    if err == nil {
        c.queue.Forget(key)
        return
    }
    
    if c.queue.NumRequeues(key) < 5 {
        c.queue.AddRateLimited(key)
        return
    }
    
    c.queue.Forget(key)
    runtime.HandleError(err)
}
```

#### Operator 模式 (`controller-examples/operator-pattern/`)
**包含示例**:
- [x] `custom-resource-definition.yaml` - CRD 定义
- [x] `operator-controller.go` - Operator 控制器
- [x] `reconcile-logic.go` - 调和逻辑实现
- [x] `status-management.go` - 状态管理

#### Webhook 控制器 (`controller-examples/webhook-controller/`)
**包含示例**:
- [x] `mutating-webhook.go` - 变更 Webhook
- [x] `validating-webhook.go` - 验证 Webhook
- [x] `webhook-server.go` - Webhook 服务器

### 🔗 Webhook 示例

#### 准入 Webhook (`webhook-examples/admission-webhook/`)
```go
// admission-webhook.go - 准入 Webhook 示例
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    
    admissionv1 "k8s.io/api/admission/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
    scheme = runtime.NewScheme()
    codecs = serializer.NewCodecFactory(scheme)
)

func main() {
    http.HandleFunc("/mutate", mutateHandler)
    http.HandleFunc("/validate", validateHandler)
    
    fmt.Println("Webhook server starting...")
    if err := http.ListenAndServeTLS(":8443", "tls.crt", "tls.key", nil); err != nil {
        panic(err)
    }
}

func mutateHandler(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    admissionReview := admissionv1.AdmissionReview{}
    if err := json.Unmarshal(body, &admissionReview); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    pod := corev1.Pod{}
    if err := json.Unmarshal(admissionReview.Request.Object.Raw, &pod); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 添加标签的变更
    patch := []map[string]interface{}{
        {
            "op":    "add",
            "path":  "/metadata/labels/mutated",
            "value": "true",
        },
    }
    
    patchBytes, _ := json.Marshal(patch)
    
    admissionResponse := &admissionv1.AdmissionResponse{
        UID:     admissionReview.Request.UID,
        Allowed: true,
        Patch:   patchBytes,
        PatchType: func() *admissionv1.PatchType {
            pt := admissionv1.PatchTypeJSONPatch
            return &pt
        }(),
    }
    
    admissionReview.Response = admissionResponse
    respBytes, _ := json.Marshal(admissionReview)
    
    w.Header().Set("Content-Type", "application/json")
    w.Write(respBytes)
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
    // 验证逻辑实现
    // ...
}
```

### 📅 调度器插件

#### 自定义调度插件 (`scheduler-plugins/`)
**包含示例**:
- [x] `filter-plugin.go` - 过滤插件
- [x] `score-plugin.go` - 评分插件
- [x] `bind-plugin.go` - 绑定插件
- [x] `scheduler-config.yaml` - 调度器配置

### 🔧 工具和实用程序

#### 开发工具 (`tools/`)
**包含示例**:
- [x] `code-generator.sh` - 代码生成脚本
- [x] `mock-generator.go` - Mock 生成器
- [x] `test-helpers.go` - 测试辅助函数
- [x] `debugging-utils.go` - 调试工具

#### 性能测试 (`performance/`)
**包含示例**:
- [x] `benchmark-client.go` - 客户端性能测试
- [x] `load-testing.go` - 负载测试
- [x] `memory-profiling.go` - 内存分析

## 🚀 使用指南

### 环境准备
```bash
# 确保有可用的 Kubernetes 集群
kubectl cluster-info

# 安装依赖
go mod init k8s-examples
go get k8s.io/client-go@v0.28.0
go get k8s.io/api@v0.28.0
go get k8s.io/apimachinery@v0.28.0
```

### 运行示例
```bash
# 进入示例目录
cd examples/client-go-examples/basic-operations/

# 编译和运行
go build -o basic-client basic-client.go
./basic-client

# 或直接运行
go run basic-client.go
```

### 调试技巧
```bash
# 启用详细日志
KLOG_V=4 go run example.go

# 使用 delve 调试
dlv exec ./example
```

## 📝 示例规范

### 代码规范
- 使用清晰的变量和函数命名
- 添加必要的注释和文档
- 遵循 Go 语言惯例
- 包含错误处理

### 文档要求
每个示例包含：
- `README.md` - 功能说明和使用方法
- `main.go` - 主要代码
- `Makefile` - 构建脚本（如需要）
- `test/` - 测试代码（如适用）

### 测试标准
- 提供单元测试
- 包含集成测试
- 验证核心功能
- 性能基准测试

## 🎯 学习路径

### 初学者
1. 从 `basic-operations` 开始
2. 学习 `informer-controller` 基础
3. 尝试 `simple-controller` 实现

### 进阶开发者
1. 掌握 `advanced-patterns`
2. 实现 `operator-pattern`
3. 开发 `webhook-examples`

### 高级开发者
1. 贡献 `scheduler-plugins`
2. 优化 `performance` 示例
3. 开发自定义工具

## 🤝 贡献指南

### 添加新示例
1. 选择合适的分类目录
2. 创建清晰的文档
3. 提供完整的代码
4. 添加测试用例

### 改进现有示例
1. 修复错误和问题
2. 优化代码质量
3. 更新依赖版本
4. 改进文档说明

### 代码审查
- 确保代码可以正常运行
- 验证最佳实践的应用
- 检查文档的完整性
- 测试示例的有效性

---

*通过这些实用示例，快速掌握 Kubernetes 开发的精髓！*
