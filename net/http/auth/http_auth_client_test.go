package auth

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

// 带有Basic认证的客户端请求
func basicAuth() {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8000/auth", nil)
	if err != nil {
		return
	}
	request.SetBasicAuth("abc", "123")
	response, err := client.Do(request)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	if response.StatusCode == http.StatusUnauthorized {
		log.Fatalf("认证失败. Response StatusCode：%v ，Response Content：%s ", response.StatusCode, string(b))
		return
	}
	defer response.Body.Close()
	log.Printf("认证成功。Response StatusCode：%v Response Content：%s ", response.StatusCode, string(b))
}

/**
进行TLS验证

使用OpenSSL命令生成相关密钥：
CA：
openssl genrsa -out client-ca.key 2048
openssl req -x509 -new -nodes -key client-ca.key -subj "/CN=ca.com" -days 5000 -out client-ca.crt
Client：
openssl genrsa -out client.key 2048
openssl req -new -key client.key -subj "/CN=client" -out client.csr
openssl x509 -req -in client.csr -CA client-ca.crt -CAkey client-ca.key -CAcreateserial -out client.crt -days 5000
*/
func tlsAuth() {
	pool := x509.NewCertPool()
	caCertFile := "ca.crt"

	caCrt, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Println("client read ca.crt error：", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	cliCrt, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Print("client loadx509keypair error：", err)
		return
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://localhost/")
	if err != nil {
		log.Println("ssl client error：", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	log.Println(string(body))
}

func TestHttpAuthClient(t *testing.T) {
	tlsAuth()
}
