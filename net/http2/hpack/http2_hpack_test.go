package hpack

import (
	"bytes"
	"golang.org/x/net/http2/hpack"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// 头信息是以键值对的形式存在
func encode() *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	encoder := hpack.NewEncoder(buf)
	encoder.WriteField(hpack.HeaderField{Name: ":status", Value: "200"})
	encoder.WriteField(hpack.HeaderField{Name: ":scheme", Value: "https"})
	encoder.WriteField(hpack.HeaderField{Name: "date", Value: time.Now().UTC().Format(http.TimeFormat)})
	encoder.WriteField(hpack.HeaderField{Name: "content-length", Value: strconv.Itoa(len("ok"))})
	encoder.WriteField(hpack.HeaderField{Name: "content-type", Value: "text/html"})
	return buf
}

func decode(buf *bytes.Buffer) {
	decoder := hpack.NewDecoder(2048, nil)
	headerFiled, _ := decoder.DecodeFull(buf.Bytes())
	for _, h := range headerFiled {
		log.Printf("header fields ==> %s\n", h.Name+" : "+h.Value)
	}

}

/**
为了弥补HTTP/1.1中存在重复且冗余的Header信息而造成的带宽和时间上的浪费，HTTP/2参考SPDY协议，
专门设计了一套头部压缩算法HPACK，HPACK使用静态索引表、动态索引表和哈夫曼编码把头部信息映射成为一个索引值，
当发送头信息时用索引值替代，从而达到压缩的目的。其中，静态索引表是由61个常规header域和一些预定义
的values组成的预定义字典表。动态索引表是在连接中遇到的实际header域的列表。动态索引表有大小限制，
新的key进来，旧的key可能会被移除。
当浏览器请求服务时，浏览器首先去查静态表，如果有完全与header信息匹配的键值对，就把它
对应的索引值取出来，例如":method: GET"对应的索引值是2，只需要传输2即可；如果只有相匹配的名称，
而value则不同，需要添加到动态表中，比如”cookie:xxxx“，并且浏览器要通知服务端同步动态表，后续
通讯可使用动态表所对应的字符即可。

HPACK：https://tools.ietf.org/html/rfc7541
*/
func TestHttp2Hpack(t *testing.T) {
	buf := encode()
	decode(buf)
}
