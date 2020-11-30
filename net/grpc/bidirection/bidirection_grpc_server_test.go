package serverside

import (
	pb "go-sample/net/grpc/bidirection/proto"
	"io"
	"log"
	"net"
	"strconv"
	"testing"

	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:50054"
)

type bidirectionalServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *bidirectionalServer) SaveAll(stream pb.HelloService_SaveAllServer) error {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("Server Receive %v", res)

		stream.Send(&pb.HelloResponse{Message: strconv.Itoa(int(res.Id)) + " has received."})
	}
	return nil
}

// Bidirectional Streaming RPC，双向流式RPC，客户端发起流式PRC请求，服务端同样以流式返回给客户端
func TestBidirectionalGrpcServer(t *testing.T) {
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
	pb.RegisterHelloServiceServer(s, &bidirectionalServer{})
	log.Println("Listen on " + address)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
