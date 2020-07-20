package socket

import (
	"fmt"
	"go-sample/util"
	"net"
	"testing"
)

func UDPServer() {
	// 1.获取udp地址
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	util.IsNilError("Server端udp地址出错：", err)

	// 2.获取服务端连接
	conn, err := net.ListenUDP("udp", addr)
	util.IsNilError("Server端监听地址出错：", err)
	defer conn.Close()

	for {
		// 3.接收数据，由于UDP是面向无连接的，所以没有Accept()这样的函数来获取Client连接的操作，直接接收Client端发来的数据即可
		var data [1024]byte
		n, addr, err := conn.ReadFromUDP(data[:])
		util.IsNilError("Server端udp读取数据出错：", err)
		fmt.Printf("服务端接收数据:%v addr:%v count:%v\n", string(data[:n]), addr, n)

		//4.回写数据
		_, err = conn.WriteToUDP(data[:n], addr)
		util.IsNilError("Server端udp发送数据出错：", err)

		fmt.Printf("服务端回写数据:%v addr:%v count:%v\n", string(data[:n]), addr, n)
	}

}

func TestUdpServer(t *testing.T) {
	UDPServer()
}
