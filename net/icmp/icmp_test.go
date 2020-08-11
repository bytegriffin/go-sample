package icmp

import (
	"testing"
)

/**
  网际控制报文协议：是一种面向无连接的网络层协议，用于在IP主机、路由器之间传递出错、查询报告控制消息。
  常用网络命令：ping、tracert和traceroute，ICMP基于一个“错误侦测与回报机制”，让人们
  能够检测网络的连线状况，常见差错报告：终点不可达、时间超过、ICMP重定向、参数问题、源点抑制。
  ICMP主要是通过不同的类别（Type）与代码（Code）让机器来识别不同的连线状况，具体结构如下：
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |     Type      |     Code      |          Checksum             |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |                             unused                            |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |      Internet Header + 64 bits of Original Data Datagram      |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

  ICMP：https://tools.ietf.org/html/rfc792
*/
func TestICMP(t *testing.T) {
	//ping("www.baidu.com")
	//ping("192.168.1.1")
	traceroute("www.baidu.com")
}
