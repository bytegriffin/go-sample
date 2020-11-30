package grpcurl

import (
	"context"
	pb "go-sample/net/grpc/grpcurl/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type testServer struct {
	pb.UnimplementedHelloServer
}

func (*testServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", req.GetName())
	return &pb.HelloResponse{Message: "hello"}, nil
}

/**
  GrpcUrl是一个检验Grpc运行状态的命令行工具
  查询服务列表：grpcurl -plaintext 127.0.0.1:8080 list
  查询具体服务：grpcurl -plaintext 127.0.0.1:8080 Hello/SayHello
  查看方法定义：grpcurl -plaintext 127.0.0.1:8080 describe Hello.SayHello
  调用服务方法：grpcurl -plaintext -d '{"name": "gopher"}' 127.0.0.1:8080 Hello/SayHello
  注意：windows下调用服务需要转义json字符串：grpcurl -plaintext -d "{\"name\": \"gopher\"}" 127.0.0.1:8080 Hello/SayHello
*/
func TestGrpcurlServer(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Failed to listen：%v", err)
	}
	s := grpc.NewServer()
	//注册GrpcUrl所需的reflection服务
	reflection.Register(s)

	pb.RegisterHelloServer(s, &testServer{})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve：%v", err)
	}

}
