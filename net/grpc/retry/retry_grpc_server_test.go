package retry

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "go-sample/net/grpc/retry/proto"
)

type failingServer struct {
	pb.UnimplementedHelloServer
	mu sync.Mutex

	reqCounter uint //请求次数
	reqModulo  uint //取模
}

const (
	Address = ":50152"
)

// 该方法成功调用reqModule次RPC，失败reqModule-1次
func (s *failingServer) maybeFailRequest() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reqCounter++
	//判断请求数量是否是4的倍数
	if (s.reqModulo > 0) && (s.reqCounter%s.reqModulo == 0) {
		return nil
	}

	return status.Errorf(codes.Unavailable, "maybeFailRequest: failing it")
}

func (s *failingServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	if err := s.maybeFailRequest(); err != nil {
		log.Println("request failed count:", s.reqCounter)
		return nil, err
	}

	log.Println("request succeeded count:", s.reqCounter)
	return &pb.HelloResponse{Message: req.GetName()}, nil
}

/**
  gRPC支持两种retry策略：重试策略（重试失败的RPC）、对冲策略（并行地多次发送相同的RPC）。
  但是两者不能同时进行。两种都是以service config进行的：
  重试策略例子：
  "retryPolicy": {
	  "maxAttempts": 4,
	  "initialBackoff": "0.1s",
	  "maxBackoff": "1s",
	  "backoffMultiplier": 2,
	  "retryableStatusCodes": [
		"UNAVAILABLE"
	  ]
	}
  对冲策略例子：
  "hedgingPolicy": {
	 "maxAttempts": 4,
	 "hedgingDelay": "0.5s",
	 "nonFatalStatusCodes": [
		"UNAVAILABLE",
		"INTERNAL",
		"ABORTED"
	 ]
  }

  https://github.com/grpc/proposal/blob/master/A6-client-retries.md
*/
func TestRetryGrpcServer(t *testing.T) {
	lis, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("listen on address", Address)

	s := grpc.NewServer()

	// 配置服务端每接收client端请求四次，便成功执行一次PRC调用
	failingService := &failingServer{
		reqCounter: 0,
		reqModulo:  4,
	}

	pb.RegisterHelloServer(s, failingService)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
