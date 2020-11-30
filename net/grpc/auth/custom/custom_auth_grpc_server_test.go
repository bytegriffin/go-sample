package tls

import (
	"context"
	pb "go-sample/net/grpc/auth/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	Address  = "127.0.0.1:50056"
	CertFile = "../server.pem"
	KeyFile  = "../server.key"
)

type authService struct {
	pb.UnimplementedHelloServer
}

func (h authService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	// 解析metadata中的信息并验证
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	//验证各字段
	var (
		AppID  string
		AppKey string
	)
	if val, ok := md["app_id"]; ok {
		AppID = val[0]
	}
	if val, ok := md["app_key"]; ok {
		AppKey = val[0]
	}
	if AppID != "1024" || AppKey != "test" {
		return nil, status.Errorf(codes.Unauthenticated, "Token认证出错，请检查: AppId=%s, AppKey=%s", AppID, AppKey)
	}

	log.Printf("Server received：%v", in.GetName())
	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

/**
  自定义认证方式：
  Client端负责实现PerRPCCredentials两个接口方法，
  Server端用错误处理进行判断
*/
func TestCustomAuthGrpcServer(t *testing.T) {
	// 监听端口
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// TLS认证
	creds, err := credentials.NewServerTLSFromFile(CertFile, KeyFile)
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
