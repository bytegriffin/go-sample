package serverside

import (
	pb "go-sample/net/grpc/serverside/proto"
	"log"
	"net"
	"strconv"
	"testing"
	"time"

	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:50052"
)

type serverSideServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *serverSideServer) GetAll(req *pb.HelloRequest, stream pb.HelloService_GetAllServer) error {
	for n := 0; n < 5; n++ {
		err := stream.Send(&pb.HelloResponse{
			Code:    200,
			Message: "hello" + strconv.Itoa(n),
		})
		// 为了更好地体验效果，增加时间间隔来发送
		time.Sleep(2 * time.Second)
		if err != nil {
			return err
		}
	}
	return nil
}

// Server-side Streaming RPC，服务端流式RPC，客户端发起一次普通PRC，服务端响应并多次地流式发送给客户端
func TestServerSideGrpcServer(t *testing.T) {
	//监听端口
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//使用证书建立grpc服务
	//creds, _ := credentials.NewServerTLSFromFile("server.pem", "key.pem")
	//options := []grpc.ServerOption{grpc.Creds(creds)}
	//s := grpc.NewServer(options ...)

	//不使用证书建立grpc服务
	s := grpc.NewServer()

	//将SimpleHello注册到grpc中
	pb.RegisterHelloServiceServer(s, &serverSideServer{})
	log.Println("Listen on " + address)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
