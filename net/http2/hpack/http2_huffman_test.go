package hpack

import (
	"encoding/hex"
	"golang.org/x/net/http2/hpack"
	"log"
	"testing"
)

func huffmanEncoding() {
	items := [][2]string{
		{"a8eb10649cbf", "no-cache"},
		{"f1e3c2e5f23a6ba0ab90f4ff", "www.example.com"},
		{"25a849e95ba97d7f", "custom-key"},
		{"25a849e95bb8e8b4bf", "custom-value"},
		{"6402", "302"},
	}

	for _, item := range items {
		encodedHex := []byte(item[0])
		encoded := make([]byte, len(encodedHex)/2)
		_, err := hex.Decode(encoded, encodedHex)
		if err != nil {
			log.Println(err)
		}
		data := hpack.HuffmanEncodeLength(item[1])
		if err != nil {
			log.Println(err)
		}
		log.Println(item[1], data)
	}
}

func huffmanDecoding() {
	items := [][2]string{
		{"a8eb10649cbf", "no-cache"},
		{"f1e3c2e5f23a6ba0ab90f4ff", "www.example.com"},
		{"25a849e95ba97d7f", "custom-key"},
		{"25a849e95bb8e8b4bf", "custom-value"},
	}

	for _, item := range items {
		encodedHex := []byte(item[0])
		encoded := make([]byte, len(encodedHex)/2)
		_, err := hex.Decode(encoded, encodedHex)
		if err != nil {
			log.Println(err)
		}
		decoded, err := hpack.HuffmanDecodeToString(encoded)
		if err != nil {
			log.Println(err)
		}
		log.Println(item[1], decoded)
	}

}

/**
  如果HPACK在静态或动态字典表中查找不到相匹配的内容，需要使用哈夫曼算法对相应的内容进行编码
  以便减小体积。HTTP/2协议中使用了一份静态哈夫曼码表，需要同时内置在客户端和服务端。

  huffman code：https://httpwg.org/specs/rfc7541.html#huffman.code
*/
func TestHttp2Huffman(t *testing.T) {
	huffmanEncoding()
	huffmanDecoding()
}
