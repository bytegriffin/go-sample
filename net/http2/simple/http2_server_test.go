package simple

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

/**
  H2模式：需要TLS加密进行HTTP/2连接，在TLS握手期间会顺带完成HTTPS/2协议的协商，
  即：当客户端发送Client Hello时便指定ALPN Next Protocol为h2或http/1.1说明客户端支持的协议，
  如果双方协商失败（比如客户端或者服务端不支持），则会使用HTTPS/1.1继续通讯。

  使用OpenSSL为localhost生成私钥和证书，在win10下需要打开管理员模式的PowerShell：
  .\openssl req -x509 -out server.crt -keyout server.key -newkey rsa:2048 -nodes -sha256 -config localhost.cnf
  如果希望浏览器能访问正常，而非golang程序客户端，可以将证书安装到”受信任的根证书颁发机构“中，
  注意此时如果还想用golang程序客户端访问，就算不带TLS认证也能访问正确，因为证书之前已经安装到操作系统中了。
*/
func startH2Server() {

	//cert, err := tls.LoadX509KeyPair("public.crt", "private.key")
	//if err != nil {
	//	log.Fatal(err)
	//}

	server := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		//TLSConfig: &tls.Config{
		//	Certificates: []tls.Certificate{cert},
		//	ServerName:   "localhost",
		//},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Protocol: " + r.Proto)
			w.Write([]byte("Protocol: " + r.Proto))
		}),
	}

	//if err := http2.ConfigureServer(server, &http2.Server{}); err != nil {
	//	log.Fatal(err)
	//}

	log.Printf("start H2 server [localhost:8000]...\n")
	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatal(err)
	}
}

/**
  H2C模式：HTTP/2 ClearText，不需要TLS加密也可以进行HTTP/2连接，而是使用基于HTTP的握手
  来完成HTTP/2的升级，客户端使用HTTP Upgrade机制请求升级，如果服务端不支持HTTP/2，那么
  它会忽略Upgrade字段，直接返回HTTP/1.1的响应。如果服务器同意升级，那么会返回HTTP/1.1 101 Switching Protocols
*/
func startH2CServer() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %v, http: %v", r.URL.Path, r.TLS == nil)
	})
	server := &http.Server{
		Addr:    ":8000",
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}
	log.Printf("start H2C server [localhost:8000]...\n")
	server.ListenAndServe()
}

func TestHttp2Server(t *testing.T) {
	startH2Server()
}
