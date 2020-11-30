package unary

import (
	"context"
	pb "go-sample/net/grpc/compression/proto"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc/encoding/gzip"

	"google.golang.org/grpc"
)

func TestCompressionUnaryGrpcClient(t *testing.T) {
	//压缩方式一：所有的RPC请求都通过一个client端发出
	//conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure(),
	//	grpc.WithDefaultCallOptions(
	//		grpc.UseCompressor(gzip.Name),
	//	),
	//)

	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)

	//压缩方式二：具体压缩的消息字段
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req := &pb.HelloRequest{Name: "grpc-compression"}
	res, err1 := c.SayHello(ctx, req, grpc.UseCompressor(gzip.Name))
	if err1 != nil {
		log.Fatalln(err)
	}
	log.Printf("Client receive：%v", res.Message)
}
