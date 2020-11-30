package unary

import (
	"context"
	pb "go-sample/net/grpc/deadlines/proto"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

const (
	ServerAddress = "127.0.0.1:50081"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", in.GetName())
	// 模拟多次请求，假设请求一次是5秒
	for i := 0; i < 5; i++ {
		// 判断client端在请求期间（即：没超时之前），是否提前取消请求
		// 比如：client端请求期间中断client端，此时Server端应该停止正在进行的操作，避免资源浪费
		if ctx.Err() == context.Canceled {
			log.Println("The client canceled the request!")
			return nil, status.Errorf(codes.Canceled, "The client canceled the request")
		}
		//模拟环境加入睡眠时间
		time.Sleep(1 * time.Second)
	}
	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

//gRPC超时
func TestDeadlinesGrpcServer(t *testing.T) {
	//监听端口
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//不使用证书建立grpc服务
	s := grpc.NewServer()

	//将SimpleHello注册到grpc中
	pb.RegisterHelloServer(s, &unaryServer{})
	log.Println("Listen on " + ServerAddress)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
