# Kubernetes 技术栈详解

## 🔧 编程语言和基础工具

### Go 语言 (v1.24+)

Kubernetes 几乎完全使用 Go 语言编写，需要掌握以下高级特性：

#### 并发编程
```go
// Goroutine 和 Channel 广泛使用
func (c *Controller) Run(stopCh <-chan struct{}) {
    defer utilruntime.HandleCrash()
    defer c.workqueue.ShutDown()
    
    go wait.Until(c.runWorker, time.Second, stopCh)
    <-stopCh
}
```

#### 接口和反射
```go
// runtime.Object 是核心接口
type Object interface {
    GetObjectKind() schema.ObjectKind
    DeepCopyObject() Object
}

// 反射用于序列化和验证
func (s *Scheme) Convert(in, out interface{}, context interface{}) error {
    // 使用 reflect 包进行类型转换
}
```

#### Context 管理
```go
// 贯穿整个请求生命周期
func (r *Request) Do(ctx context.Context) Result {
    if ctx == nil {
        ctx = context.Background()
    }
    // ...
}
```

### 构建工具

#### Make
```makefile
# 主要构建命令
all: hack/make-rules/build.sh $(WHAT)
test: hack/make-rules/test.sh $(WHAT) $(TESTS)
verify: hack/make-rules/verify.sh
```

#### Bash 脚本
```bash
# hack/make-rules/build.sh - 核心构建脚本
KUBE_BUILD_PLATFORMS=${KUBE_BUILD_PLATFORMS:-linux/amd64}
```

## 🌐 网络和通信框架

### gRPC + Protocol Buffers

#### 定义服务
```proto
// 容器运行时接口 (CRI)
service RuntimeService {
    rpc Version(VersionRequest) returns (VersionResponse) {}
    rpc CreateContainer(CreateContainerRequest) returns (CreateContainerResponse) {}
}
```

#### 客户端调用
```go
// kubelet 与容器运行时通信
conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
client := runtimeapi.NewRuntimeServiceClient(conn)
```

### REST API

#### Cobra CLI 框架
```go
// kubectl 命令定义
var getCmd = &cobra.Command{
    Use:   "get [(-o|--output=)json|yaml|name|...]",
    Short: "Display one or many resources",
    Run: func(cmd *cobra.Command, args []string) {
        // 命令实现
    },
}
```

#### HTTP 路由
```go
// API 路由注册
func (s *APIGroupVersion) InstallREST(container *restful.Container) error {
    for _, path := range paths {
        route := ws.GET(path).To(handler)
        ws.Route(route)
    }
}
```

## 🗄️ 存储和数据

### etcd 分布式存储

#### 客户端配置
```go
import "go.etcd.io/etcd/client/v3"

cfg := clientv3.Config{
    Endpoints:   []string{"localhost:2379"},
    DialTimeout: 5 * time.Second,
}
client, err := clientv3.New(cfg)
```

#### 键值操作
```go
// 存储 Kubernetes 对象
key := "/registry/pods/default/my-pod"
value := serializedPod
_, err := client.Put(context.Background(), key, value)
```

### 序列化格式

#### JSON/YAML
```go
// runtime.Codec 接口
type Codec interface {
    Encoder
    Decoder
}

// YAML 解析
obj := &v1.Pod{}
err := yaml.Unmarshal(data, obj)
```

#### Protocol Buffers
```go
// 高效的二进制序列化
import "k8s.io/apimachinery/pkg/runtime/serializer/protobuf"

codec := protobuf.NewCodec(scheme)
```

## 🔌 接口和插件系统

### Container Runtime Interface (CRI)

#### 运行时实现
```go
// 容器运行时接口
type RuntimeService interface {
    Version(apiVersion string) (*RuntimeVersion, error)
    CreateContainer(config *ContainerConfig, sandboxConfig *PodSandboxConfig) (string, error)
    StartContainer(containerID string) error
}
```

#### 支持的运行时
- **containerd** - 主流容器运行时
- **CRI-O** - 轻量级运行时
- **Docker** - 传统支持（已弃用）

### Container Network Interface (CNI)

#### 网络插件配置
```json
{
  "cniVersion": "0.4.0",
  "name": "bridge",
  "type": "bridge",
  "bridge": "cni0",
  "isGateway": true,
  "ipMasq": true,
  "ipam": {
    "type": "host-local",
    "subnet": "10.244.0.0/16"
  }
}
```

#### 网络管理
```go
// CNI 插件调用
func (plugin *cniNetworkPlugin) addToNetwork(network *cniNetwork, podName string, podNamespace string, podSandboxID kubecontainer.ContainerID, podNetnsPath string) error {
    rt := &libcni.RuntimeConf{
        ContainerID: podSandboxID.ID,
        NetNS:       podNetnsPath,
    }
    
    return plugin.cniConfig.AddNetworkList(context.Background(), network.NetworkConfig, rt)
}
```

### Container Storage Interface (CSI)

