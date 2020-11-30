package load_balancing

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/load_balancing/proto"
)

var (
	address = []string{":50051", ":50052"}
)

type ecServer struct {
	pb.UnimplementedHelloServer
	addr string
}

func (s *ecServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: fmt.Sprintf("%s (from %s)", req.Name, s.addr)}, nil
}

func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, &ecServer{addr: addr})
	log.Printf("serving on %s\n", addr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

/**
  gRPC的负载均衡机制是基于外部负载均衡器，一个外部的负载均衡器提供简单的客户端
  和一个最新的服务器列表。gRPC官方提供两种负载均衡策略：round_robin（轮询调度）、grpclb。
  如果任何一个被解析器返回的地址是均衡器地址，那么client会使用grpclb策略，否则
  client将会使用service config中配置的策略，如果client中没有在service config
  中配置负载均衡策略，那么将会默认使用pick_first，即客户端会取第一个可用服务器地址。
  https://github.com/grpc/grpc/blob/master/doc/load-balancing.md
*/
func TestLoadBalancingGrpcServer(t *testing.T) {
	var wg sync.WaitGroup
	for _, addr := range address {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			startServer(addr)
		}(addr)
	}
	wg.Wait()
}
