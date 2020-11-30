package validator

import (
	"context"
	pb "go-sample/net/grpc/middleware/proto"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	Address  = "127.0.0.1:50062"
	CertFile = "../server.pem"
	keyFile  = "../server.key"
)

type unaryServer struct {
	pb.UnimplementedValidateServiceServer
}

func (s *unaryServer) ValidatorRPC(ctx context.Context, request *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	log.Printf("Server received：%v", request.GetName())
	return &pb.ValidateResponse{Code: 200, Message: "hello，" + request.GetName()}, nil
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

// 自定义错误
func RecoveryInterceptor() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	})
}

/**
  将gRPC中的panic转成error，从而恢复程序
*/
func TestRecoveryInterceptorGrpcServer(t *testing.T) {
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
		// 2.Validator Interceptor
		grpc.ChainUnaryInterceptor(
			grpc_recovery.UnaryServerInterceptor(RecoveryInterceptor()),
			unaryLogServerInterceptor,
		),
	}
	s := grpc.NewServer(opts...)

	//将SimpleHello注册到grpc中
	pb.RegisterValidateServiceServer(s, &unaryServer{})
	log.Println("Listen on " + Address)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
