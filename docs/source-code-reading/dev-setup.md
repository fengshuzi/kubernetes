# Kubernetes 开发环境搭建指南

## 🎯 环境要求

### 系统要求
- **操作系统**: Linux (推荐 Ubuntu 20.04+) 或 macOS
- **内存**: 至少 16GB (推荐 32GB)
- **磁盘**: 至少 100GB 可用空间
- **CPU**: 8核心以上 (构建会很慢否则)

### 必要软件

#### Go 环境
```bash
# 安装 Go 1.24+
wget https://go.dev/dl/go1.24.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.24.6.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GO111MODULE=on' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

#### Docker 环境
```bash
# Ubuntu 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 启动 Docker
sudo systemctl start docker
sudo systemctl enable docker

# 验证安装
docker version
```

#### 其他工具
```bash
# Git
sudo apt update && sudo apt install -y git

# Make
sudo apt install -y build-essential

# 网络工具
sudo apt install -y curl wget jq

# 容器运行时 (可选)
sudo apt install -y containerd.io
```

## 🔧 源码获取和编译

### 1. 克隆源码
```bash
# 创建工作目录
mkdir -p $GOPATH/src/k8s.io
cd $GOPATH/src/k8s.io

# 克隆 Kubernetes 源码
git clone https://github.com/kubernetes/kubernetes.git
cd kubernetes

# 查看当前版本
git describe --tags --abbrev=0
```

### 2. 构建环境准备
```bash
# 安装依赖
make verify-go-version

# 验证构建环境
make verify

# 快速验证 (跳过耗时检查)
make quick-verify
```

### 3. 编译核心组件
```bash
# 编译所有组件 (耗时较长)
make all

# 编译特定组件
make WHAT=cmd/kubectl
make WHAT=cmd/kube-apiserver
make WHAT=cmd/kubelet

# 查看编译产物
ls -la _output/bin/
```

### 4. 快速编译 (开发模式)
```bash
# 只编译 Linux/amd64 版本
export KUBE_FASTBUILD=true
make all

# 或者使用并行编译
make -j$(nproc) all
```

## 🚀 本地集群搭建

### 方式一：使用 local-up-cluster.sh

#### 启动本地集群
```bash
# 进入 hack 目录
cd hack

# 启动本地集群 (单节点)
./local-up-cluster.sh

# 后台运行
nohup ./local-up-cluster.sh > cluster.log 2>&1 &
```

#### 配置 kubectl
```bash
# 设置 kubeconfig
export KUBECONFIG=/var/run/kubernetes/admin.kubeconfig

# 或者复制到默认位置
mkdir -p ~/.kube
cp /var/run/kubernetes/admin.kubeconfig ~/.kube/config

# 验证集群
kubectl get nodes
kubectl get pods -A
```

### 方式二：使用 kind (推荐)

#### 安装 kind
```bash
# 下载 kind
GO111MODULE="on" go install sigs.k8s.io/kind@v0.20.0

# 验证安装
kind version
```

#### 创建集群配置
```yaml
# kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 30000
    hostPort: 30000
  - containerPort: 30001
    hostPort: 30001
- role: worker
- role: worker
```

#### 启动集群
```bash
# 创建集群
kind create cluster --name k8s-dev --config kind-config.yaml

# 加载本地构建的镜像
kind load docker-image my-image:latest --name k8s-dev

# 删除集群
kind delete cluster --name k8s-dev
```

### 方式三：使用 kubeadm (生产环境)

#### 准备节点
```bash
# 禁用 swap
sudo swapoff -a
sudo sed -i '/ swap / s/^\(.*\)$/#\1/g' /etc/fstab

# 配置内核模块
cat <<EOF | sudo tee /etc/modules-load.d/k8s.conf
overlay
br_netfilter
EOF

sudo modprobe overlay
sudo modprobe br_netfilter

# 配置内核参数
cat <<EOF | sudo tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables  = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward                 = 1
EOF

sudo sysctl --system
```

#### 安装 kubeadm
```bash
# 安装 kubeadm, kubelet, kubectl
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl

curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-archive-keyring.gpg

echo "deb [signed-by=/etc/apt/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

sudo apt-get update
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl
```

#### 初始化集群
```bash
# 初始化 master 节点
sudo kubeadm init --pod-network-cidr=192.168.0.0/16

# 配置 kubectl
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

# 安装网络插件 (Calico)
kubectl apply -f https://docs.projectcalico.org/manifests/calico.yaml

# 允许在 master 节点调度 Pod (单节点集群)
kubectl taint nodes --all node-role.kubernetes.io/master-
```

## 🛠️ 开发工具配置

### VS Code 配置

#### 安装扩展
```bash
# 必要扩展
code --install-extension golang.go
code --install-extension ms-kubernetes-tools.vscode-kubernetes-tools
code --install-extension ms-vscode.vscode-json

# 推荐扩展
code --install-extension eamodio.gitlens
code --install-extension ms-vscode.vscode-yaml
code --install-extension redhat.vscode-yaml
```

#### 配置文件
```json
// .vscode/settings.json
{
    "go.useLanguageServer": true,
    "go.languageServerExperimentalFeatures": {
        "diagnostics": true,
        "documentLink": true
    },
    "go.lintTool": "golangci-lint",
    "go.lintFlags": [
        "--fast"
    ],
    "go.vetFlags": [
        "-composites=false"
    ],
    "go.buildTags": "",
    "go.testFlags": ["-v"],
    "go.testTimeout": "30s",
    "files.exclude": {
        "**/_output": true,
        "**/vendor": true,
        "**/.git": true
    }
}
```

#### 调试配置
```json
// .vscode/launch.json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug kube-apiserver",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/kube-apiserver",
            "args": [
                "--etcd-servers=http://127.0.0.1:2379",
                "--service-cluster-ip-range=10.0.0.0/24",
                "--insecure-bind-address=127.0.0.1",
                "--insecure-port=8080"
            ]
        },
        {
            "name": "Debug kubectl",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/kubectl",
            "args": ["get", "pods"]
        }
    ]
}
```

### GoLand 配置

#### 项目设置
```bash
# 设置 GOROOT 和 GOPATH
File -> Settings -> Go -> GOROOT: /usr/local/go
File -> Settings -> Go -> GOPATH: $HOME/go

