package reverse

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"testing"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Proxy server path: %v \n", r.Host+r.URL.Path)
	log.Printf(" Proxy server path: %v \n", r.Host+r.URL.Path)
	u, _ := url.Parse("http://127.0.0.1:9091/")
	proxy := httputil.NewSingleHostReverseProxy(u)

	// 重新设置Real Server的请求路径
	r.URL.Host = u.Host
	r.URL.Scheme = u.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = u.Host

	proxy.ServeHTTP(w, r)
}

/**
  Http反向代理：服务端代理，隐藏真正的服务器地址，常见的有LVS、nginx等。
  请求过程：Client ==》 Proxy Server ==》 Real Server
  使用浏览器或curl访问地址：http://127.0.0.1:8080/hello
*/
func TestHttpReverseProxyServer(t *testing.T) {
	http.HandleFunc("/hello", sayHello)
	log.Print("start http server 8080...")
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		fmt.Printf("Http Server [8080] failed, err：%v", err)
		return
	}

}
