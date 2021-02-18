package socks5

import (
	"golang.org/x/net/context"
	"golang.org/x/net/proxy"
	"log"
	"net"
	"net/http"
	"testing"
)

/**
  http socks5代理：客户端与目标服务器之间通讯的透明传递，
  设计最初是为了让有权限的用户可以穿过防火墙的限制，访问外部资源。
  socks4只支持TCP，socks5支持TCP和UDP。

  https://tools.ietf.org/html/rfc1928
  https://tools.ietf.org/html/rfc1929
*/
func TestHttpSocks5Client(t *testing.T) {

	httpTransport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// 带用户名/密码的socks5代理
			// auth := proxy.Auth{User: "user", Password: "password"}
			dial, _ := proxy.SOCKS5("tcp", "118.89.94.176:80", nil, proxy.Direct)
			return dial.Dial(network, addr)
		},
	}

	client := &http.Client{Transport: httpTransport}
	resp, err := client.Get("https://www.google.com")
	if err != nil {
		log.Println("err: ", err)
		return
	}
	log.Println(resp)
}
