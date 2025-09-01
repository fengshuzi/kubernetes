// Day 2 - client-go 基础实战 (Java开发者版)
// 运行前需要: go mod init && go get k8s.io/client-go@latest k8s.io/api@latest k8s.io/apimachinery@latest
// 运行: go run day2-client-go-basic.go

package main

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// K8sManager - 封装K8s操作，类似Java的Service类
type K8sManager struct {
	clientset kubernetes.Interface
	namespace string
}

// NewK8sManager - 构造函数，类似Java的构造器
func NewK8sManager(namespace string) (*K8sManager, error) {
	// 构建kubeconfig路径
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	
	// 创建配置
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("构建配置失败: %v", err)
	}
	
	// 创建客户端
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建客户端失败: %v", err)
	}
	
	return &K8sManager{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

// 1. CRUD操作 - 类似Java的DAO操作

// CreateJavaApp - 创建Java应用 (类似数据库插入操作)
func (k *K8sManager) CreateJavaApp(name, image string, port int32) error {
	fmt.Printf("🚀 创建Java应用: %s\n", name)
	
	// 创建Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: k.namespace,
			Labels: map[string]string{
				"app":  name,
				"type": "java-app",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  name,
						"type": "java-app",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: port,
									Name:          "http",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "JAVA_OPTS",
									Value: "-Xmx512m -Xms256m",
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    "500m",
									corev1.ResourceMemory: "512Mi",
								},
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    "200m",
									corev1.ResourceMemory: "256Mi",
								},
							},
						},
					},
				},
			},
		},
	}
	
	// 创建Deployment
	_, err := k.clientset.AppsV1().Deployments(k.namespace).Create(
		context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("创建Deployment失败: %v", err)
	}
	
	// 创建Service
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name + "-service",
			Namespace: k.namespace,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": name,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(int(port)),
				},
			},
			Type: corev1.ServiceTypeClusterIP,
		},
	}
	
	_, err = k.clientset.CoreV1().Services(k.namespace).Create(
		context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("创建Service失败: %v", err)
	}
	
	fmt.Printf("✅ Java应用 %s 创建成功\n", name)
	return nil
}

// GetJavaApps - 获取Java应用列表 (类似数据库查询)
func (k *K8sManager) GetJavaApps() ([]string, error) {
	fmt.Println("📋 获取Java应用列表...")
	
	deployments, err := k.clientset.AppsV1().Deployments(k.namespace).List(
		context.TODO(), metav1.ListOptions{
			LabelSelector: "type=java-app",
		})
	if err != nil {
		return nil, fmt.Errorf("获取Deployment列表失败: %v", err)
	}
	
	var apps []string
	for _, dep := range deployments.Items {
		status := "Unknown"
		if dep.Status.ReadyReplicas == *dep.Spec.Replicas {
			status = "Running"
		} else if dep.Status.ReadyReplicas == 0 {
			status = "Pending"
		} else {
			status = "Partial"
		}
		
		apps = append(apps, fmt.Sprintf("%s (状态: %s, 副本: %d/%d)", 
			dep.Name, status, dep.Status.ReadyReplicas, *dep.Spec.Replicas))
	}
	
	return apps, nil
}

// UpdateJavaApp - 更新Java应用 (类似数据库更新)
func (k *K8sManager) UpdateJavaApp(name string, replicas int32) error {
	fmt.Printf("🔄 更新Java应用 %s 副本数为 %d\n", name, replicas)
	
	// 获取现有Deployment
	deployment, err := k.clientset.AppsV1().Deployments(k.namespace).Get(
		context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("获取Deployment失败: %v", err)
	}
	
	// 更新副本数
	deployment.Spec.Replicas = &replicas
	
	// 应用更新
	_, err = k.clientset.AppsV1().Deployments(k.namespace).Update(
		context.TODO(), deployment, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新Deployment失败: %v", err)
	}
	
	fmt.Printf("✅ Java应用 %s 更新成功\n", name)
	return nil
}

// DeleteJavaApp - 删除Java应用 (类似数据库删除)
func (k *K8sManager) DeleteJavaApp(name string) error {
	fmt.Printf("🗑️ 删除Java应用: %s\n", name)
	
	// 删除Service
	err := k.clientset.CoreV1().Services(k.namespace).Delete(
		context.TODO(), name+"-service", metav1.DeleteOptions{})
	if err != nil {
		fmt.Printf("⚠️ 删除Service失败: %v\n", err)
	}
	
	// 删除Deployment
	err = k.clientset.AppsV1().Deployments(k.namespace).Delete(
		context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("删除Deployment失败: %v", err)
	}
	
	fmt.Printf("✅ Java应用 %s 删除成功\n", name)
	return nil
}

