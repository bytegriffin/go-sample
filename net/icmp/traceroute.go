package icmp

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func traceroute(host string) {
	if runtime.GOOS == "windows" {
		log.Fatalln("该命令不支持Windows系统")
	}
	// 1. 查询DNS A记录
	ips, err := net.LookupIP(host)
	if err != nil {
		log.Fatal(err)
	}

	var dst net.IPAddr
	for _, ip := range ips {
		if ip.To4() != nil {
			dst.IP = ip
			log.Printf("使用 %v 跟踪 %s\n", dst.IP, host)
			break
		}
	}
	if dst.IP == nil {
		log.Fatal("没有发现IP")
	}

	// 2.开启监听ICMP应答
	pc, err := net.ListenPacket("ip4:1", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	p := ipv4.NewPacketConn(pc)
	if err := p.SetControlMessage(ipv4.FlagTTL|ipv4.FlagSrc|ipv4.FlagDst|ipv4.FlagInterface, true); err != nil {
		log.Fatal(err)
	}

	// 3.生成一个ICMP消息包
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Data: []byte(""),
		},
	}

	rb := make([]byte, 1500)
	for i := 1; i <= 64; i++ { //上限是64跳
		wm.Body.(*icmp.Echo).Seq = i
		wb, err := wm.Marshal(nil)
		if err != nil {
			log.Fatal(err)
		}
		if err := p.SetTTL(i); err != nil {
			log.Fatal(err)
		}

		// 每一跳通常有多条流量路径，每一跳都需要进行多次探测
		begin := time.Now()
		if _, err := p.WriteTo(wb, nil, &dst); err != nil {
			log.Fatal(err)
		}
		if err := p.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
			log.Fatal(err)
		}
		n, cm, peer, err := p.ReadFrom(rb)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				fmt.Printf("%v\t*\n", i)
				continue
			}
			log.Fatal(err)
		}

		// 6.解析应答
		// iana参考：https://godoc.org/golang.org/x/net/internal/iana
		const ProtocolICMP = 1
		rm, err := icmp.ParseMessage(ProtocolICMP, rb[:n])
		if err != nil {
			log.Fatal(err)
		}
		rtt := time.Since(begin)

		switch rm.Type {
		case ipv4.ICMPTypeTimeExceeded:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d\t%v %+v %v\n\t%+v\n", i, peer, names, rtt, cm)
		case ipv4.ICMPTypeEchoReply:
			names, _ := net.LookupAddr(peer.String())
			fmt.Printf("%d\t%v %+v %v\n\t%+v\n", i, peer, names, rtt, cm)
			return
		default:
			log.Printf("unknown ICMP message: %+v\n", rm)
		}
	}

}
