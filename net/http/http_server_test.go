package http

import (
	"encoding/json"
	"fmt"
	"go-sample/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type indexHandler struct {
	content string
}

func (i indexHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(i.content))
}

// 日志中间件：可循环嵌套handler
func loggingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("开始访问 %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("结束访问 %s in %v", r.URL.Path, time.Since(start))
	})
}

// 处理Json格式
func jsonHandler(res http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	var user User
	// 反序列化：将Json字符串 转换成 结构体
	err := json.Unmarshal(body, &user)
	util.IsNilError("jsonHandler is failed.", err)
	user.Id = 456
	user.Name = "def"
	// 序列化：将struct 转换成 json字符串
	ret, _ := json.Marshal(user)
	fmt.Fprint(res, string(ret))
}

// 上传文件
func uploadFileHandler(res http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("content-type")
	contentLength := req.ContentLength
	log.Printf("upload content-type:%s,content-length:%d \n", contentType, contentLength)
	if !strings.Contains(contentType, "multipart/form-data") {
		res.Write([]byte("Content-Type must be multipart/form-data"))
		return
	}
	if contentLength > 10*1024*1024 {
		res.Write([]byte("file to large,limit 10MB"))
		return
	}
	err := req.ParseMultipartForm(10 * 1024 * 1024)
	if err != nil {
		res.Write([]byte("ParseMultipartForm error；" + err.Error()))
		return
	}
	if len(req.MultipartForm.File) == 0 {
		res.Write([]byte("have not any file."))
		return
	}

	for name, files := range req.MultipartForm.File {
		if name == "" {
			res.Write([]byte("File data is null."))
			return
		}
		for _, file := range files {
			f, _ := file.Open()
			path := "./" + file.Filename
			dst, _ := os.Create(path)
			io.Copy(dst, f)
			log.Printf("upload successful. fileName=%s fileSize=%v byte. savePath=%s",
				file.Filename, file.Size, path)
			// 往客户端回写数据
			data, _ := ioutil.ReadAll(req.Body)
			fmt.Fprint(res, data)
		}
	}
}

// Cookie管理
func cookieHandler(res http.ResponseWriter, req *http.Request) {

	cookie, err := req.Cookie("cookieName")
	if err != nil {
		// 设置Cookie
		cookie = &http.Cookie{
			Name:    "cookieName",
			Value:   "cookieValue",
			Path:    "/",
			Expires: time.Now().AddDate(1, 0, 0),
		}
		log.Println("第一次访问，服务端设置cookie值。")
	} else {
		cookie = &http.Cookie{
			Name:    "cookieName",
			Value:   "updateCookieValue",
			Path:    "/",
			Expires: time.Now().AddDate(1, 0, 0),
		}
		log.Println("第二次访问，服务端修改cookie值。")
	}

	// 删除Cookie，不用设置value值，将MaxAge设置成-1即可
	//cookie := http.Cookie{
	//	Name: "cookieName",
	//	Path: "/",
	//	MaxAge: -1,
	//}
	http.SetCookie(res, cookie)

	cookieStr := cookie.String()
	log.Println("Server端已经设置好cookie值：", cookieStr)
	res.Write([]byte(cookieStr))

	// 读取Cookie
	//ck, err := req.Cookie("cookieName")
	//if err != nil {
	//	return
	//}
	//log.Println("cookieName的值为：", ck.Value)
}

func TestHttpServer(t *testing.T) {
	// 注册路由的实现方法一：注入一个匿名回调函数
	http.HandleFunc("/form", func(res http.ResponseWriter, req *http.Request) {
		// 设置Response Header
		res.Header().Set("Content-Type", "text/html")
		// 回写网页内容方式一
		res.Write([]byte("测试页面."))
		// 回写网页内容方式二
		//fmt.Fprintln(res, "hello world")
	})

	// 注册路由的实现方法二：实现ServeHTTP接口
	http.Handle("/", &indexHandler{content: "首页内容"})
	http.Handle("/log", loggingHandler(&indexHandler{content: "日志管理页面"}))
	http.HandleFunc("/json", jsonHandler)
	http.HandleFunc("/upload", uploadFileHandler)
	http.HandleFunc("/cookie", cookieHandler)
	http.ListenAndServe("127.0.0.1:8000", nil)

	// 注册路由的实现方法三：方法一和二是使用DefaultServeMux，以下使用多路复用ServeMux来自定义实现
	//mux := http.NewServeMux()
	//mux.Handle("/", &indexHandler{content: "首页内容"})
	//mux.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
	//	fmt.Fprintln(res, "hello world")
	//})
	//http.ListenAndServe("127.0.0.1:8000", mux)
}
