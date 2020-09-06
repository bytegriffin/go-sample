package echo

import (
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"nhooyr.io/websocket/wsjson"

	"nhooyr.io/websocket"
)

/**
  与gorilla/websocket相比，nhooyr.io/websocket具有以下优点：
  1.支持并发写，2.更友好的API，3.支持Wasm, 4.Close handshake
  5.掩码算法比gorilla/websocket快1.75倍  6.支持压缩扩展
  压缩扩展：https://tools.ietf.org/html/rfc7692
  Name Registry：https://tools.ietf.org/html/rfc7936
*/
func TestNhooyrWSServer(t *testing.T) {
	http.HandleFunc("/echo", func(writer http.ResponseWriter, request *http.Request) {
		conn, err := websocket.Accept(writer, request, &websocket.AcceptOptions{
			Subprotocols: []string{"echo"}, //子协议
			//InsecureSkipVerify: true,             //是否跳过同源验证
			//OriginPatterns:     []string{"example.com"}, //跨域

			//deflate压缩扩展，是否启用上下文接管。具体在HTTP Header中表现为以下内容：
			//Sec-WebSocket-Extensions: permessage-deflate;
			//  client_max_window_bits; server_max_window_bits=10;
			//  client_no_context_takeover; server_no_context_takeover
			//启用的好处是上一条消息的滑动窗口会对下一条消息的内容进行编码，这样可提高压缩比，缺点是上下文开销大，连接开始占用内存一直到连接关闭。
			//禁用的好处是客户端和服务端在不同消息之间可重置上下文，大大减少了连接开销，缺点是影响压缩性能。
			CompressionMode: websocket.CompressionNoContextTakeover,
			//压缩阈值,CompressionNoContextTakeover默认是512个字节，CompressionContextTakeover默认是128个字节
			CompressionThreshold: 128,
		})
		if err != nil {
			log.Fatalln(err)
		}
		defer conn.Close(websocket.StatusInternalError, "Server端内部出错了。")

		//判断客户端请求的子协议是否正确
		//if conn.Subprotocol() != "echo" {
		//	conn.Close(websocket.StatusPolicyViolation, "客户端必须是echo子协议才能进行协商。")
		//	return
		//}

		ctx, cancel := context.WithTimeout(request.Context(), time.Minute*10)
		defer cancel()

		for {
			var msg interface{}
			err = wsjson.Read(ctx, conn, &msg)
			if err != nil {
				log.Fatalln(err)
			}

			err = wsjson.Write(ctx, conn, &msg)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("server read %s", msg)
		}

	})
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
