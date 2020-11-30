package tls

import (
	"context"
	pb "go-sample/net/grpc/auth/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	Address = "127.0.0.1:50055"
)

type authService struct {
	pb.UnimplementedHelloServer
}

func (h authService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", in.GetName())
	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

/**
  SSL/TLS认证方式：
  1.制作私钥
  openssl genrsa -out server.key 2048
  openssl ecparam -genkey -name secp384r1 -out server.key
  2.制作自签名公钥(x509)
  openssl req -new -x509 -sha256 -key server.key -out server.pem -days 3650
*/
func TestTlsGrpcServer(t *testing.T) {
	// 监听端口
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile("../server.pem", "../server.key")
	if err != nil {
		log.Fatalf("Failed to generate credentials %v", err)
	}

	// 实例化grpc服务，并开启认证
	s := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterHelloServer(s, &authService{})
	log.Println("Listen on " + Address)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Auth Grpc Server is failed.")
	}

}
