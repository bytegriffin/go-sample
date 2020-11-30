package profiling

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/profiling/proto"

	profiling "google.golang.org/grpc/profiling/service"
)

const (
	profilingPort = ":50052"
)

func setupClientProfiling() error {
	lis, err := net.Listen("tcp", profilingPort)
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return err
	}
	log.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()

	// 注册带分析功能的gRPC Server
	pc := &profiling.ProfilingConfig{
		Server:          s,
		Enabled:         true,
		StreamStatsSize: 1024,
	}
	if err = profiling.Init(pc); err != nil {
		log.Printf("error calling profsvc.Init: %v\n", err)
		return err
	}

	go s.Serve(lis)
	return nil
}

/**
  client远程开启/关闭分析功能：
  go run google.golang.org/grpc/profiling/cmd -address localhost:50031 -enable-profiling
  go run google.golang.org/grpc/profiling/cmd -address localhost:50031 -disable-profiling
*/
func TestProfilingGrpcClient(t *testing.T) {
	if err := setupClientProfiling(); err != nil {
		log.Fatalf("error setting up profiling: %v\n", err)
	}

	// 1.以同步方式连接server端
	conn, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := c.SayHello(ctx, &pb.HelloRequest{Name: "hello, profiling"})
	fmt.Printf("Server returned %q, %v\n", res.GetMessage(), err)
	if err != nil {
		log.Fatalf("Server returned error: %v", err)
	}

	log.Printf("Sleeping for 30 seconds with exposed profiling service on :%v \n", profilingPort)
	time.Sleep(30 * time.Second)
}
