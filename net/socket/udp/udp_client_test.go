package socket

import (
	"fmt"
	"go-sample/util"
	"net"
	"testing"
)

func UdpClient() {
	// 1. 获取udp地址
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:1234")
	util.IsNilError("Client端获取UDP地址失败。", err)

	// 2. 连接服务端
	conn, err := net.DialUDP("udp", nil, udpAddr)
	util.IsNilError("Client端连接UDP地址失败。", err)
	defer conn.Close()

	// 3.发送数据
	sendData := []byte("hello world")
	_, err = conn.Write(sendData)
	util.IsNilError("Client端发送UDP数据失败。", err)

	// 4.接收数据
	recData := make([]byte, 4096)
	n, rAddr, err := conn.ReadFromUDP(recData)
	util.IsNilError("Client端接收Server端的UDP数据失败。", err)
	fmt.Printf("客户端接收:%v addr:%v count:%v\n", string(recData[:n]), rAddr, n)
}

func TestUdpClient(t *testing.T) {
	UdpClient()
}
