# Kubernetes 源码学习路径

## 🎯 学习目标设定

根据不同的目标，选择相应的学习路径：

### 🎓 基础理解型 (2-3个月)
- **目标**: 理解 K8s 整体架构和核心概念
- **适合**: 运维工程师、初级开发者
- **重点**: API 使用、基本原理

### 🚀 开发贡献型 (4-6个月)  
- **目标**: 能够阅读代码、修复 bug、提交 PR
- **适合**: 中级开发者、开源贡献者
- **重点**: 核心组件源码、开发流程

### 🏗️ 架构专家型 (6-12个月)
- **目标**: 深度理解设计思想，能够架构类似系统
- **适合**: 高级工程师、架构师
- **重点**: 分布式系统设计、性能优化

## 📚 阶段一：基础准备 (4-6周)

### Week 1-2: Go 语言强化

#### 必须掌握的 Go 特性
```go
// 1. 并发编程
func worker(jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        select {
        case result := <-processJob(job):
            results <- result
        case <-time.After(timeout):
            results <- Result{Error: "timeout"}
        }
    }
}

// 2. 接口和类型系统
type Controller interface {
    Run(stopCh <-chan struct{})
    HasSynced() bool
}

// 3. Context 和取消
func (c *Controller) Run(ctx context.Context) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-c.processingComplete:
        return nil
    }
}
```

#### 学习资源
- 📖 《The Go Programming Language》
- 💻 [Go by Example](https://gobyexample.com/)
- 🎯 实践：实现一个简单的 HTTP 服务器和客户端

### Week 3-4: 容器技术基础

#### Docker 深入理解
```bash
# 理解容器运行时
docker run --rm -it \
  --pid=host \
  --net=host \
  --privileged \
  alpine:latest

# 查看容器内核名字空间
lsns -t pid,net,mnt,user
```

#### 学习重点
- Linux 容器技术 (namespaces, cgroups)
- 容器镜像格式和存储
- 容器网络基础
- 容器运行时 (containerd, CRI-O)

#### 实践项目
```go
// 简单的容器管理器
type ContainerManager interface {
    CreateContainer(config *ContainerConfig) (string, error)
    StartContainer(id string) error
    StopContainer(id string) error
    ListContainers() ([]*Container, error)
}
```

### Week 5-6: Kubernetes 基础

#### 集群搭建和基本操作
```bash
# 使用 kubeadm 搭建集群
kubeadm init --pod-network-cidr=192.168.0.0/16

# 基本 kubectl 操作
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --port=80 --type=NodePort
kubectl get pods -o wide
```

#### 核心概念理解
- Pod、Service、Deployment
- ConfigMap、Secret
- Namespace、Label、Selector
- RBAC 权限控制

## 🔍 阶段二：源码入门 (6-8周)

### Week 7-8: client-go 深入

#### 基础客户端使用
```go
package main

import (
    "context"
    "fmt"
    
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    // 加载 kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", "/path/to/kubeconfig")
    if err != nil {
        panic(err)
    }
    
    // 创建客户端
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

#### Informer 和 Controller 模式
```go
import (
    "k8s.io/client-go/informers"
    "k8s.io/client-go/tools/cache"
)

func setupInformer(clientset kubernetes.Interface) {
    factory := informers.NewSharedInformerFactory(clientset, 0)
    podInformer := factory.Core().V1().Pods()
    
    podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            pod := obj.(*v1.Pod)
            fmt.Printf("Pod added: %s\n", pod.Name)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            pod := newObj.(*v1.Pod)
            fmt.Printf("Pod updated: %s\n", pod.Name)
        },
        DeleteFunc: func(obj interface{}) {
            pod := obj.(*v1.Pod)
            fmt.Printf("Pod deleted: %s\n", pod.Name)
        },
    })
    
    factory.Start(stopCh)
    factory.WaitForCacheSync(stopCh)
}
```

### Week 9-10: kubectl 源码分析

#### 命令结构分析
```go
// staging/src/k8s.io/kubectl/pkg/cmd/cmd.go
func NewKubectlCommand(in io.Reader, out, err io.Writer) *cobra.Command {
    cmds := &cobra.Command{
        Use:   "kubectl",
        Short: "kubectl controls the Kubernetes cluster manager",
    }
    
    // 添加子命令
    cmds.AddCommand(create.NewCmdCreate(f, ioStreams))
    cmds.AddCommand(get.NewCmdGet("kubectl", f, ioStreams))
    cmds.AddCommand(apply.NewCmdApply("kubectl", f, ioStreams))
    
    return cmds
}
```

#### 重点分析文件
- `staging/src/k8s.io/kubectl/pkg/cmd/get/get.go`
- `staging/src/k8s.io/kubectl/pkg/cmd/apply/apply.go`
- `staging/src/k8s.io/cli-runtime/pkg/resource/builder.go`

### Week 11-12: API Machinery

#### 类型系统
```go
// k8s.io/apimachinery/pkg/runtime
type Object interface {
    GetObjectKind() schema.ObjectKind
    DeepCopyObject() Object
}

