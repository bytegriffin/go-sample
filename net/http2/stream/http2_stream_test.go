package stream

import (
	"bytes"
	"log"
	"testing"

	"golang.org/x/net/http2"
)

/**
  Stream是HTTP/2协议中一个较为完整的逻辑处理单元，即：一次完整的资源请求-响应数据交换流程，
  它表示客户端和服务器之间交换的独立双向帧序列。流具有以下特点：
  1.一个HTTP/2连接Conn可同时保持多个打开的流，两端可交叉发送不同流的Frame。
  2.客户端或服务端都可创建新的Stream
  3.客户端或服务端都可先关闭Stream
  4.同一个Stream内的Frame保证有序
  5.Stream使用整数值标记，流ID为奇数表示客户端发送HEADER/DATA消息，流ID为偶数表示服务端发送PUSH消息
  流的状态分为：idle、reserved、open、half-closed、closed。
  Stream的状态机：
                             +--------+
                     send PP |        | recv PP
                    ,--------|  idle  |--------.
                   /         |        |         \
                  v          +--------+          v
           +----------+          |           +----------+
           |          |          | send H /  |          |
    ,------| reserved |          | recv H    | reserved |------.
    |      | (local)  |          |           | (remote) |      |
    |      +----------+          v           +----------+      |
    |          |             +--------+             |          |
    |          |     recv ES |        | send ES     |          |
    |   send H |     ,-------|  open  |-------.     | recv H   |
    |          |    /        |        |        \    |          |
    |          v   v         +--------+         v   v          |
    |      +----------+          |           +----------+      |
    |      |   half   |          |           |   half   |      |
    |      |  closed  |          | send R /  |  closed  |      |
    |      | (remote) |          | recv R    | (local)  |      |
    |      +----------+          |           +----------+      |
    |           |                |                 |           |
    |           | send ES /      |       recv ES / |           |
    |           | send R /       v        send R / |           |
    |           | recv R     +--------+   recv R   |           |
    | send R /  `----------->|        |<-----------'  send R / |
    | recv R                 | closed |               recv R   |
    `----------------------->|        |<----------------------'
                             +--------+

       send:   endpoint sends this frame
       recv:   endpoint receives this frame

       H:  HEADERS frame (with implied CONTINUATIONs)
       PP: PUSH_PROMISE frame (with implied CONTINUATIONs)
       ES: END_STREAM flag
       R:  RST_STREAM frame

  Stream：https://tools.ietf.org/html/rfc7540#page-16
*/
func TestHttp2Stream(t *testing.T) {
	buf := new(bytes.Buffer)
	fr := http2.NewFramer(buf, buf)

	// 以HeadersFrame为例
	hp := http2.HeadersFrameParam{
		//每个Frame都有一个属性Stream Identifier，来表示属于哪个流，保证收发有序。
		//straemID本身会递增，不会被重用。当客户端发送StreamID为2Header帧时，并且未使用过1号帧发送，那么1号帧会编程close状态。
		StreamID: 2,
		//是否是流的终点，该值为true发送端的流为half-closed(local)状态，接收端的流为half-closed(remote)
		//在这个状态上发送RST_STREAM frame可以使状态立刻变成closed，而closed的流上不允许发送除PRIORITY之外的Frame。
		EndStream:     true, //当前帧是流的最后一个帧
		EndHeaders:    true, //当前帧是头信息的最后一个帧
		PadLength:     0,    //填充标志，在数据Payload里填充无用信息，用于干扰信道监听
		BlockFragment: []byte("test"),
		//请求优先级设置：每个stream都可以设置依赖和权重，可以按依赖树分配优先级，解决了关键请求被阻塞的问题
		Priority: http2.PriorityParam{
			StreamDep: 2,     //依赖，0表示不依赖
			Exclusive: false, //是否独占
			Weight:    255,   //权重，与StreamDep参数一起设置或都不设置，取值空间1-256
		},
	}
	fr.WriteHeaders(hp)

	f, err := fr.ReadFrame()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(f.Header().String())
}
