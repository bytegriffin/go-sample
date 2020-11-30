package unary

import (
	"context"

	pb "go-sample/net/grpc/interceptor/proto"
	"log"
	"testing"

	"google.golang.org/grpc/grpclog"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

func unaryAuthClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("****** before unary auth client interceptor ******")
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("****** after unary auth client interceptor ******")
	return err
}

func unaryLogClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("====== before unary log client interceptor ======")
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("====== after unary log client interceptor ======")
	return err
}

func TestUnaryInterceptorGrpcClient(t *testing.T) {
	creds, err := credentials.NewClientTLSFromFile(CertFile, "127.0.0.1")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	opts := []grpc.DialOption{
		// 1. TLS Credential
		grpc.WithTransportCredentials(creds),
		// 2. Client Unary Interceptors
		grpc.WithChainUnaryInterceptor(
			unaryAuthClientInterceptor,
			unaryLogClientInterceptor,
		),
	}

	// 获取连结
	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用通讯方法
	req := &pb.HelloRequest{Name: "unary-interceptor"}
	res, err := c.SayHello(context.Background(), req)

	if err != nil {
		grpclog.Fatalln(err)
	}
	log.Printf("Client receive：%v", res.Message)
}
