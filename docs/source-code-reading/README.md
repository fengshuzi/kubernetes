# Kubernetes 源码阅读指南

## 📚 文档导航

欢迎来到 Kubernetes 源码分析文档！这里提供了系统性的源码阅读指南和详细的技术分析。

## 🗂️ 文档结构

### 📖 源码阅读指南
- [project-structure.md](./project-structure.md) - 项目结构分析
- [components/README.md](./components/README.md) - 组件分析总览

### 🔧 HTTP 框架分析
- [http-framework/README.md](./http-framework/README.md) - HTTP 框架分析总览
- [http-framework/kube-apiserver-http-framework.md](./http-framework/kube-apiserver-http-framework.md) - API Server HTTP 框架详解

### 🏗️ 组件分析
- [components/kube-apiserver.md](./components/kube-apiserver.md) - API Server 详细分析
- [components/kube-controller-manager.md](./components/kube-controller-manager.md) - 控制器管理器分析
- [components/kube-scheduler.md](./components/kube-scheduler.md) - 调度器源码分析
- [components/kubelet.md](./components/kubelet.md) - Kubelet 节点代理分析
- [components/kube-proxy.md](./components/kube-proxy.md) - 网络代理组件分析

## 🎯 学习路径

### 初学者路径
1. **项目结构** - 了解整体代码组织
2. **HTTP 框架** - 理解 API Server 的请求处理机制
3. **API Server** - 掌握核心 API 服务器
4. **kubectl** - 学习客户端交互

### 进阶路径
1. **调度器** - 深入调度算法
2. **控制器** - 掌握控制器模式
3. **网络代理** - 理解网络实现
4. **节点管理** - 学习 Kubelet 机制

### 专家路径
1. **性能优化** - 分析性能瓶颈
2. **安全机制** - 深入认证授权
3. **存储系统** - 理解数据持久化
4. **扩展开发** - 学习插件机制

## 🔍 快速导航

### 按技术栈分类
- **HTTP 框架**: Go-Restful + 标准库
- **存储系统**: etcd + 自定义存储
- **网络**: CNI + kube-proxy
- **调度**: 多级调度算法
- **安全**: RBAC + 准入控制

### 按功能模块分类
- **API 层**: RESTful API 设计
- **控制层**: 控制器模式实现
- **数据层**: 资源存储和管理
- **网络层**: 服务发现和负载均衡
- **安全层**: 认证授权机制

## 📊 技术栈概览

| 组件 | 主要技术 | 关键特性 |
|------|----------|----------|
| **API Server** | Go-Restful + HTTP/2 | RESTful API、自动路由 |
| **Scheduler** | 多级调度算法 | 资源感知、亲和性 |
| **Controller** | 控制器模式 | 声明式管理、事件驱动 |
| **Kubelet** | CRI 接口 | 容器生命周期管理 |
| **Proxy** | iptables/ipvs | 服务发现、负载均衡 |

## 🛠️ 开发工具

### 代码分析工具
```bash
# 查找函数定义
grep -r "func FunctionName" pkg/

# 查找接口实现
grep -r "type.*Interface" pkg/

# 查看调用关系
go list -f '{{.ImportPath}} {{.Imports}}' ./...
```

### 调试技巧
```bash
# 启用详细日志
export KUBE_LOG_LEVEL=5

# 查看 API 请求
kubectl get pods -v=8

# 分析性能
go tool pprof http://localhost:8080/debug/pprof/profile
```

## 📝 贡献指南

欢迎为文档贡献内容！

### 文档规范
- 使用中文撰写
- 保持结构清晰
- 添加代码示例
- 包含架构图表

### 内容要求
- **准确性**: 确保代码分析正确
- **完整性**: 覆盖主要功能模块
- **实用性**: 提供实际的学习价值
- **时效性**: 基于最新版本分析

---

*开始您的 Kubernetes 源码学习之旅！通过系统性的阅读和分析，您将深入理解这个强大的容器编排平台。*