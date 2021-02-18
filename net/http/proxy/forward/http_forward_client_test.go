package forward

import (
	"log"
	"net/http"
	"net/url"
	"testing"
)

/**
  Http正向代理：客户端代理，可隐藏客户端IP，比如翻墙。
*/
func TestHttpForwardClient(t *testing.T) {
	// 通过环境变量来指定代理
	// os.Setenv("HTTP_PROXY", "http://127.0.0.1:9743")

	//使用系统默认的代理
	// proxy := http.ProxyFromEnvironment

	// 使用代理
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse("http://103.146.176.124:80")
	}

	transport := &http.Transport{Proxy: proxy}

	client := &http.Client{Transport: transport}
	resp, err := client.Get("http://httpbin.org/get")

	if err != nil {
		log.Println(err)
		return
	}

	log.Println(resp)
}
