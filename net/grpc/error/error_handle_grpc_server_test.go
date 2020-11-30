package error

import (
	"context"
	pb "go-sample/net/grpc/error/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

const (
	ServerAddress = "127.0.0.1:50091"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

// 该实现方法中带错误处理，比如可以用作Server端验证
func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received Id Value：%v", in.GetId())
	if len(in.GetId()) > 5 {
		log.Println("-----------------------------")
		return nil, status.Errorf(codes.InvalidArgument, "Length of `Id` cannot be more than 5 characters")
	}
	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

// gRPC错误处理
func TestErrorHandleGrpcServer(t *testing.T) {
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterHelloServer(s, &unaryServer{})
	log.Println("Listen on " + ServerAddress)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
