package protocal

import (
	"bufio"
	"fmt"
	"go-sample/util"
	"io"
	"net"
	"testing"
)

func TestProtoTcpServer(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:12345")
	util.IsNilError("Server端启动监听失败。", err)
	defer listen.Close()

	for { // 循环接入所有请求的客户端
		conn, err := listen.Accept()
		util.IsNilError("Server端accept failed。", err)
		go func(conn net.Conn) {
			defer conn.Close()
			reader := bufio.NewReader(conn)
			for { // 循环读取完整消息
				msg, err := Decode(reader)
				if err == io.EOF {
					return
				}
				util.IsNilError("Server端解码出错。", err)
				fmt.Println("Server端接收到Client端发来的数据：", msg)
			}
		}(conn)
	}
}