// 序列化和反序列化
type Codec interface {
    Encoder
    Decoder
}
```

#### Scheme 和版本转换
```go
// 注册 API 类型
func addKnownTypes(scheme *runtime.Scheme) error {
    scheme.AddKnownTypes(SchemeGroupVersion,
        &Pod{},
        &PodList{},
        &Service{},
        &ServiceList{},
    )
    return nil
}
```

### Week 13-14: kube-apiserver 入门

#### API 服务器启动流程
```go
// cmd/kube-apiserver/app/server.go
func Run(completeOptions completedServerRunOptions, stopCh <-chan struct{}) error {
    server, err := CreateServerChain(completeOptions, stopCh)
    if err != nil {
        return err
    }
    
    return server.PrepareRun().Run(stopCh)
}
```

#### REST API 处理
```go
// pkg/registry/core/pod/storage/storage.go
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
    pod := obj.(*api.Pod)
    
    // 验证
    if createValidation != nil {
        if err := createValidation(ctx, obj); err != nil {
            return nil, err
        }
    }
    
    // 存储到 etcd
    return r.store.Create(ctx, pod, options)
}
```

## 🏗️ 阶段三：核心组件深入 (8-12周)

### Week 15-18: Controller Manager

#### Deployment Controller
```go
// pkg/controller/deployment/deployment_controller.go
func (dc *DeploymentController) syncDeployment(key string) error {
    namespace, name, err := cache.SplitMetaNamespaceKey(key)
    if err != nil {
        return err
    }
    
    deployment, err := dc.dLister.Deployments(namespace).Get(name)
    if err != nil {
        return err
    }
    
    // 获取关联的 ReplicaSet
    rsList, err := dc.getReplicaSetsForDeployment(deployment)
    if err != nil {
        return err
    }
    
    // 同步期望状态
    return dc.sync(deployment, rsList)
}
```

#### 学习重点
- `pkg/controller/deployment/`
- `pkg/controller/replicaset/`
- `pkg/controller/job/`
- `pkg/controller/namespace/`

### Week 19-22: Scheduler

#### 调度框架
```go
// pkg/scheduler/framework/interface.go
type Framework interface {
    PreFilter(ctx context.Context, state *CycleState, pod *v1.Pod) (*PreFilterResult, *Status)
    Filter(ctx context.Context, state *CycleState, pod *v1.Pod, nodeInfo *NodeInfo) *Status
    Score(ctx context.Context, state *CycleState, pod *v1.Pod, nodeName string) (int64, *Status)
    Bind(ctx context.Context, state *CycleState, pod *v1.Pod, nodeName string) *Status
}
```

#### 调度算法
```go
// pkg/scheduler/core/generic_scheduler.go
func (g *genericScheduler) Schedule(ctx context.Context, prof *profile.Profile, state *framework.CycleState, pod *v1.Pod) (result ScheduleResult, err error) {
    // 1. 预过滤
    preFilterResult, s := prof.RunPreFilterPlugins(ctx, state, pod)
    if !s.IsSuccess() {
        return result, s.AsError()
    }
    
    // 2. 过滤节点
    feasibleNodes, filteredNodesStatuses, s := g.findNodesThatFitPod(ctx, prof, state, pod)
    if !s.IsSuccess() {
        return result, s.AsError()
    }
    
    // 3. 打分选择
    priorityList, s := g.prioritizeNodes(ctx, prof, state, pod, feasibleNodes)
    if !s.IsSuccess() {
        return result, s.AsError()
    }
    
    // 4. 选择最优节点
    host, err := g.selectHost(priorityList)
    return ScheduleResult{SuggestedHost: host}, err
}
```

### Week 23-26: Kubelet

#### Pod 管理
```go
// pkg/kubelet/kubelet.go
func (kl *Kubelet) syncPod(o syncPodOptions) error {
    pod := o.pod
    
    // 创建 Pod 目录
    if err := kl.makePodDataDirs(pod); err != nil {
        return err
    }
    
    // 拉取镜像
    if err := kl.imagePuller.EnsureImageExists(pod); err != nil {
        return err
    }
    
    // 启动容器
    result := kl.containerRuntime.SyncPod(pod, o.apiPodStatus, o.podStatus, o.pullSecrets, o.backOff)
    
    return result.Error()
}
```

#### 容器运行时集成
```go
// pkg/kubelet/kuberuntime/kuberuntime_manager.go
func (m *kubeGenericRuntimeManager) SyncPod(pod *v1.Pod, podStatus *kubecontainer.PodStatus, pullSecrets []v1.Secret, backOff *flowcontrol.Backoff) (result kubecontainer.PodSyncResult) {
    // 创建 Pod Sandbox
    podSandboxID, err := m.createPodSandbox(pod, backOff)
    if err != nil {
        return result
    }
    
    // 启动 Init 容器
    if err := m.startInitContainers(pod, podSandboxID); err != nil {
        return result
    }
    
    // 启动应用容器
    for _, container := range pod.Spec.Containers {
        if err := m.startContainer(pod, &container, podSandboxID); err != nil {
            return result
        }
    }
    
    return result
}
```

## 🚀 阶段四：高级特性 (6-8周)

### Week 27-28: 网络和存储

#### CNI 插件系统
```go
// pkg/kubelet/dockershim/network/cni/cni.go
func (plugin *cniNetworkPlugin) addToNetwork(network *cniNetwork, podName string, podNamespace string, podSandboxID kubecontainer.ContainerID) error {
    rt := &libcni.RuntimeConf{
        ContainerID: podSandboxID.ID,
        NetNS:       netnsPath,
    }
    
    return plugin.cniConfig.AddNetworkList(context.Background(), network.NetworkConfig, rt)
}
```

#### CSI 存储插件
```go
// pkg/volume/csi/csi_attacher.go
func (c *csiAttacher) Attach(spec *volume.Spec, nodeName types.NodeName) (string, error) {
    csiSource, err := getCSISourceFromSpec(spec)
    if err != nil {
        return "", err
    }
    
    // 调用 CSI Controller 接口
    return c.k8sAttach(spec, nodeName, csiSource)
}
```

### Week 29-30: 扩展性和自定义资源

#### Custom Resource Definition
```go
// staging/src/k8s.io/apiextensions-apiserver/pkg/controller/openapi/controller.go
func (c *Controller) sync(key string) error {
    crd, err := c.crdLister.Get(key)
    if err != nil {
        return err
    }
    
    // 生成 OpenAPI Schema
    spec, err := c.buildOpenAPISpec(crd)
    if err != nil {
        return err
    }
    
    // 更新 API Server
    return c.updateOpenAPISpec(spec)
}
```

#### Operator 模式
```go
// 自定义控制器示例
type MyController struct {
    client    clientset.Interface
    myLister  listers.MyResourceLister
    workqueue workqueue.RateLimitingInterface
}

