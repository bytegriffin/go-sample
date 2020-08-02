package socket

import (
	"fmt"
	"go-sample/util"
	"net"
	"testing"
)

func TcpClient() {
	// 1.连接Server端
	conn, err := net.Dial("tcp", "127.0.0.1:1234")
	util.IsNilError("Client端连接Server端时失败。", err)
	defer conn.Close()

	// 2.给Server端发生消息
	_, err = conn.Write([]byte("hello world"))
	util.IsNilError("Client端发送消息失败。", err)

	// 3.Client端接受消息
	buf := [512]byte{}
	n, err := conn.Read(buf[:])
	util.IsNilError("Client端接收失败。", err)

	fmt.Println("Client端接收服务端回写消息：" + string(buf[:n]))
}

func TestTcpClient(t *testing.T) {
	TcpClient()
}
