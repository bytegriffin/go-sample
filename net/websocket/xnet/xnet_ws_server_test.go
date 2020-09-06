package xnet

import (
	"log"
	"net/http"
	"testing"

	"golang.org/x/net/websocket"
)

/**
  WebSocket是HTML5提供的一种在单个TCP连接上进行全双工通讯的协议。之前的网站实现推送技术大多采用
  HTTP Long Polling或Comet技术，但是效率低、浪费资源，而WebSocket具有性能开销小、
  服务端可主动向客户端发送数据、更好地二进制支持、通信高效等优点。
  浏览器通过HTTP/1.1协议Header中的Upgrade向服务端进行握手协商，申请协议升级，服务端同意后会返回
  给浏览器Header头写一个101 Switching Protocols的状态码确认，表示转换协议成功，之后双端就可主动发送消息。
  frame是websocket中最小通信单位，一个或多个frame又可封装成一个message接口供调用。具体格式如下：
   0                   1                   2                   3
   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
  +-+-+-+-+-------+-+-------------+-------------------------------+
  |F|R|R|R| opcode|M| Payload len |    Extended payload length    |
  |I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
  |N|V|V|V|       |S|             |   (if payload len==126/127)   |
  | |1|2|3|       |K|             |                               |
  +-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
  |     Extended payload length continued, if payload len == 127  |
  + - - - - - - - - - - - - - - - +-------------------------------+
  |                               |Masking-key, if MASK set to 1  |
  +-------------------------------+-------------------------------+
  | Masking-key (continued)       |          Payload Data         |
  +-------------------------------- - - - - - - - - - - - - - - - +
  :                     Payload Data continued ...                :
  + - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
  |        IANA             Payload Data continued ...            |
  +---------------------------------------------------------------+

  WebSocket：https://tools.ietf.org/html/rfc6455
  HTTP/2 WebSocket：https://tools.ietf.org/html/rfc8441
*/
func TestXNetWSServer(t *testing.T) {

	http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		//io.Copy(conn, conn)

		for {
			msg := make([]byte, 512)
			//读取消息，如果消息体不大的话，它将填充消息体，下一次读取将读取帧数据的其余部分。
			//n, err := conn.Read(msg)

			//采用Message格式接收
			err := websocket.Message.Receive(conn, &msg)

			//采用JSON格式接收
			//var user User
			//err := websocket.JSON.Receive(conn, &user)

			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Receive: %s\n", msg)
			// 发送方法一：
			sendMsg := "[" + string(msg) + "]"
			_, err = conn.Write([]byte(sendMsg))
			// 发送方法二：
			//err = websocket.Message.Send(ws, message)
			// 发送方法三：支持JSON
			//user := &User{Name: "aaa", Password: "123"}
			//websocket.JSON.Marshal(&user)
			//err = websocket.JSON.Send(ws, message)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Send: %s\n", sendMsg)
		}
	}))
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
