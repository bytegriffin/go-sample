package main

import (
//	"syscall/js"
//	"time"
)

/**
  WebAssembly 的Dom交互操作
  需要引入syscall/js包，引入方法：如果使用Goland编译器，打开
  File-Settings-Go-Build Tags&Vendoring中OS选择js，Arch选择wasm。
  否则会报错，找不到syscall/js包，为了程序报错，先注释掉。

  https://godoc.org/syscall/js
*/
func main() {
	c := make(chan struct{}, 0)
	//js.Global().Get("console").Call("log", "Hello world Go/wasm!")
	//js.Global().Get("document").Call("getElementById", "app").Set("innerText", time.Now().String())

	// js显式调用Go代码 将Go中的writeHello方法注册到js的方法hello中
	//js.Global().Set("hello", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	//	return sayHello(args[0].String())
	//}))

	<-c
}

func sayHello(str string) string {
	return "hello " + str
}
