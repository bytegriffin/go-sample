package metadata

import (
	"context"
	pb "go-sample/net/grpc/metadata/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

const (
	ServerAddress = "127.0.0.1:50051"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {

	//Server端收到RPC之后，并且在处理业务之前，向client端发送metadata
	header := metadata.Pairs("server-header", "server-header-val")
	//unary模式
	grpc.SendHeader(ctx, header)
	//stream模式
	//stream.SendHeader(header)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Printf("get metadata error")
	}
	// 注意：返回的是一个切片，而不是一个值
	if val1, ok := md["key1"]; ok {
		log.Printf("val from metadata: %v", val1[0])
	}
	if val2, ok := md["key2"]; ok {
		log.Printf("val from metadata: %v, %v", val2[0], val2[1])
	}
	if val3, ok := md["key3"]; ok {
		log.Printf("val from metadata: %v", val3[0])
	}

	log.Printf("Server received：%v", in.GetName())

	//Server端收到RPC之后，并且在处理业务之后，向client端发送metadata
	trailer := metadata.Pairs("server-trailer", "server-trailer-val")
	//unary模式
	grpc.SetTrailer(ctx, trailer)
	//stream模式
	//stream.SendHeader(header)

	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

func TestMetadataUnaryGrpcServer(t *testing.T) {
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterHelloServer(s, &unaryServer{})

	log.Println("Listen on " + ServerAddress)
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
