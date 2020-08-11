package icmp

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// IPv4 ping
func ipv4Ping(addr string) (*net.IPAddr, time.Duration, error) {
	// 1.开启监听ICMP应答
	pc, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return nil, 0, err
	}
	defer pc.Close()

	// 2.如果是域名解析出IP地址
	dst, err := net.ResolveIPAddr("ip4", addr)
	if err != nil {
		return nil, 0, err
	}

	// 3.生成一个ICMP消息包
	m := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(""),
		},
	}
	b, err := m.Marshal(nil)
	if err != nil {
		return dst, 0, err
	}

	// 4.发送消息
	start := time.Now()
	n, err := pc.WriteTo(b, dst)
	if err != nil {
		return dst, 0, err
	} else if n != len(b) {
		return dst, 0, fmt.Errorf("got %v; want %v", n, len(b))
	}

	// 5.等待应答
	reply := make([]byte, 1500)
	err = pc.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return dst, 0, err
	}
	n, peer, err := pc.ReadFrom(reply)
	if err != nil {
		return dst, 0, err
	}
	duration := time.Since(start)

	// 6.解析应答
	// iana参考：https://godoc.org/golang.org/x/net/internal/iana
	const ProtocolICMP = 1
	rm, err := icmp.ParseMessage(ProtocolICMP, reply[:n])
	if err != nil {
		return dst, 0, err
	}
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return dst, duration, nil
	default:
		return dst, 0, fmt.Errorf("来自 %v 的回复，%+v ", peer, rm)
	}

}

func ping(addr string) {
	dst, dur, err := ipv4Ping(addr)
	if err != nil {
		log.Printf("Ping %s (%s): %s\n", addr, dst, err)
		return
	}
	log.Printf("Ping %s (%s): %s\n", addr, dst, dur)
}
