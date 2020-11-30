package wait_for_ready

import (
	"context"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/wait_for_ready/proto"
)

const (
	address = "localhost:50051"
)

type server struct {
	pb.UnimplementedHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: req.GetName()}, nil
}

/**
  在gRPC中调用“wait for ready”
*/
func TestWaitForReadyGrpc(t *testing.T) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)

	var wg sync.WaitGroup
	wg.Add(3)

	// 情况1：如果没用启动"Wait for ready"，那么RPC会立即失败，将返回"Unavailable"错误代码.
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Wait for ready isn't enabled"})

		got := status.Code(err)
		log.Printf("[1] wanted = %v, got = %v\n", codes.Unavailable, got)
	}()

	// 情况2：如果启动"Wait for ready"，则RPC将等待服务器
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Wait for ready enabled"}, grpc.WaitForReady(true))

		got := status.Code(err)
		log.Printf("[2] wanted = %v, got = %v\n", codes.OK, got)
	}()

	// 情况3：如果启动"Wait for ready" 但是超时，那么将在超时之后失败。
	go func() {
		defer wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Wait for ready enabled and context dies"}, grpc.WaitForReady(true))

		got := status.Code(err)
		log.Printf("[3] wanted = %v, got = %v\n", codes.DeadlineExceeded, got)
	}()

	// 需要延迟3秒启动gRPC服务
	time.Sleep(3 * time.Second)
	go func() {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pb.RegisterHelloServer(s, &server{})

		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Wait()
}
