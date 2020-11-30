package grpcgateway

import (
	"log"
	"net/http"
	"path"
	"strings"
	"testing"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	"google.golang.org/grpc/credentials"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	gw "go-sample/net/grpc/http/grpc_gateway/proto"
)

const (
	// http地址
	HTTPAddress = ":8080"
)

/**
  启动H2C模式：不需要认证
*/
func startH2CGrpcGateServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	opts := []grpc.DialOption{grpc.WithInsecure()}
	gwmux := runtime.NewServeMux()
	err := gw.RegisterHelloHandlerFromEndpoint(ctx, gwmux, GrpcAddress, opts)
	if err != nil {
		log.Fatalf("failed to register gw server: %v", err)
	}

	// 2.启动http server
	httpmux := http.NewServeMux()
	httpmux.HandleFunc("/swagger", func(w http.ResponseWriter, req *http.Request) {
		if !strings.HasSuffix(req.URL.Path, "swagger.json") {
			log.Printf("Not Found: %s", req.URL.Path)
			http.NotFound(w, req)
			return
		}
		p := strings.TrimPrefix(req.URL.Path, "/swagger/")
		p = path.Join("proto", p)
		log.Printf("Serving swagger-file: %s", p)
		http.ServeFile(w, req, p)
	})
	//Aseet和AssetDir可能会在IDE下报错，但是不影响运行
	fileServer := http.FileServer(&assetfs.AssetFS{
		//Asset:    swagger.Asset,
		//AssetDir: swagger.AssetDir,
		Prefix: "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	httpmux.Handle(prefix, http.StripPrefix(prefix, fileServer))
	httpmux.Handle("/", gwmux)

	log.Println("开启Grpc-gateway服务...")

	if err := http.ListenAndServe(HTTPAddress, gwmux); err != nil {
		log.Fatal(err)
	}
}

/**
  启动H2模式：需要认证
*/
func startH2GrpcGateServer() {
	ctx := context.Background()

	// 1.启动gateway server
	//带证书认证，此时的gateway作为client会访问GatewayServer
	creds, err := credentials.NewClientTLSFromFile(CertFile, "localhost")
	if err != nil {
		log.Fatalf("failed to create client TLS credentials: %v", err)
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}

	gwmux := runtime.NewServeMux()
	//注册grpc server endpoint
	err = gw.RegisterHelloHandlerFromEndpoint(ctx, gwmux, GrpcAddress, opts)
	if err != nil {
		log.Fatalf("failed to register gw server: %v", err)
	}

	// 2.启动http server
	httpmux := http.NewServeMux()
	httpmux.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("pong"))
		//if !strings.HasSuffix(req.URL.Path, "swagger.json") {
		//	log.Printf("Not Found: %s", req.URL.Path)
		//	http.NotFound(w, req)
		//	return
		//}
		//p := strings.TrimPrefix(req.URL.Path, "/swagger/")
		//p = path.Join("message.swagger.json", p)
		//
		//log.Printf("Serving swagger-file: %s", p)
		//
		//http.ServeFile(w, req, p)
	})
	httpmux.Handle("/", gwmux)

	log.Println("开启Grpc-gateway服务...")
	if err := http.ListenAndServeTLS(HTTPAddress, CertFile, KeyFile, httpmux); err != nil {
		log.Fatal(err)
	}
}

/**
  测试接口：http://localhost:8080/get/test
          http://localhost:8080/post
  H2C模式测试命令：
          curl http://localhost:8080/get/test
          curl -X POST -k http://localhost:8080/post -d "{\"name\": \" world\"}"
*/
func TestGrpcGatewayServer(t *testing.T) {
	startH2CGrpcGateServer()
}
