package unary

import (
	"context"
	pb "go-sample/net/grpc/unary/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
)

const (
	ServerAddress = "127.0.0.1:50051"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", in.GetName())
	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

// Unary RPC，一元RPC，客户端发起一次PRC，服务端响应并返回给客户端
func TestUnaryGrpcServer(t *testing.T) {
	//监听端口
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//使用证书建立grpc服务
	//creds, _ := credentials.NewServerTLSFromFile("server.pem", "key.pem")
	//options := []grpc.ServerOption{grpc.Creds(creds)}
	//s := grpc.NewServer(options ...)

	//不使用证书建立grpc服务
	s := grpc.NewServer()

	//将unaryServer注册到grpc中
	pb.RegisterHelloServer(s, &unaryServer{})
	log.Println("Listen on " + ServerAddress)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
