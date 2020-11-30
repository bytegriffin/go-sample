package debug

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"testing"
	"time"

	"google.golang.org/grpc"
	channelzService "google.golang.org/grpc/channelz/service"

	pb "go-sample/net/grpc/debug/proto"
)

var (
	Addresses          = []string{":10001", ":10002", ":10003"}
	serverChannelzPort = ":50051"
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

// slowServer是模拟server的一个延迟版的server
type slowServer struct {
	pb.UnimplementedHelloServer
}

func (s *slowServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	// Delay 100ms ~ 200ms before replying
	time.Sleep(time.Duration(100+Intn(100)) * time.Millisecond)
	return &pb.HelloResponse{Message: "Hello " + in.Name}, nil
}

/**
  Grpc官方为用户提供了两种debug方式：GrpcLog、Channelz。
  1.GrpcLog可以通过设置日志级别（Info、Warning、Error、Fatal）来显示跟踪日志。
  开启GrpcLog需要事先设置好环境变量：
  GRPC_GO_LOG_VERBOSITY_LEVEL=99 //默认值为0
  GRPC_GO_LOG_SEVERITY_LEVEL=info
  2.Channelz可以提供不同Channel级别的调试信息。每个channel表示一个DAG（有向无环图），
  它可以包含多个子channel，每个子channel还可以包含多个socket，每个channel可
  同时拥有子channel和socket，但不是同时拥有子channel和socket。

  https://grpc.io/blog/a-short-introduction-to-channelz/
  https://github.com/grpc/proposal/blob/master/A14-channelz.md
*/
func TestDebugGrpcServer(t *testing.T) {
	//1.创建server端的channelz服务
	listen, err := net.Listen("tcp", serverChannelzPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listen.Close()
	s := grpc.NewServer()
	channelzService.RegisterChannelzServiceToServer(s)
	go s.Serve(listen)
	defer s.Stop()

	//2.开启三个server，并将其中一个设置为slowServer
	var listeners []net.Listener
	var svrs []*grpc.Server
	for i := 0; i < 3; i++ {
		lis, err := net.Listen("tcp", Addresses[i])
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		listeners = append(listeners, lis)
		s := grpc.NewServer()
		svrs = append(svrs, s)
		if i == 2 {
			pb.RegisterHelloServer(s, &slowServer{})
		} else {
			pb.RegisterHelloServer(s, &server{})
		}
		go s.Serve(lis)
	}

	//3.等待用户使用Ctrl+C退出程序
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// 直到收到信号之前一直堵塞
	<-ch
	for i := 0; i < 3; i++ {
		svrs[i].Stop()
		listeners[i].Close()
	}
}
