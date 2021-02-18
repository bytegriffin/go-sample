package main

import (
	"log"
	"net/http"
)

/**
由于hello.go文件引入syscall/js包，设置了环境变量，
因此运行此server需要再次设置下
set GOARCH=amd64
set GOOS=windows
*/
func main() {
	log.Println("----start webassembly http server---")
	log.Fatal(http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`E:\gopath\src\go-sample\net\webassembly\dom\`))))
}