func (c *MyController) processNextWorkItem() bool {
    obj, shutdown := c.workqueue.Get()
    if shutdown {
        return false
    }
    defer c.workqueue.Done(obj)
    
    err := c.syncHandler(obj.(string))
    if err == nil {
        c.workqueue.Forget(obj)
        return true
    }
    
    c.workqueue.AddRateLimited(obj)
    return true
}
```

### Week 31-32: 安全和监控

#### RBAC 实现
```go
// pkg/registry/rbac/role/storage/storage.go
func (r *REST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error) {
    role := obj.(*rbac.Role)
    
    // RBAC 验证
    if err := rbacvalidation.ValidateRole(role); err != nil {
        return nil, err
    }
    
    return r.store.Create(ctx, role, options)
}
```

#### 审计日志
```go
// staging/src/k8s.io/apiserver/pkg/audit/event.go
func LogRequestObject(ctx context.Context, obj runtime.Object, verb string, level auditinternal.Level) {
    ac := AuditContextFrom(ctx)
    if ac == nil || ac.Event.Level.Less(level) {
        return
    }
    
    // 记录审计事件
    ac.Event.ObjectRef = GetObjectRef(obj)
    ac.Event.Verb = verb
}
```

## 📈 学习进度检查

### 阶段性目标检查

#### 基础阶段完成标志
- [ ] 能够独立搭建 Kubernetes 集群
- [ ] 熟练使用 kubectl 和基本 API
- [ ] 理解 Pod、Service、Deployment 等核心概念
- [ ] 能够编写简单的 client-go 程序

#### 进阶阶段完成标志  
- [ ] 能够阅读和理解 kubectl 源码
- [ ] 掌握 Informer 和 Controller 模式
- [ ] 理解 API Server 的基本工作原理
- [ ] 能够编写简单的 Controller

#### 高级阶段完成标志
- [ ] 深入理解各核心组件的设计思想
- [ ] 能够分析和解决复杂的源码问题
- [ ] 具备贡献代码到上游的能力
- [ ] 能够设计和实现类似的分布式系统

### 实践项目建议

#### 入门项目
```go
// 1. Pod 监控器
// 监控集群中 Pod 的状态变化
func main() {
    // 使用 client-go 监控 Pod 事件
    // 输出 Pod 创建、更新、删除日志
}
```

#### 进阶项目
```go
// 2. 自定义控制器
// 实现一个管理自定义资源的控制器
type MyController struct {
    // 控制器实现
}

