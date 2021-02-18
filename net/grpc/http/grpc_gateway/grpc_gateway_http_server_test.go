package grpcgateway

import (
	pb "go-sample/net/grpc/http/grpc_gateway/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"golang.org/x/net/context"

	gw "go-sample/net/grpc/http/grpc_gateway/proto"
)

const (
	// grpc地址
	GrpcAddress = ":9192"
	CertFile    = "server.pem"
	KeyFile     = "server.key"
)

type grpcServer struct {
	pb.UnimplementedHelloServer
}

func (s *grpcServer) GetHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received get method：%v", req.GetName())
	return &pb.HelloResponse{Code: 200, Message: "hello get"}, nil
}

func (s *grpcServer) PostHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received post method：%v", req.GetName())
	return &pb.HelloResponse{Code: 200, Message: "hello post"}, nil
}

/**
  内部开启RPC，外部开启Restful Api，只需要部署一套代码
  Grpc-Gateway内部会生成一个反向代理服务器，会将外部的Restful api转换成rpc。
  一、生成密钥：
  1.生成RSA私钥：openssl genrsa -out server.key 2048
  2.生成ECC私钥：openssl ecparam -genkey -name secp384r1 -out server.key
  3.自定义公钥：openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650

  二、安装Grpc-gateway：
  1.将grpc-gateway/third_party/googleapis/google/api下的proto文件复制到本地proto目录下
  2.生成gRPC stubs：message.pb.go文件
  protoc -I . --go_out . --go_opt paths=source_relative --go-grpc_out . --go-grpc_opt paths=source_relative message.proto
  3.使用protoc-gen-grpc-gateway生成反向代理message.pb.gw.go文件
  protoc -I . --grpc-gateway_out ./  --grpc-gateway_opt logtostderr=true  --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true message.proto
  4.使用protoc-gen-openapiv2生成OpenAPI，即：message.swagger.json文件
  protoc -I . --openapiv2_out ./ --openapiv2_opt logtostderr=true message.proto
  5.拷贝swagger ui到third_party/swagger-ui目录下，安装go-bindata，并将swagger-ui转换为Go源码文件
  go-bindata --nocompress -pkg swagger -o pkg/swagger/datafile.go third_party/swagger-ui/...
*/
func TestGrpcGatewayHttpServer(t *testing.T) {
	startH2CGrpcGatewayHTTPServer()
}

/**
  H2C模式：无需认证
*/
func startH2CGrpcGatewayHTTPServer() {
	//监听端口
	listen, err := net.Listen("tcp", GrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	//将SimpleHello注册到grpc中
	//pb.RegisterHelloServer(s, &grpcServer{})
	gw.RegisterHelloServer(s, &grpcServer{})
	log.Println("开启Grpc服务：" + GrpcAddress)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}
}

/**
  H2模式：需要认证
*/
func startH2GrpcGatewayHTTPServer() {
	//监听端口
	listen, err := net.Listen("tcp", GrpcAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//使用证书建立grpc服务
	creds, _ := credentials.NewServerTLSFromFile(CertFile, KeyFile)
	options := []grpc.ServerOption{grpc.Creds(creds)}
	s := grpc.NewServer(options...)

	//将SimpleHello注册到grpc中
	pb.RegisterHelloServer(s, &grpcServer{})
	log.Println("开启Grpc服务：" + GrpcAddress)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}
}
