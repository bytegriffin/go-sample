package protocal

import (
	"go-sample/util"
	"net"
	"strconv"
	"testing"
)

func TestProtoClient(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	util.IsNilError("Client Dial failed。", err)
	defer conn.Close()

	for i := 0; i < 10; i++ {
		data, err := Encode("hello world " + strconv.Itoa(i))
		util.IsNilError("Client encode failed。", err)
		conn.Write(data)
	}

}