// 2. Watch机制 - 类似Java的观察者模式
func (k *K8sManager) WatchPods(duration time.Duration) error {
	fmt.Printf("👀 开始监听Pod变化 (%v)...\n", duration)
	
	// 创建ListWatch - 类似Java的EventListener
	listWatcher := cache.NewListWatchFromClient(
		k.clientset.CoreV1().RESTClient(),
		"pods",
		k.namespace,
		fields.Everything(),
	)
	
	// 创建Informer
	_, controller := cache.NewInformer(
		listWatcher,
		&corev1.Pod{},
		time.Second*30,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pod := obj.(*corev1.Pod)
				if hasLabel(pod, "type", "java-app") {
					fmt.Printf("➕ Java Pod创建: %s (状态: %s)\n", 
						pod.Name, pod.Status.Phase)
				}
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				oldPod := oldObj.(*corev1.Pod)
				newPod := newObj.(*corev1.Pod)
				if hasLabel(newPod, "type", "java-app") && 
				   oldPod.Status.Phase != newPod.Status.Phase {
					fmt.Printf("🔄 Java Pod状态变化: %s (%s -> %s)\n", 
						newPod.Name, oldPod.Status.Phase, newPod.Status.Phase)
				}
			},
			DeleteFunc: func(obj interface{}) {
				pod := obj.(*corev1.Pod)
				if hasLabel(pod, "type", "java-app") {
					fmt.Printf("➖ Java Pod删除: %s\n", pod.Name)
				}
			},
		},
	)
	
	// 启动controller
	stopCh := make(chan struct{})
	go controller.Run(stopCh)
	
	// 运行指定时间后停止
	time.Sleep(duration)
	close(stopCh)
	
	fmt.Println("✅ 监听结束")
	return nil
}

// 工具函数
func int32Ptr(i int32) *int32 { return &i }

func hasLabel(pod *corev1.Pod, key, value string) bool {
	if pod.Labels == nil {
		return false
	}
	return pod.Labels[key] == value
}

// 演示函数 - 类似Java的单元测试
func demonstrateBasicOperations(manager *K8sManager) {
	fmt.Println("\n=== 基础CRUD操作演示 ===")
	
	appName := "demo-java-app"
	
	// 1. 创建应用
	if err := manager.CreateJavaApp(appName, "openjdk:11-jre-slim", 8080); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		return
	}
	
	// 2. 等待一下让资源创建
	time.Sleep(2 * time.Second)
	
	// 3. 查询应用
	apps, err := manager.GetJavaApps()
	if err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
		return
	}
	
	fmt.Println("📋 当前Java应用:")
	for _, app := range apps {
		fmt.Printf("   - %s\n", app)
	}
	
	// 4. 更新应用
	if err := manager.UpdateJavaApp(appName, 2); err != nil {
		fmt.Printf("❌ 更新失败: %v\n", err)
		return
	}
	
	// 5. 等待更新生效
	time.Sleep(3 * time.Second)
	
	// 6. 再次查询
	apps, err = manager.GetJavaApps()
	if err != nil {
		fmt.Printf("❌ 查询失败: %v\n", err)
		return
	}
	
	fmt.Println("📋 更新后的Java应用:")
	for _, app := range apps {
		fmt.Printf("   - %s\n", app)
	}
	
	// 7. 清理资源
	if err := manager.DeleteJavaApp(appName); err != nil {
		fmt.Printf("❌ 删除失败: %v\n", err)
		return
	}
}

func main() {
	fmt.Println("🎯 Day 2: client-go 基础实战")
	fmt.Println("这个程序演示了如何使用client-go操作Kubernetes资源")
	
	// 创建K8s管理器
	manager, err := NewK8sManager("default")
	if err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		fmt.Println("请确保:")
		fmt.Println("1. kubectl 配置正确 (kubectl get nodes)")
		fmt.Println("2. 有访问K8s集群的权限")
		return
	}
	
	fmt.Println("✅ 连接K8s集群成功")
	
	// 演示基础操作
	demonstrateBasicOperations(manager)
	
	// 演示Watch机制
	fmt.Println("\n=== Watch机制演示 ===")
	fmt.Println("💡 现在去另一个终端创建/删除一些Pod，观察这里的输出")
	fmt.Println("例如: kubectl run test-pod --image=nginx --labels=type=java-app")
	
	if err := manager.WatchPods(10 * time.Second); err != nil {
		fmt.Printf("❌ Watch失败: %v\n", err)
	}
	
	// 总结
	fmt.Println("\n=== 学习总结 ===")
	fmt.Println("✅ 掌握了client-go的基本使用")
	fmt.Println("✅ 理解了K8s资源的CRUD操作")
	fmt.Println("✅ 学会了Watch机制监听资源变化")
	fmt.Println("✅ 对比了K8s客户端操作与Java数据库操作的相似性")
	
	fmt.Println("\n🎯 下一步: 学习Controller模式和Informer机制")
	fmt.Println("运行: go run day3-controller-pattern.go")
}

/*
对Java开发者的学习提示:

1. client-go vs JDBC/JPA:
   client-go.Clientset      ≈ DataSource/EntityManager
   Create/Get/Update/Delete ≈ CRUD操作
   Watch                    ≈ EventListener/Observer
   ListOptions              ≈ Query条件

2. 资源操作模式:
   Kubernetes资源操作类似数据库操作:
   - Deployment/Service等 ≈ 数据库表
   - Labels/Selectors ≈ WHERE条件
   - Watch ≈ 数据库触发器/事件

3. 异步和并发:
   - Goroutine处理Watch事件 ≈ Java异步处理
   - Channel通信 ≈ BlockingQueue
   - Context取消 ≈ Future.cancel()

4. 最佳实践:
   - 总是检查错误返回值
   - 使用正确的context
   - 合理设置超时时间
   - 注意资源清理
*/
