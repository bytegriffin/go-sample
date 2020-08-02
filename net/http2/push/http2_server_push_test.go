package push

import (
	"log"
	"net/http"
	"testing"
)

/**
Http/2 server push
用浏览器或者nghttp2访问地址：https://localhost:9090/
使用OpenSSL命令生成私钥和证书：
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
*/
func TestHttp2ServerPush(t *testing.T) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 浏览器在正常Response之前会优先发送PUSH_PROMISE frame，来说明要为客户端推送资源，
		// 客户端收到请求后，针对这些资源就不会主动向服务端发请求。
		pusher, ok := w.(http.Pusher)
		if ok {
			if err := pusher.Push("/static/app.js", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
			if err := pusher.Push("/static/style.css", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
			if err := pusher.Push("/static/avatar.jpg", nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(`<html>
			<head>
				<title>Hello World</title>
				<script src="/static/app.js"></script>
				<link rel="stylesheet" href="/static/style.css"">
			</head>
			<body>
			Hello <span id="data"></span> !<br/>
			<img src="/static/avatar.jpg" />
			</body>
			</html>
			`))
	})

	log.Println("listening port 9090")
	log.Fatal(http.ListenAndServeTLS("localhost:9090", "server.crt", "server.key", nil))
}
