/*
Http支持的几种认证方式：
1.Http Basic基本认证：采用Base64编码用户名和密码（明文）进行传输，容易被拦截并解码或被重放攻击，安全性较低。
2.Http Digest摘要认证：为了弥补Basic认证中发送明文的缺点，Http1.1开始采用Nonce随机数进行
  非可逆加密（例如MD5，SHA1），即：客户端将用户名、密码和nonce随机数生成的Salt值等字段一起加密，并将其计算
  出的值与服务端计算出的值进行比对是否相同，就算通讯期间被截取也不容易解码，但缺点是用户可以伪装身份访问资源。
3.SSL客户端认证：安全级别较高，一般使用时会跟搭配其他认证一起使用，申请认证者除了密码还需提供其他特有信息，
  可避免用户被第三者冒充的问题，但是需要花费一定费用去购买安全证书。
4.Http表单认证：一般采用基于Cookie、Session的方式，大多数人都在使用它创建动态网站，倘若遇到XSS攻击，可使用
  text/template包进行转码或者直接在Cookie中添加httponly属性。但是这种认证扩展性不佳，比如涉及到跨域就
  必须要求Session共享，每台服务器都要读取Session信息。
5.Bearer Token认证：它随着OAuth流行而流行，Token编码一般采用JWT格式，客户端和服务端之间通过Json对象来通讯，
  服务端不再保存任何Session数据，变成了Stateless无状态，从而提高了扩展性，实现了跨域认证。JWT比较适合分布式API调用、
  微服务之间调用、用户申请注册后的验证链接等，OAuth比较适合多账号登录授权，其中Access_code模式可以采用JWT的格式生成code，
  其缺点是由于服务端是无状态的，一旦被服务端颁布Token后，在到期之前不能被中止，一直有效。
*/

package auth

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

/*
Http Basic认证服务器

Http Basic认证过程：
1.客户端（浏览器）请求服务端
2.服务端响应客户端，发现客户端请求的资源需要验证才能访问，于是向浏览器返回状态码为401的Response，
  Response Header中会新增以下内容：WWW-Authenticate: Basic realm="realm是web资源的逻辑划分"
3.浏览器会弹出对话框，要求输入用户名和密码。输入完成点确定后，客户端会将用户名和密码以guest:guest的形式进行
  Base64加密并发送给服务端，此时客户端中的Request Header中会包含类似以下内容：Authorization Basic dGVzdedIder39
4.服务端接收到新的请求并对Authorization中的值进行解码并验证，如果验证成功则将相应的资源返回，否则返回401。
*/
func basicAuthHandler(resp http.ResponseWriter, req *http.Request) {
	username, password, _ := req.BasicAuth()
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		log.Println("认证失败，用户名或密码不能为空。")
		writeUnAuthorized(resp, "用户名或密码不能为空。")
		return
	}
	if username != "abc" || password != "123" {
		log.Println("认证失败，用户名或密码输入错误。Username: ", username, " ,password: ", password)
		writeUnAuthorized(resp, "用户名或密码输入错误。")
		return
	}
	log.Println("认证成功。")
	io.WriteString(resp, "hello, world!\n")
}

// 回写状态码和错误信息
func writeUnAuthorized(resp http.ResponseWriter, info string) {
	resp.WriteHeader(http.StatusUnauthorized)
	resp.Write([]byte(info))
}

// 开启 basic auth server
func startBasicAuthServer() {
	http.HandleFunc("/auth", basicAuthHandler)
	http.ListenAndServe("127.0.0.1:8000", nil)
}

type indexHandler struct {
}

func (i indexHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Print("client hello.")
	writer.Write([]byte("server hello"))
}

/*
开启 TLS auth server 进行双端认证，服务端只允许特定的客户端进行访问，

使用OpenSSL命令生成相关密钥：
CA：
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -subj "/CN=ca.com" -days 5000 -out ca.crt
server：
openssl genrsa -out server.key 2048
openssl req -new -key server.key -subj "/CN=server" -out server.csr
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 5000

TLS v1.3：https://tools.ietf.org/html/draft-ietf-tls-tls13-23
*/
func startTLSAuthServer() {
	pool := x509.NewCertPool()
	caCertFile := "client-ca.crt"
	caCrt, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		log.Println("server read crt error：", err)
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	s := &http.Server{
		Addr:    "localhost:443",
		Handler: &indexHandler{},
		TLSConfig: &tls.Config{
			ClientCAs:              pool,
			ClientAuth:             tls.RequireAndVerifyClientCert,
			MaxVersion:             tls.VersionTLS13,
			SessionTicketsDisabled: false,
		},
	}
	log.Fatalln(s.ListenAndServeTLS("server.crt", "server.key"))
}

func TestHttpAuthServer(t *testing.T) {
	startTLSAuthServer()
}
