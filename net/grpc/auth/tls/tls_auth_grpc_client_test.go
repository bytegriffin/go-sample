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

func TestTlsGrpcClient(t *testing.T) {
	// TLS连接，go1.15版本开始废弃CommonName，推荐使用SANs证书，因此serverName不要填CommonName
	creds, err := credentials.NewClientTLSFromFile("../server.pem", "127.0.0.1")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}
	// 获取连结
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: #{err}")
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用通讯方法
	req := &pb.HelloRequest{Name: "tls-grpc"}
	res, err := c.SayHello(context.Background(), req)

	if err != nil {
		grpclog.Fatalln(err)
	}
	log.Printf("Client receive：%v", res.Message)
}
