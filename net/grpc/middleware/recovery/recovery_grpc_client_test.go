package validator

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	pb "go-sample/net/grpc/middleware/proto"
	"log"
	"testing"

	"google.golang.org/grpc/grpclog"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

func unaryLogClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("====== before unary log client interceptor ======")
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("====== after unary log client interceptor ======")
	return err
}

func TestRecoveryRetryGrpcClient(t *testing.T) {
	creds, err := credentials.NewClientTLSFromFile(CertFile, "127.0.0.1")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	retryOps := []grpc_retry.CallOption{
		//最大重试次数
		grpc_retry.WithMax(10),
		//重试间隔
		grpc_retry.WithPerRetryTimeout(time.Second * 2),
		//退避时间
		grpc_retry.WithBackoff(grpc_retry.BackoffLinearWithJitter(time.Second/2, 0.2)),
	}
	retryInterceptor := grpc_retry.UnaryClientInterceptor(retryOps...)

	opts := []grpc.DialOption{
		// 1. TLS Credential
		grpc.WithTransportCredentials(creds),
		// 2. Client Unary Interceptors
		//这里用的grpc_middleware的ChainUnaryClient，而不是grpc的WithChainUnaryInterceptor
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(
				retryInterceptor,
				unaryLogClientInterceptor,
			),
		),
	}

	// 获取连结
	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewValidateServiceClient(conn)

	// 调用通讯方法
	req := &pb.ValidateRequest{Id: 199, Name: "a"}
	res, err := c.ValidatorRPC(context.Background(), req)

	if err != nil {
		grpclog.Fatalln(err)
	}
	log.Printf("Client receive：%v", res.Message)
}
