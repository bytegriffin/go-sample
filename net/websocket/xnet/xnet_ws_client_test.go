package xnet

import (
	"crypto/tls"
	"log"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

/**
  大多情况下，客户端是浏览器
*/
func TestXNetWSClient(t *testing.T) {
	// 连接方式一：可自定义配置
	config, err := websocket.NewConfig("ws://localhost:8080/ws", "http://localhost:8080/")
	if err != nil {
		log.Fatalln(err)
	}
	config.TlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	//config.Header.Add("Authorization", fmt.Sprintf("Basic %s",
	//	base64.StdEncoding.EncodeToString([]byte("convox:123456"))))
	config.Version = websocket.ProtocolVersionHybi13
	ws, err := websocket.DialConfig(config)

	// 连接方式二：直接连接
	//ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost:8080/")

	if err != nil {
		log.Fatalln(err)
	}
	//defer ws.Close()

	time.Sleep(time.Second * 2)
	message := []byte("hello, world!")
	// 发送方法一：
	//_, err = ws.Write(message)
	// 发送方法二：
	err = websocket.Message.Send(ws, message)
	// 发送方法三：支持JSON
	//user := &User{Name: "aaa", Password: "123"}
	//websocket.JSON.Marshal(&user)
	//err = websocket.JSON.Send(ws, message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Send: %s\n", message)

	msg := make([]byte, 512)
	// 接收方法一：
	//m, err := ws.Read(msg)
	// 接收方法二：
	err = websocket.Message.Receive(ws, &msg)
	// 接收方法三：支持JSON
	// err = websocket.JSON.Receive(ws, &user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Receive: %s\n", msg)

	ws.SetReadDeadline(time.Now().Add(time.Second * time.Duration(10)))

	//关闭连接，服务端也会关掉
	//ws.Close()
}
