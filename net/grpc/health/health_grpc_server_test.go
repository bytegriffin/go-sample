package health

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	pb "go-sample/net/grpc/health/proto"

	healthPB "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	Address = ":50551"

	system = "" // 空字符串表示系统运行状况
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (e *unaryServer) UnaryEcho(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", req.GetName())
	return &pb.HelloResponse{
		Message: fmt.Sprintf("hello from localhost: %v", Address),
	}, nil
}

var _ pb.HelloServer = &unaryServer{}

/**
  gRPC健康检查
  Client端有两种监控Server端的运行状况：
  check()：面向流式RPC，探测服务器的运行状况。
  watch()：面向一元RPC，探测服务器的变化。

  https://github.com/grpc/proposal/blob/master/A17-client-side-health-checking.md
*/
func TestHealthGrpcServer(t *testing.T) {
	lis, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	healthCheck := health.NewServer()
	healthPB.RegisterHealthServer(s, healthCheck)
	pb.RegisterHelloServer(s, &unaryServer{})

	go func() {
		// 异步检查依赖项，并根据需要切换服务状态
		next := healthPB.HealthCheckResponse_SERVING

		for {
			healthCheck.SetServingStatus(system, next)

			if next == healthPB.HealthCheckResponse_SERVING {
				next = healthPB.HealthCheckResponse_NOT_SERVING
			} else {
				next = healthPB.HealthCheckResponse_SERVING
			}

			time.Sleep(time.Second * 5)
		}
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
