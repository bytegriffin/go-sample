package simple

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"

	"golang.org/x/net/http2"
)

// h2 client
func startH2Client() {
	crt, err := ioutil.ReadFile("server.crt")
	if err != nil {
		log.Fatal(err)
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)

	client := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            rootCAs,
				InsecureSkipVerify: false,
				ServerName:         "localhost",
			},
			DisableCompression: true,
			AllowHTTP:          false,
		},
	}

	resp, err := client.Get("https://localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()

	certInfo := resp.TLS.PeerCertificates[0]
	log.Printf("过期时间: %v \n", certInfo.NotAfter)
	log.Printf("组织信息: %v \n", certInfo.Subject)

	log.Printf("Response Http Proto: %s Content: %s", resp.Proto, string(bytes))
}

// 没有相应的认证，直接访问会出错
func startH2ClientNoTls() {
	resp, err := http.Get("https://localhost:8000")
	if err != nil {
		log.Fatalln(err)
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Println(string(bytes))
}

// h2c client
func startH2CClient() {
	client := http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true, //充许非加密的http连接
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return net.Dial(network, addr)
			},
		},
	}
	resp, _ := client.Get("http://localhost:8000")
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	log.Printf("Response Http Proto: %s Content: %s", resp.Proto, string(bytes))
}

func TestHttp2Client(t *testing.T) {
	startH2Client()
}
