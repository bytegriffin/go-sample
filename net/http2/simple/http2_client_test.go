package simple

import (
	"crypto/tls"
	"crypto/x509"
	"golang.org/x/net/http2"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
)

// h2 client
func startH2Client() {
	crt, err := ioutil.ReadFile("public.crt")
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

	log.Printf("Response Http Proto: %s Content: %s", resp.Proto, string(bytes))
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
