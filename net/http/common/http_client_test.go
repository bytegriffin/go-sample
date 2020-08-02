package http

import (
	"bytes"
	"encoding/json"
	"go-sample/util"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"
)

func get() {
	// 方法一：利用http来进行Get请求
	//resp, _ := http.Get("http://127.0.0.1:8000/form")
	//defer resp.Body.Close()
	//log.Println("Response Header：", resp.Header)

	// 方法二：利用http.Client和http.NewRequest
	client := &http.Client{}
	// client.Get("") // 也可以使用client来Get请求
	request, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8000/form", nil)
	request.Header.Set("Content-Type", "application/json") //设置header
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	log.Println("Response Header：", resp.Header)

	buf := make([]byte, 1024)
	n, err := resp.Body.Read(buf)
	res := util.IsHttpNilError("http request error.", err)

	if res {
		log.Println("Http Content：", string(buf[:n]))
	}
}

// Post表单数据
func postForm() {
	data := url.Values{"id": {"123"}, "name": {"abc"}}
	resp, _ := http.Post("http://127.0.0.1:8000/form",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()))

	resp, err := http.PostForm("http://example.com/form",
		url.Values{"key": {"Value"}, "id": {"123"}})

	defer resp.Body.Close()
	log.Println("Response Header：", resp.Header)

	buf := make([]byte, 1024)
	n, err := resp.Body.Read(buf)
	res := util.IsHttpNilError("http request error.", err)
	if res {
		log.Println("Http Content：", string(buf[:n]))
	}
}

// Post Json格式的数据内容
func postJson() {
	// 设置request，序列化
	var user User
	user.Id = 123
	user.Name = "abc"
	jsonStr, _ := json.Marshal(user)
	req, err := http.NewRequest("Post", "http://127.0.0.1:8000/json", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	util.IsNilError("PostJson() request is error.", err)
	defer req.Body.Close()
	// 通过httpclient请求
	client := &http.Client{Timeout: 5 * time.Second}
	res, err := client.Do(req)
	util.IsNilError("PostJson() response is error.", err)
	defer res.Body.Close()
	// 读取io内容
	content, _ := ioutil.ReadAll(res.Body)
	log.Println("postJson() response content：", string(content))
}

// 上传文件
// 测试服务器地址：http://httpbin.org/post
func uploadFile() {
	// 模拟form表单中的一个form文件
	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)
	writer.WriteField("field", "this is a field")
	w, _ := writer.CreateFormFile("fieldName", "test.txt")
	w.Write([]byte("this is file content."))
	writer.Close()
	//模拟浏览器post请求
	resp, _ := http.Post("http://httpbin.org/post", writer.FormDataContentType(), buff)
	defer resp.Body.Close()
	// 接收服务器端口回写的数据
	data, _ := ioutil.ReadAll(resp.Body)
	log.Println("uploadFile() response content：", string(data))
}

// 设置cookie
// 测试服务器地址：http://httpbin.org/cookies/set?username=abc&password=123
func setCookie() {
	// 第一次请求：服务端会自动生成cookie，并返回给客户端
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}
	u, _ := url.Parse("http://127.0.0.1:8000/cookie")
	res, _ := client.Get(u.String())
	defer res.Body.Close()
	io.Copy(os.Stdout, res.Body)
	log.Println()

	// 打印第一次请求的cookie值
	for _, cookie := range jar.Cookies(u) {
		log.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}

	//第二次请求：客户端会读取服务端返回的cookie值
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, _ := client.Do(req)
	// 打印第二次请求的cookie值 方法一：使用cookiejar
	for _, cookie := range jar.Cookies(u) {
		log.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
	// 方法二：使用response同样可以打印出cookie值
	for _, cookie := range resp.Cookies() {
		log.Printf("  %s: %s\n", cookie.Name, cookie.Value)
	}
}

// net/http包自动会处理redirect跳转，如果要禁止用的话开启CheckRedirect
func redirect() {
	// 下列方法可以禁止redirect
	//client := &http.Client{
	//	CheckRedirect: func(req *http.Request, via []*http.Request) error {
	//		return http.ErrUseLastResponse
	//	},
	//}

	//要把http://www.baidu.com中的http:进行转码，否则报错
	queryEncode := url.QueryEscape("http://www.baidu.com")
	log.Println(queryEncode)
	res, err := http.Get("http://127.0.0.1:8000/redirect?url=" + queryEncode)
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()
	log.Println(res.Status)
	// 跳转到baidu首页
	data, _ := ioutil.ReadAll(res.Body)
	log.Println("uploadFile() response content：", string(data))
}

func TestHttpClient(t *testing.T) {
	redirect()
}
