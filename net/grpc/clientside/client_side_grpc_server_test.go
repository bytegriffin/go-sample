package serverside

import (
	pb "go-sample/net/grpc/clientside/proto"
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:50053"
)

type clientSideServer struct {
	pb.UnimplementedHelloServiceServer
}

func (s *clientSideServer) UploadPhoto(stream pb.HelloService_UploadPhotoServer) error {
	// 设置metadata数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		log.Printf("Data: %s\n", md["no"][0])
	}

	var image []byte
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			log.Printf("File Size：%d\n", len(image))
			return stream.SendAndClose(&pb.HelloResponse{Message: "image has received."})
		}
		if err != nil {
			return err
		}
		log.Printf("File received：%d\n", len(data.Data))
		image = append(image, data.Data...)
	}

	return nil
}

// Client-side Streaming RPC，客户端流式RPC，客户端多次发起流式PRC，服务端响应一次给客户端
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
	pb.RegisterHelloServiceServer(s, &clientSideServer{})
	log.Println("Listen on " + address)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
