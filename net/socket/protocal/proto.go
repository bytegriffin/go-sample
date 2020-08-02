package protocal

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

/*
TCP是以字节流的方式传输数据（而UDP是面向报文），传输的最小单位是一个报文段（Segment）。TCP Header中有个Option标识位，
常见的标识位为MSS（MaximuM Segment Size）最大分段大小：一般是由发送端向接收段确认每个分段数据大小，
MSS长度=MTU长度-IP Header-TCP Header，链路层每次传输数据时有个MTU（Maximum Transmission Unit）最大传输单元，
一般是1500字节，超过这个大小就要分成多个报文段，

发生粘包、拆包的原因：
	1.应用程序写入的数据字节大小 > 套接字发送缓冲区大小，这会导致拆包
	2.应用程序写入的数据字节大小 < 套接字发送缓冲区大小，网卡将应用程序多次写入的数据发送到网络上，这会导致粘包。
	3.进行MSS（最大报文长度）大小的TCP分段，当TCP报文长度 - TCP头部长度 > MSS 时会导致拆包。
	4.接收方不能及时读取套接字缓冲区数据，这会导致粘包。
处理TCP粘包问题，通常有几种方法：
	1.定长分隔：每个数据包最大为该长度，缺点是数据不足时会浪费带宽
	2.特定字符分隔：如rn，缺点是如果正文中包含rn字符就会导致问题
	3.每次发送完就断开连接，比如http1.0，缺点是每次都需要打开和关闭资源，开销大
	4.自定义消息格式：将消息分为消息头和消息体，并在消息头中添加消息长度：主要采用这种方式
*/

// 消息编码：将消息分成消息（包）头和消息（包）体，
// 消息头包含消息长度，消息体是消息内容。
func Encode(message string) ([]byte, error) {
	// 1.读取消息长度，转换成int32类型（占4个字节）
	length := int32(len(message))
	pkg := new(bytes.Buffer)
	// 2.生成消息头 小端模式
	err := binary.Write(pkg, binary.LittleEndian, length)
	if err != nil {
		return nil, err
	}
	// 3.生成消息体 小端模式
	err = binary.Write(pkg, binary.LittleEndian, []byte(message))
	if err != nil {
		return nil, err
	}
	// 4.返回消息
	return pkg.Bytes(), nil
}

// 消息解码：
func Decode(reader *bufio.Reader) (string, error) {
	// 1.读取消息长度，读取前4个字节
	len, _ := reader.Peek(4)
	buf := bytes.NewBuffer(len)
	var length int32
	err := binary.Read(buf, binary.LittleEndian, &length)
	if err != nil {
		return "", err
	}
	// 2.buffer返回缓冲中现有的可读取的字节数
	if int32(reader.Buffered()) < length+4 {
		return "", err
	}

	// 3.读取真正的消息体
	pack := make([]byte, int(4+length))
	_, err = reader.Read(pack)
	if err != nil {
		return "", err
	}
	return string(pack[4:]), nil
}
