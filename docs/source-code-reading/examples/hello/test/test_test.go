package test

import "testing"

func TestSayHello(t *testing.T) {
	// 这里可以测试 SayHello 函数
	// 因为在同一个包内，所以可以访问
	SayHello("测试")

	// 也可以调用非导出函数
	sayGoodbye()
}
