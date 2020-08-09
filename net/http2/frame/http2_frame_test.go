package frame

import (
	"bytes"
	"log"
	"testing"

	"golang.org/x/net/http2"
)

/**
写Frame Header，表示打开一个流，此外还携带一个header块片段，HEADER Frame格式：
+---------------+
|Pad Length? (8)|
+-+-------------+-----------------------------------------------+
|E|                 Stream Dependency? (31)                     |
+-+-------------+-----------------------------------------------+
|  Weight? (8)  |
+-+-------------+-----------------------------------------------+
|                   Header Block Fragment (*)                 ...
+---------------------------------------------------------------+
|                           Padding (*)                       ...
+---------------------------------------------------------------+
*/
func writeFrameHeader(framer *http2.Framer) *http2.Framer {
	err := framer.WriteHeaders(http2.HeadersFrameParam{
		StreamID:      42,
		BlockFragment: []byte("test"),
		Priority:      http2.PriorityParam{},
		EndStream:     false,
		EndHeaders:    true,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
  写Frame Data，表示传送与一个流关联的任意的可变长度的字节序列，DATA Frame格式：
  +---------------+
  |Pad Length? (8)|
  +---------------+-----------------------------------------------+
  |                            Data (*)                         ...
  +---------------------------------------------------------------+
  |                           Padding (*)                       ...
  +---------------------------------------------------------------+

*/
func writeFrameData(framer *http2.Framer) *http2.Framer {
	if err := framer.WriteData(2, true, []byte("ok")); err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
写Frame Push Promise，表示在发送者打算初始化流之前通知客户端，PUSH_PROMISE Frame格式：
+---------------+
|Pad Length? (8)|
+-+-------------+-----------------------------------------------+
|R|                  Promised Stream ID (31)                    |
+-+-----------------------------+-------------------------------+
|                   Header Block Fragment (*)                 ...
+---------------------------------------------------------------+
|                           Padding (*)                       ...
+---------------------------------------------------------------+
*/
func writeFramePushPromise(framer *http2.Framer) *http2.Framer {
	pushPromise := http2.PushPromiseParam{
		StreamID:      1,
		PromiseID:     3,
		BlockFragment: []byte("test"),
	}
	if err := framer.WritePushPromise(pushPromise); err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
写Priority，表示流发送者建议的优先级，PRIORITY Frame格式：
+-+-------------------------------------------------------------+
|E|                  Stream Dependency (31)                     |
+-+-------------+-----------------------------------------------+
|   Weight (8)  |
+-+-------------+
*/
func writePriority(framer *http2.Framer) *http2.Framer {
	err := framer.WritePriority(12, http2.PriorityParam{
		StreamDep: 4,
		Exclusive: true,
		Weight:    34,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
写RST Stream，表示允许立即终止一个流，RST_STREAM Frame格式：
+---------------------------------------------------------------+
|                        Error Code (32)                        |
+---------------------------------------------------------------+
*/
func writeRSTStream(framer *http2.Framer) *http2.Framer {
	err := framer.WriteRSTStream(123, http2.ErrCode(http2.ErrCodeStreamClosed))
	if err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
写Settings，表示传达影响端点通信方式的配置参数，SETTINGS Frame格式：
+-------------------------------+
|       Identifier (16)         |
+-------------------------------+-------------------------------+
|                        Value (32)                             |
+---------------------------------------------------------------+
*/
func writeSettings(framer *http2.Framer) *http2.Framer {
	err := framer.WriteSettings([]http2.Setting{{1, 11}, {2, 22}}...)
	if err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
写Ping，表示测量来自发送方的最小往返时间，以及确定空闲连接是否仍然有效，PING Frame格式：
+---------------------------------------------------------------+
|                      Opaque Data (64)                         |
+---------------------------------------------------------------+
*/
func writePing(framer *http2.Framer) *http2.Framer {
	err := framer.WritePing(true, [8]byte{1, 2, 3, 4, 5, 6, 7, 8})
	if err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
写GoAway，表示启动连接关闭或发出严重错误状态信号，GOAWAY Frame格式：
+-+-------------------------------------------------------------+
|R|                  Last-Stream-ID (31)                        |
+-+-------------------------------------------------------------+
|                      Error Code (32)                          |
+---------------------------------------------------------------+
|                  Additional Debug Data (*)                    |
+---------------------------------------------------------------+
*/
func writeGoAway(framer *http2.Framer) *http2.Framer {
	if err := framer.WriteGoAway(123456, http2.ErrCode(0x123423), []byte("test")); err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
写Continuation，表示继续发送header块片段的序列，CONTINUATION Frame格式：
+---------------------------------------------------------------+
|                   Header Block Fragment (*)                 ...
+---------------------------------------------------------------+
*/
func writeContinuation(framer *http2.Framer) *http2.Framer {
	err := framer.WriteContinuation(42, true, []byte("end"))
	if err != nil {
		log.Fatalln(err)
	}
	return framer
}

/**
写WindowUpdate，用于流量控制，WINDOW_UPDATE Frame格式：
+-+-------------------------------------------------------------+
|R|              Window Size Increment (31)                     |
+-+-------------------------------------------------------------+
*/
func writeWindowUpdate(framer *http2.Framer) *http2.Framer {
	if err := framer.WriteWindowUpdate(23, 44); err != nil {
		log.Fatalln(err)
	}
	return framer
}

// 读Frame
func readFrame(framer *http2.Framer) {
	// 针对读取CONTINUATION帧需要设置为true，否则报错
	framer.AllowIllegalReads = true
	fr, _ := framer.ReadFrame()
	log.Printf("Frame type is [%s],Frame Header is %s \n", fr.Header().Type, fr.Header().String())

	//headerFrame, ok := fr.(*http2.HeadersFrame)
	//if ok {
	//	decoder := hpack.NewDecoder(2048, nil)
	//	hf, _ := decoder.DecodeFull(headerFrame.HeaderBlockFragment())
	//	for _, h := range hf {
	//		log.Printf("%s\n", h.Name+":"+h.Value)
	//	}
	//}
}

/**
Frame是HTTP/2协议中最小传输单位，每个Frame包含一个9字节大小的FrameHeader和不定长的Frame Payload：
Payload结构和内容取决于Type类型，如果Type类型是SETTINGS，Payload的大小将受SETTINGS_MAX_FRAME_SIZE参数限制。
Type类型包含：Headers、Data、PushPromise、Priority、Settings、Ping、RSTStream、GoAway、Continuation、WindowUpdate。
具体内容格式如下：
 +-----------------------------------------------+
 |                 Length (24)                   |
 +---------------+---------------+---------------+
 |   Type (8)    |   Flags (8)   |
 +-+-------------+---------------+-------------------------------+
 |R|                 Stream Identifier (31)                      |
 +=+=============================================================+
 |                   Frame Payload (0...)                      ...
 +---------------------------------------------------------------+

 HTTP/2 Frame：https://tools.ietf.org/html/rfc8336
 HTTP/2: https://tools.ietf.org/html/rfc7540#page-12
*/
func TestHttp2Frame(t *testing.T) {
	buf := new(bytes.Buffer)
	fr := http2.NewFramer(buf, buf)

	readFrame(writeFrameHeader(fr))
	readFrame(writeFrameData(fr))
	readFrame(writeFramePushPromise(fr))
	readFrame(writeRSTStream(fr))
	readFrame(writeSettings(fr))
	readFrame(writePriority(fr))
	readFrame(writePing(fr))
	readFrame(writeGoAway(fr))
	readFrame(writeContinuation(fr))
	readFrame(writeWindowUpdate(fr))
}
