package profiling

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/profiling/proto"

	profiling "google.golang.org/grpc/profiling/service"
)

const (
	Address = ":50031"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("SayHello called with name %q\n", in.GetName())
	return &pb.HelloResponse{Message: in.GetName()}, nil
}

/**
  带分析功能的gRPC服务
*/
func TestProfilingGrpcServer(t *testing.T) {
	lis, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &unaryServer{})

	// 注册带分析功能的gRPC Server
	pc := &profiling.ProfilingConfig{
		Server:          s,
		Enabled:         true, //分析开关
		StreamStatsSize: 1024,
	}
	if err = profiling.Init(pc); err != nil {
		log.Printf("error calling profiling.Init: %v\n", err)
		return
	}

	s.Serve(lis)
}
