package socket

import (
	"fmt"
	"go-sample/util"
	"net"
	"testing"
)

func TcpServer() {
	// 1.监听端口
	listen, err := net.Listen("tcp", "127.0.0.1:1234")
	util.IsNilError("Server端监听失败。", err)
	defer listen.Close()

	// 2.等待连接
	conn, err := listen.Accept()
	util.IsNilError("Server端Accept失败。", err)
	defer conn.Close()

	// 3. 接受客户端信息
	var buf [128]byte
	n, err := conn.Read(buf[:])
	util.IsNilError("Server读取客户端失败。", err)
	rec := string(buf[:n])
	fmt.Println("Server接收Client端发来的数据：", rec)

	// 4.返回给客户端信息
	_, err = conn.Write([]byte("Server端回写数据：" + rec))
	util.IsNilError("Server端回写数据时出错。", err)
}

func TestTcpServer(t *testing.T) {
	TcpServer()
}
