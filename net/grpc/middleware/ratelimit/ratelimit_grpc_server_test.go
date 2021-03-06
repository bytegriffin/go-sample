package unary

import (
	"context"
	pb "go-sample/net/grpc/middleware/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	Address  = "127.0.0.1:50071"
	CertFile = "../server.pem"
	keyFile  = "../server.key"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", in.GetName())
	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

func unaryLogServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("====== before unary log server interceptor ======")
	m, err := handler(ctx, req)
	if err != nil {
		log.Printf("failed to handler Unary RPC: %v\n", err)
	}
	log.Printf("====== after unary log server interceptor ======")
	return m, err
}

type mockLimiter struct{}

func (*mockLimiter) Limit() bool {
	return true //是否限流，true表示限流
}

type Limiter interface {
	Limit() bool
}

func UnaryServerInterceptor(limiter Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if limiter.Limit() {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected by grpc_ratelimit middleware, please retry later.", info.FullMethod)
		}
		return handler(ctx, req)
	}
}

/**
  自定义限流
*/
func TestRatelimitGrpcServer(t *testing.T) {
	//监听端口
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(CertFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load certificates: %v", err)
	}

	opts := []grpc.ServerOption{
		// 1. TLS Credential
		grpc.Creds(creds),
		// 2.Interceptor
		grpc.ChainUnaryInterceptor(unaryLogServerInterceptor, UnaryServerInterceptor(&mockLimiter{})),
	}
	s := grpc.NewServer(opts...)

	//将SimpleHello注册到grpc中
	pb.RegisterHelloServer(s, &unaryServer{})
	log.Println("Listen on " + Address)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
