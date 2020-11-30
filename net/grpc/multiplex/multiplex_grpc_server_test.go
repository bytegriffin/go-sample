package multiplex

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"

	hpb "go-sample/net/grpc/multiplex/proto"
	epb "go-sample/net/grpc/multiplex/proto2"
)

const (
	ServerAddress = "127.0.0.1:60051"
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
  gRPC多路复用：将两个服务注册到同一个gRPC Server中
*/
func TestMultiplexGrpcServer(t *testing.T) {
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	// 注册hello server
	hpb.RegisterHelloServer(s, &helloServer{})
	// 注册echo server
	epb.RegisterEchoServer(s, &echoServer{})

	log.Println("Listen on " + ServerAddress)

	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
