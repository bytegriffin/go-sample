package validator

import (
	"context"
	pb "go-sample/net/grpc/middleware/proto"
	"log"
	"net"
	"testing"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	Address  = "127.0.0.1:50061"
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

/**
  服务端验证
  编译validator.proto文件命令：protoc  --proto_path=. --go_out=. --go-grpc_out=. --govalidators_out=.  validator.proto
*/
func TestValidateInterceptorGrpcServer(t *testing.T) {
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
			grpc_validator.UnaryServerInterceptor(),
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
