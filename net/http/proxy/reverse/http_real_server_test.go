package reverse

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestHttpRealServer(t *testing.T) {
	// real server的请求路径与proxy server的请求路径有关，
	// 可以写成 / 或者 /hello，否则会输出404 page not found
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Real server path: %v \n", r.Host+r.URL.Path)
		log.Printf("Real server path: %v \n", r.Host+r.URL.Path)
	})
	log.Print("start http server 9091...")
	if err := http.ListenAndServe("127.0.0.1:9091", nil); err != nil {
		log.Printf("Http Server [9091] failed, err：%v", err)
		return
	}
}
