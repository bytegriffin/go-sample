package name_resolving

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/name_resolving/proto"
)

const (
	ServerAddress = "127.0.0.1:51051"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
	address string
}

func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received ：%v", in.GetName())
	return &pb.HelloResponse{Message: fmt.Sprintf("%s (from %s)", in.Name, s.address)}, nil
}

/**
  名称解析器：可以看成是Map[server-name][]backend-ip-address，
  一个常用的名称解析程序的例子就是DNS。
  https://github.com/grpc/grpc/blob/master/doc/naming.md
*/
func TestNameResolvingGrpcServer(t *testing.T) {
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterHelloServer(s, &unaryServer{address: ServerAddress})

	log.Println("Listen on " + ServerAddress)

	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
