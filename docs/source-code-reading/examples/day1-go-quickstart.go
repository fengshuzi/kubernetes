// Day 1 - Go语言快速入门 (为Java开发者定制)
// 这个文件可以直接运行: go run day1-go-quickstart.go

package main

import (
	"fmt"
	"sync"
	"time"
)

// 1. 接口定义 - 类似Java interface，但是隐式实现
type Worker interface {
	DoWork() string
	GetID() int
}

// 2. 结构体 - 类似Java的类
type JavaDeveloper struct {
	ID       int
	Name     string
	Language string
}

// 3. 方法实现 - 隐式实现Worker接口
func (j JavaDeveloper) DoWork() string {
	return fmt.Sprintf("%s 正在用 %s 开发", j.Name, j.Language)
}

func (j JavaDeveloper) GetID() int {
	return j.ID
}

// 4. Go开发者
type GoDeveloper struct {
	ID   int
	Name string
}

func (g GoDeveloper) DoWork() string {
	return fmt.Sprintf("%s 正在学习 Kubernetes 源码", g.Name)
}

func (g GoDeveloper) GetID() int {
	return g.ID
}

// 5. 工厂函数 - 类似Java的工厂模式
func CreateWorker(workerType string, id int, name string) Worker {
	switch workerType {
	case "java":
		return JavaDeveloper{ID: id, Name: name, Language: "Java"}
	case "go":
		return GoDeveloper{ID: id, Name: name}
	default:
		return nil
	}
}

// 6. 并发处理 - Go的核心特性
func demonstrateGoroutines() {
	fmt.Println("\n=== Goroutine 演示 (类似Java CompletableFuture) ===")
	
	var wg sync.WaitGroup
	results := make(chan string, 3)
	
	workers := []Worker{
		CreateWorker("java", 1, "张三"),
		CreateWorker("go", 2, "李四"),
		CreateWorker("java", 3, "王五"),
	}
	
	// 启动goroutine - 类似Java的异步执行
	for _, worker := range workers {
		wg.Add(1)
		go func(w Worker) {
			defer wg.Done()
			// 模拟工作耗时
			time.Sleep(time.Millisecond * 100)
			results <- w.DoWork()
		}(worker)
	}
	
	// 等待所有goroutine完成
	go func() {
		wg.Wait()
		close(results)
	}()
	
	// 收集结果
	for result := range results {
		fmt.Printf("✅ %s\n", result)
	}
}

// 7. 错误处理 - Go的显式错误处理方式
func processWorker(worker Worker) (string, error) {
	if worker == nil {
		return "", fmt.Errorf("worker 不能为空")
	}
	
	if worker.GetID() <= 0 {
		return "", fmt.Errorf("worker ID 必须大于0")
	}
	
	return worker.DoWork(), nil
}

// 8. 指针和内存管理
func demonstratePointers() {
	fmt.Println("\n=== 指针演示 (Java开发者需要理解的概念) ===")
	
	// 值传递 vs 引用传递
	dev := JavaDeveloper{ID: 1, Name: "原始名字", Language: "Java"}
	fmt.Printf("修改前: %s\n", dev.Name)
	
	// 通过指针修改
	changeNameByPointer(&dev, "修改后的名字")
	fmt.Printf("修改后: %s\n", dev.Name)
}

func changeNameByPointer(dev *JavaDeveloper, newName string) {
	dev.Name = newName
}

// 9. 类型断言 - 类似Java的instanceof和类型转换
func demonstrateTypeAssertion() {
	fmt.Println("\n=== 类型断言演示 ===")
	
	var worker Worker = CreateWorker("java", 1, "Java开发者")
	
	// 类型断言 - 类似Java的 (JavaDeveloper) worker
	if javaDev, ok := worker.(JavaDeveloper); ok {
		fmt.Printf("这是一个Java开发者: %s, 使用语言: %s\n", javaDev.Name, javaDev.Language)
	} else {
		fmt.Println("这不是Java开发者")
	}
}

// 10. 主函数 - 程序入口
func main() {
	fmt.Println("🚀 欢迎Java老兵学习Go语言!")
	fmt.Println("这个程序演示了Go语言的核心特性，帮助你快速上手")
	
	// 1. 基本用法演示
	fmt.Println("\n=== 基本接口使用 ===")
	workers := []Worker{
		CreateWorker("java", 1, "资深Java开发"),
		CreateWorker("go", 2, "Go新手"),
	}
	
	for _, worker := range workers {
		result, err := processWorker(worker)
		if err != nil {
			fmt.Printf("❌ 错误: %v\n", err)
		} else {
			fmt.Printf("✅ %s\n", result)
		}
	}
	
	// 2. 并发演示
	demonstrateGoroutines()
	
	// 3. 指针演示
	demonstratePointers()
	
	// 4. 类型断言演示
	demonstrateTypeAssertion()
	
	// 5. Go vs Java 对比总结
	fmt.Println("\n=== Go vs Java 关键差异 ===")
	fmt.Println("1. 接口: Go是隐式实现，Java是显式implements")
	fmt.Println("2. 并发: Go用goroutine+channel，Java用Thread+并发包")
	fmt.Println("3. 错误: Go返回error，Java抛出Exception")
	fmt.Println("4. 内存: Go有指针但有GC，Java纯GC")
	fmt.Println("5. 继承: Go用组合，Java用继承")
	
	fmt.Println("\n🎯 下一步: 运行 'go run day2-client-go-basic.go' 开始K8s之旅!")
}

// 额外的学习提示
/*
对于Java开发者的Go学习重点:

1. 思维转换:
   - 从面向对象转向接口+组合
   - 从异常处理转向错误返回
   - 从线程池转向goroutine

2. 语法重点:
   - defer 语句 (类似Java的finally)
   - range 遍历 (类似Java的for-each)
   - make 和 new 的区别
   - 空接口 interface{} (类似Java的Object)

3. 最佳实践:
   - 错误处理要显式检查
   - 使用go fmt格式化代码
   - 利用go vet检查问题
   - 善用go mod管理依赖

4. K8s相关:
   - client-go库的使用方式
   - context包的重要性
   - yaml标签的使用
   - controller模式的实现
*/
