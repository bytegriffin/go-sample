package client

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"testing"
	"time"

	"golang.org/x/net/context"
)

func retry(attempts int, timeout time.Duration, client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("client倒数第%d次连接出错，出错信息：%s ", attempts, err)
		if attempts--; attempts > 0 {
			time.Sleep(timeout)
			return retry(attempts, timeout, client, req)
		}
		return nil, err
	}
	return resp, err
}

func httpProxy() func(_ *http.Request) (*url.URL, error) {
	u := url.URL{}
	urlProxy, err := u.Parse("http://110.243.11.176:9999")
	if err != nil {
		log.Println("http proxy error ", err)
		return nil
	}
	return http.ProxyURL(urlProxy)
}

func client() {
	//Transport表示一个Http事务，是RoundTripper的一个实现，本身已经实现了连接复用
	tr := &http.Transport{
		DisableCompression:     false,            //是否取消gzip压缩
		DisableKeepAlives:      false,            // 是否取消长连接
		MaxIdleConns:           10,               //所有host的idle状态的最大连接数目
		MaxIdleConnsPerHost:    1,                //每个host的idle状态的最大连接数目
		MaxConnsPerHost:        10,               //每个host上的最大连接数目，包含dialing/active/idle状态的连接数。http2中每个host只允许有一个idle的connection
		IdleConnTimeout:        30 * time.Second, //连接保持idle状态的最大时间，0表示不受限制
		ResponseHeaderTimeout:  30 * time.Second, //发送完request后等待serve response的时间
		ExpectContinueTimeout:  1 * time.Second,  //从客户端发送包含Expect:100-continues请求头到响应后继续发送post data的间隔时间
		MaxResponseHeaderBytes: 0,                //Server端响应头Header的最大字节数，0表示默认
		WriteBufferSize:        0,                //向transport中写缓冲大小，默认0表示4k
		ReadBufferSize:         0,                //从transport中读缓冲大小，默认0表示4k
		//如果配置了DialContext或TLSClientConfig参数时，默认会关闭http2，
		//如果想要使用自定义Dialer或TLS，又想打开HTTP2的话，那么可以将此参数设置为true
		ForceAttemptHTTP2: false,
		//Proxy: http.ProxyFromEnvironment,//系统环境设置 os.Setenv("HTTP_PROXY", "http://127.0.0.1:9743")
		//Proxy: httpProxy(), // http代理
		// 创建未加密的Dial连接
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// socks5代理
			// dial, _ := proxy.SOCKS5("tcp","136.244.78.50:32082",nil,proxy.Direct)
			dial := net.Dialer{
				Timeout:   30 * time.Second,                 //dial等待连接建立的最大时长，如果已经设置了Deadline参数，它将更早地失败
				Deadline:  time.Now().Add(30 * time.Second), //dial等待连接建立的截至时间点，如果已经设置了Timeout参数，它将更早地失败
				KeepAlive: 30 * time.Millisecond,            //网络保活的心跳间隔，默认15s
				//LocalAddr: &net.TCPAddr{//本地网络地址
				//	IP: net.ParseIP("192.168.1.103"),
				//	Port: 0,
				//	Zone: "",//ipv6 zone
				//},
				FallbackDelay: 30 * time.Millisecond, //ipv6配置错误回退到ipv4之前需要等待ipv6成功的时间间隔，默认是300ms
				// Resolver: ,//可选项，可重新设置
				// 控制器：在创建网络连接后dial之前会被执行
				Control: func(network, address string, c syscall.RawConn) error {
					log.Println("======在dial之前先打印这里======")
					return nil
				},
			}
			return dial.Dial(network, addr)
		},

		// 非代理模式的https连接，如果该值为nil，那么需使用DialContext和TLSClientConfig开启https，如果不为nil，那么程序
		// 将会忽略TLSClientConfig和TLSHandshakeTimeout参数，DialContext将不能被用于HTTPS请求。
		//DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		//	return nil, nil
		//},
		// 客户端的TLS配置，默认为nil，如果不为nil，那么http2将不会被开启
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, //是否开启中间人攻击来进行测试，如果为true则表示客户端将不再对服务端的证书进行校验
		},
		TLSHandshakeTimeout: 0, //0表示TLS握手无时间限制
	}

	req, err := http.NewRequest("GET", "http://httpbin.org/get", nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second, //0表示无限制
	}

	resp, err := retry(2, 0*time.Second, client, req)
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("Response Content: ", string(bytes))

}

func TestHttpClient(t *testing.T) {
	client()
}
