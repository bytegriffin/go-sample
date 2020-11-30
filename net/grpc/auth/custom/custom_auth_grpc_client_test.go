package tls

import (
	"context"
	pb "go-sample/net/grpc/auth/proto"
	"log"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	// OpenTLS 是否开启TLS认证
	OpenTLS = true
)

// 自定义认证
type customCredential struct {
	AppID  string
	AppKey string
}

// 获取当前请求的元数据
func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_id":  c.AppID,
		"app_key": c.AppKey,
	}, nil
}

// 自定义认证是否开启TLS
func (c customCredential) RequireTransportSecurity() bool {
	return OpenTLS
}

/**
  PerRPCCredentials 是gRPC提供给自定义认证的接口，
  它包含了两个接口方法：GetRequestMetadata、RequireTransportSecurity
*/
func TestCustomAuthGrpcClient(t *testing.T) {
	custom := customCredential{
		AppID:  "custom-id",
		AppKey: "custom-key",
	}
	var err error
	var opts []grpc.DialOption
	if OpenTLS {
		// TLS连接
		creds, err := credentials.NewClientTLSFromFile(CertFile, "127.0.0.1")
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		//不能跟grpc.WithCredentialsBundle一起使用
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// TLS连接，go1.15版本开始废弃CommonName，推荐使用SANs证书，因此serverName不要填CommonName
	opts = append(opts, grpc.WithPerRPCCredentials(&custom))

	// 获取连结
	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用通讯方法
	req := &pb.HelloRequest{Name: "custom-grpc"}
	res, err := c.SayHello(context.Background(), req)

	if err != nil {
		grpclog.Fatalln(err)
	}
	log.Printf("Client receive：%v", res.Message)
}
