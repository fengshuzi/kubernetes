package test

import "fmt"

// 导出的函数（首字母大写）
func SayHello(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// 非导出的函数（首字母小写）
func sayGoodbye() {
	fmt.Println("Goodbye!")
}

// 导出的变量
var Message = "This is from test package"