// 3. Kubectl 插件
// 扩展 kubectl 功能
func main() {
    // 实现 kubectl-my-plugin
}
```

#### 高级项目
```go
// 4. Admission Webhook
// 实现准入控制逻辑
func admissionHandler(w http.ResponseWriter, r *http.Request) {
    // Webhook 实现
}

// 5. Scheduler 插件
// 自定义调度策略
type MySchedulerPlugin struct{}

func (p *MySchedulerPlugin) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *framework.NodeInfo) *framework.Status {
    // 自定义过滤逻辑
}
```

## 📚 推荐学习资源

### 官方文档
- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [Kubernetes 开发指南](https://github.com/kubernetes/community/tree/master/contributors/devel)
- [API 参考文档](https://kubernetes.io/docs/reference/)

### 书籍推荐
- 《Kubernetes in Action》
- 《Programming Kubernetes》  
- 《Kubernetes: Up and Running》

### 在线资源
- [Kubernetes 源码分析](https://github.com/kubernetes/kubernetes)
- [client-go 示例](https://github.com/kubernetes/client-go/tree/master/examples)
- [Kubernetes 博客](https://kubernetes.io/blog/)

### 社区参与
- [Kubernetes Slack](https://slack.k8s.io/)
- [SIG 兴趣小组](https://github.com/kubernetes/community/blob/master/sig-list.md)
- [贡献指南](https://kubernetes.io/docs/contribute/)

---

*下一步：开始实践，从 [开发环境搭建](./dev-setup.md) 开始！*
