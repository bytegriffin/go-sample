package unary

import (
	"context"
	pb "go-sample/net/grpc/compression/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
)

const (
	ServerAddress = "127.0.0.1:50151"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", in.GetName())
	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

/**
  实现Gzip压缩
  Server端只需要导入gzip包即可实现自动注册
  Client端口需要显式代码设置
*/
func TestCompressionUnaryGrpcServer(t *testing.T) {
	//监听端口
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	//将unaryServer注册到grpc中
	pb.RegisterHelloServer(s, &unaryServer{})
	log.Println("Listen on " + ServerAddress)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
