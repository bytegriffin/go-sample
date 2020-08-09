package flow

import (
	"testing"
)

/**
  由于HTTP/1的流量控制完全是在TCP层上，粒度太粗，而HTTP/2协议引入的流量控制可作用在HTTP应用层上，主要是控制发送端允许发送一定数量的DATA帧。
  流控的作用是保证同一个TCP连接上的多个流Stream之间不会相互阻塞，从而更好地支持多路复用。新建连接时，两端都有一个初始值为
  65535字节大小的流控窗口flow-control window，发送端每发送一个DATA帧，就会把window值递减，而接收端每接收到一个DATA帧，
  就会回传一个WINDOW_UPDATE帧给发送端，告诉发送端只能允许发送多大的DATA帧，并将WINDOW_UPDATE帧内的Window Size Increment值加到Window值上。
  其中流控窗口Flow-control window的大小实际上代表了接收端能够接收到数据的速度快慢。

  Flow Control：https://tools.ietf.org/html/rfc7540#page-22
*/
type flow struct {
	_ incomparable
	// n表示接收端允许发送端发送的DATA字节数，它的大小可衡量接收端的缓存能力
	n int32

	// 多个Stream流共享一个TCP连接
	conn *flow
}

type incomparable [0]func()

func (f *flow) setConnFlow(cf *flow) {
	f.conn = cf
}

func (f *flow) available() int32 {
	n := f.n
	if f.conn != nil && f.conn.n < n {
		n = f.conn.n
	}
	return n
}

func (f *flow) take(n int32) {
	if n > f.available() {
		panic("internal error: took too much")
	}
	f.n -= n
	if f.conn != nil {
		f.conn.n -= n
	}
}

func (f *flow) add(n int32) bool {
	sum := f.n + n
	if (sum > n) == (f.n > 0) {
		f.n = sum
		return true
	}
	return false
}

func TestFlow(t *testing.T) {
	var st flow
	var conn flow
	st.add(3)
	conn.add(2)

	n := st.available()
	println(n)

	st.setConnFlow(&conn)
	n = st.available()
	println(n)

	b := st.add(1)
	println(b)
	n = st.available()
	println(n)

	st.take(1)
	n = st.available()
	println(n)

}
