package main

import "fmt"

/**
  WebAssembly是”基于堆栈的Web虚拟机“，它可以是任何编程语言构建Web应用程序，
  只需要编译成wasm格式的二进制文件，就可实现跨平台运行，wasm要比javascript运行速度更快。

  操作步骤：
  1.windows下执行build.cmd脚本，生成hello.wasm文件，以便html页面导入该wasm文件
  2.启动WebAssembly Http Server
  3.打开浏览器，访问http://localhost:8080，按F12查看Console中的信息:hello world

  https://github.com/golang/go/wiki/WebAssembly
*/
func main() {
	fmt.Println("hello world.")
}
