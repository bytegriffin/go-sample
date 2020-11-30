package keepalive

import (
	"context"
	"log"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"

	"google.golang.org/grpc/keepalive"

	pb "go-sample/net/grpc/keepalive/proto"
)

const (
	Address = ":50522"
)

// keepalive.EnforcementPolicy 用于在Server端设置keepalive强制策略
// Server端将关闭与违反此策略的Client端的连接。
var kaep = keepalive.EnforcementPolicy{
	MinTime:             5 * time.Second, // server端允许client端ping的最小间隔时间，如果一个client端每ping一次不超过5秒，那么server端就会终止连接。
	PermitWithoutStream: true,            // 在没有激活Stream的情况下是否允许client端ping
}

// keepalive.ServerParameters 用于在Server端设置keepalive和max_age的参数
var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     15 * time.Second, // 最大连接的闲置时间。如果一个client闲置了15秒，则会发送一个h2 Goaway frame。
	MaxConnectionAge:      30 * time.Second, // 最大连接激活时间。如果任何一个连接的活动时间超过30秒，则会发送一个h2 Goaway frame。
	MaxConnectionAgeGrace: 5 * time.Second,  // 最大连接等待时间。当被强制关闭连接之前，允许等待5秒时间挂起RPC。
	Time:                  5 * time.Second,  // 如果client端空闲了5秒，则Server端将会主动ping客户端以确保连接仍处于激活状态。
	Timeout:               1 * time.Second,  // 在连接停止之前，ping有1秒的ack。
}

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (s *unaryServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", req.GetName())
	return &pb.HelloResponse{Message: req.GetName()}, nil
}

/**
  gRPC保活机制：
  gRPC在传输上发送http2 ping来检测连接是否断开。
  如果ping在一段时间内没有被对方确认，连接将被关闭。
*/
func TestKeepaliveGrpcServer(t *testing.T) {
	lis, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	pb.RegisterHelloServer(s, &unaryServer{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