# 配置代码格式
File -> Settings -> Editor -> Code Style -> Go
# 勾选 "Use gofmt" 和 "Use goimports"
```

#### 运行配置
```bash
# 创建运行配置
Run -> Edit Configurations -> + -> Go Build

# kube-apiserver 配置
Name: kube-apiserver
Package path: k8s.io/kubernetes/cmd/kube-apiserver
Arguments: --etcd-servers=http://127.0.0.1:2379 --service-cluster-ip-range=10.0.0.0/24

# kubectl 配置  
Name: kubectl
Package path: k8s.io/kubernetes/cmd/kubectl
Arguments: get pods
```

## 🔍 调试技巧

### 1. 使用 Delve 调试器

#### 安装 Delve
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

#### 调试示例
```bash
# 编译带调试信息的二进制
make WHAT=cmd/kubectl DBG=1

# 使用 dlv 调试
dlv exec ./_output/bin/kubectl -- get pods

# 设置断点
(dlv) break main.main
(dlv) break k8s.io/kubernetes/pkg/kubectl/cmd.(*Command).Execute

# 运行和检查
(dlv) continue
(dlv) print variableName
(dlv) stack
```

### 2. 日志调试

#### 启用详细日志
```bash
# kube-apiserver 详细日志
./_output/bin/kube-apiserver --v=4 --logtostderr

# kubectl 详细日志
./_output/bin/kubectl get pods --v=6

# 自定义日志级别
import "k8s.io/klog/v2"

klog.V(2).InfoS("Processing pod", "pod", pod.Name)
klog.ErrorS(err, "Failed to create pod", "pod", pod.Name)
```

### 3. 性能分析

#### pprof 性能分析
```go
// 添加到 main 函数
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()
```

#### 分析命令
```bash
# CPU 分析
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 内存分析
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看 goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine
```

### 4. 单元测试

#### 运行测试
```bash
# 运行所有测试
make test

# 运行特定包的测试
go test -v ./pkg/kubectl/cmd/

# 运行特定测试
go test -v ./pkg/kubectl/cmd/ -run TestApply

# 生成覆盖率报告
go test -v ./pkg/kubectl/cmd/ -coverprofile=coverage.out
go tool cover -html=coverage.out
```

#### 编写测试
```go
// example_test.go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := MyFunction(tt.input)
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

## 📊 监控和观察

### 1. 集群监控

#### 安装 Prometheus
```bash
# 使用 Helm 安装
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install prometheus prometheus-community/kube-prometheus-stack
```

#### 查看指标
```bash
# 转发 Prometheus 端口
kubectl port-forward svc/prometheus-kube-prometheus-prometheus 9090:9090

# 访问 http://localhost:9090
```

### 2. 分布式追踪

#### 安装 Jaeger
```bash
# 安装 Jaeger Operator
kubectl create namespace observability
kubectl create -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.47.0/jaeger-operator.yaml -n observability

# 创建 Jaeger 实例
kubectl apply -f - <<EOF
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: simple-prod
EOF
```

### 3. 日志聚合

#### 安装 ELK Stack
```bash
# 使用 Helm 安装 Elasticsearch
helm repo add elastic https://helm.elastic.co
helm install elasticsearch elastic/elasticsearch

# 安装 Kibana
helm install kibana elastic/kibana

# 安装 Filebeat
helm install filebeat elastic/filebeat
```

## 🚨 常见问题解决

### 构建问题

#### 内存不足
```bash
# 限制并行构建
export GOMAXPROCS=2
make all

# 使用交换空间
sudo fallocate -l 8G /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
```

#### 网络问题
```bash
# 配置 Go 代理
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn

# 或者关闭代理
export GOPROXY=direct
export GOSUMDB=off
```

### 运行时问题

#### 端口冲突
```bash
# 检查端口占用
sudo netstat -tlnp | grep :8080
sudo lsof -i :8080

# 杀死占用进程
sudo kill -9 <PID>
```

#### 权限问题
```bash
# 添加用户到 docker 组
sudo usermod -aG docker $USER
newgrp docker

# 重启 docker 服务
sudo systemctl restart docker
```

### 集群问题

#### 节点 NotReady
```bash
# 检查节点状态
kubectl describe node <node-name>

# 检查 kubelet 日志
sudo journalctl -u kubelet -f

# 重启 kubelet
sudo systemctl restart kubelet
```

#### Pod 启动失败
```bash
# 查看 Pod 事件
kubectl describe pod <pod-name>

# 查看容器日志
kubectl logs <pod-name> -c <container-name>

# 进入容器调试
kubectl exec -it <pod-name> -- /bin/sh
```

## 🎯 下一步

现在您已经搭建好了完整的开发环境，可以开始：

1. **浏览代码结构** - 使用 IDE 熟悉项目布局
2. **运行单元测试** - 验证环境配置正确
3. **启动本地集群** - 测试核心功能
4. **开始代码阅读** - 从感兴趣的组件开始

---

*环境搭建完成！接下来可以开始 [核心概念](./core-concepts.md) 的学习*
