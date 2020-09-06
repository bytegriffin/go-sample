package simple

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	WriteBufferSize:  10240,
	ReadBufferSize:   10240,
	HandshakeTimeout: 5 * time.Millisecond, //升级websocket协议所需要的握手时间

	//源头检查，如果客户端是浏览器，该字段为必选项，如果客户端不是浏览器，该字段为可选。
	//检查的目的是为了防止跨站点请求伪造 CSRF（cross-site request forgery），如果不检查或允许Origin的话，
	//则返回true；如果该值为nil，则使用安全默认值；如果Origin Host不等于请求头Host，则返回false
	//Origin：https://tools.ietf.org/html/rfc6454
	CheckOrigin: func(r *http.Request) bool {
		if r.Method != "GET" {
			log.Println("必须是Get方法")
			return false
		}
		if r.URL.Path != "/ws" {
			log.Println("请求路径出错")
			return false
		}
		return true //允许跨域
	},
	//WriteBufferPool: nil, //写缓冲池

	//按顺序指定服务支持的子协议，客户端可以指定其 Sec-WebSocket-Protocol 为其所期望
	//采用的子协议集合，而服务端则可以在此集合中选取一个并返回给客户端。
	//Subprotocols: []string{"chat","superchat"},

	//http的错误响应，该值如果为nil，那么会默认使用http.Error来当错误响应，如果服务端不想接收连接的话，
	//它必须返回适当的HTTP错误状态码（比如 403 Forbidden），并且终止WebSocket的握手过程
	//Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
	//},

	//是否支持压缩，设置为true表示客户端会向服务端申请协商是否支持压缩，如果服务端也同意，
	//则会使用deflate压缩技术来压缩消息，一般在请求头中用以下字段表示：
	//Sec-WebSocket-Extensions: permessage-deflate; client_max_window_bits
	//Compression Extensions： https://tools.ietf.org/html/rfc7692
	EnableCompression: true,
}

/**
  WebSocket服务端
  与gorilla/websocket相比，golang.org/x/net/websocket包中的websocket很长时间不更新
  而且功能不全，例如golang.org/x/net/websocket不支持用户定义连接之间的I/O缓冲区等，具体区别：
  +---------------------------------+------------------------------------------+
  |                                 |  github.com/gorilla   | golang.org/x/net |
  +----------------------------------------------------------------------------+
  | RFC 6455 Features                                                          |
  +----------------------------------------------------------------------------+
  | Passes Autobahn Test Suite     |        Yes            |        No         |
  +----------------------------------------------------------------------------+
  | Receive fragmented message     |        Yes            |        No         |
  +----------------------------------------------------------------------------+
  | Send close message             |        Yes            |        No         |
  +----------------------------------------------------------------------------+
  | Send pings and receive pongs   |        Yes            |        No         |
  +----------------------------------------------------------------------------+
  | Get the type of a received data message  |  Yes        |        Yes        |
  +----------------------------------------------------------------------------+
  | Other Features                                                             |
  +----------------------------------------------------------------------------+
  | Compression Extensions         |      Experimental     |        No         |
  +----------------------------------------------------------------------------+
  | Read message using io.Reader   |        Yes            |        No         |
  +----------------------------------------------------------------------------+
  | Write message using io.WriteCloser  |        Yes       |        No         |
  +----------------------------------------------------------------------------+
*/
func TestGorillaWSServer(t *testing.T) {

	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		//升级协议，如果升级失败返回HTTP错误响应
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Fatalln(err)
		}
		//defer conn.Close()

		for {
			conn.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(5000)))
			mt, message, err := conn.ReadMessage()
			// 读取JSON类型
			//var user User
			//err = conn.ReadJSON(user)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					log.Printf("ReadMessage from remote:%v error: %v \n", conn.RemoteAddr(), err)
				} else {
					log.Fatalln(err)
				}
			}
			//log.Printf("server receive: %s \n", message)

			sendMsg := "[" + string(message) + "]"
			err = conn.WriteMessage(mt, []byte(sendMsg))
			//写JSON类型
			//user = User{Name: "aaa", Password: "123"}
			//conn.WriteJSON(&user)
			if err != nil {
				log.Fatalln(err)
			}
			log.Printf("server send: %s \n", sendMsg)
		}
	})
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
