package file

import (
	"log"
	"net/http"
	"testing"
)

// 访问 http://localhost:9090/files/
func TestHttpFileServer(t *testing.T) {
	dir := http.FileServer(http.Dir("D:\\"))
	// 设置路由
	http.Handle("/files/", http.StripPrefix("/files/", dir))
	log.Fatal(http.ListenAndServe("localhost:9090", nil))
}
