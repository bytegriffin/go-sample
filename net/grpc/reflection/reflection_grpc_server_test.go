package multiplex

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc"

	hpb "go-sample/net/grpc/reflection/proto"
	epb "go-sample/net/grpc/reflection/proto2"
)

const (
	ServerAddress = "127.0.0.1:50051"
)

type helloServer struct {
	hpb.UnimplementedHelloServer
}

func (s *helloServer) SayHello(ctx context.Context, in *hpb.HelloRequest) (*hpb.HelloResponse, error) {
	log.Printf("HelloServer received：%v", in.GetName())
	return &hpb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

type echoServer struct {
	epb.UnimplementedEchoServer
}

func (s *echoServer) Echo(ctx context.Context, in *epb.EchoRequest) (*epb.EchoResponse, error) {
	log.Printf("EchoServer received：%v", in.GetName())
	return &epb.EchoResponse{Code: 200, Message: "echo，" + in.GetName()}, nil
}

/**
  gRPC反射服务：在gRPC Server上提供有关可公开访问的gRPC服务的信息的服务。
  使用gRPC官方提供了cli工具
  查询gRPC服务：grpc_cli ls localhost:50051
  查询具体服务命令：grpc_cli ls localhost:50051 Hello -
  列出服务的rpc方法：grpc_cli ls localhost:50051 Hello.SayHello -l
  检查消息类型：grpc_cli type localhost:50051 Hello.HelloRequest
  调用远程方法：grpc_cli call localhost:50051 SayHello "name: 'gRPC CLI'"

  https://github.com/grpc/grpc-go/blob/master/Documentation/server-reflection-tutorial.md#grpc-cli
  https://github.com/grpc/grpc/blob/master/doc/server-reflection.md
*/
func TestReflectionGrpcServer(t *testing.T) {
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// 注册hello server
	hpb.RegisterHelloServer(s, &helloServer{})
	// 注册echo server
	epb.RegisterEchoServer(s, &echoServer{})

	// 注册反射到gRPC Server上
	reflection.Register(s)

	log.Println("Listen on " + ServerAddress)

	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