#### 存储驱动
```go
// CSI 驱动接口
type ControllerServer interface {
    CreateVolume(context.Context, *CreateVolumeRequest) (*CreateVolumeResponse, error)
    DeleteVolume(context.Context, *DeleteVolumeRequest) (*DeleteVolumeResponse, error)
}
```

## 🧪 测试框架

### Ginkgo BDD 测试

#### 测试套件定义
```go
import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("Pod", func() {
    Context("when creating a pod", func() {
        It("should be scheduled successfully", func() {
            pod := &v1.Pod{...}
            Expect(client.Create(ctx, pod)).To(Succeed())
        })
    })
})
```

### Gomega 断言库

#### 常用断言
```go
// 基础断言
Expect(err).NotTo(HaveOccurred())
Expect(pod.Status.Phase).To(Equal(v1.PodRunning))

// 异步断言
Eventually(func() bool {
    return pod.Status.Phase == v1.PodRunning
}).Should(BeTrue())
```

### 测试工具

#### 模拟对象
```go
// fake 客户端
import "k8s.io/client-go/kubernetes/fake"

client := fake.NewSimpleClientset()
```

## 🔐 安全框架

### 认证 (Authentication)

#### 多种认证方式
```go
// X509 证书认证
type x509Authenticator struct {
    roots *x509.CertPool
}

// OIDC 认证
type oidcAuthenticator struct {
    issuerURL string
    clientID  string
}
```

### 授权 (Authorization)

#### RBAC (Role-Based Access Control)
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "watch", "list"]
```

#### 授权检查
```go
// 授权决策
func (r *RBACAuthorizer) Authorize(ctx context.Context, a authorizer.Attributes) (authorizer.Decision, string, error) {
    // 检查 RBAC 规则
}
```

### 准入控制 (Admission Control)

#### Admission Webhooks
```go
// 准入控制器接口
type Interface interface {
    Admit(ctx context.Context, a Attributes, o ObjectInterfaces) error
}

// Validating Admission Webhook
type ValidatingAdmissionWebhook struct {
    hookSource   hookSource
    namespaceMatcher NamespaceMatcher
}
```

## 📊 监控和可观测性

### Prometheus 指标

#### 指标定义
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    requestCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "apiserver_request_total",
            Help: "Counter of apiserver requests broken out for each verb, API resource, client, and HTTP response contentType and code.",
        },
        []string{"verb", "resource", "client", "contentType", "code"},
    )
)
```

#### 指标收集
```go
// 中间件记录指标
func WithMetrics(handler http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        handler.ServeHTTP(w, r)
        duration := time.Since(start)
        requestDuration.Observe(duration.Seconds())
    })
}
```

### OpenTelemetry 追踪

#### 分布式追踪
```go
import "go.opentelemetry.io/otel/trace"

func (r *Request) Do(ctx context.Context) Result {
    ctx, span := tracer.Start(ctx, "http-request")
    defer span.End()
    
    // 执行请求
}
```

### 日志框架

#### klog 结构化日志
```go
import "k8s.io/klog/v2"

// 结构化日志
klog.InfoS("Pod created", "pod", klog.KObj(pod), "node", nodeName)

// 错误日志
klog.ErrorS(err, "Failed to create pod", "pod", klog.KObj(pod))
```

## 🔧 开发和调试工具

### 代码生成

#### client-gen
```bash
# 生成客户端代码
go run k8s.io/code-generator/cmd/client-gen \
  --clientset-name versioned \
  --input-base "" \
  --input k8s.io/api/core/v1
```

#### informer-gen
```bash
# 生成 Informer 代码
go run k8s.io/code-generator/cmd/informer-gen \
  --versioned-clientset-package k8s.io/client-go/kubernetes \
  --listers-package k8s.io/client-go/listers
```

### 性能分析

#### pprof 性能分析
```go
import _ "net/http/pprof"

// 启动 pprof 服务
go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

#### 内存和 CPU 分析
```bash
# CPU 分析
go tool pprof http://localhost:6060/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap
```

## 🏗️ 架构模式

### 控制器模式

#### Controller 接口
```go
type Controller interface {
    Run(stopCh <-chan struct{})
}

// 控制循环
func (c *Controller) processNextWorkItem() bool {
    obj, shutdown := c.workqueue.Get()
    if shutdown {
        return false
    }
    defer c.workqueue.Done(obj)
    
    err := c.syncHandler(obj.(string))
    c.handleErr(err, obj)
    return true
}
```

### Informer 模式

#### 事件监听
```go
// SharedInformer 监听资源变化
podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
    AddFunc: func(obj interface{}) {
        controller.enqueue(obj)
    },
    UpdateFunc: func(oldObj, newObj interface{}) {
        controller.enqueue(newObj)
    },
    DeleteFunc: func(obj interface{}) {
        controller.enqueue(obj)
    },
})
```

---

*下一步：阅读 [学习路径](./learning-path.md) 制定系统化的学习计划*
