package keepalive

import (
	"context"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"

	"google.golang.org/grpc/keepalive"

	pb "go-sample/net/grpc/keepalive/proto"
)

// keepalive.ClientParameters 用于在Client端设置keepalive的参数
// 该client端参数需要与server端的参数一起配合使用，否则会因为彼此配置不正确进而导致连接失败。
var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // 如果没有激活，则每10秒发送一次ping
	Timeout:             time.Second,      // 在连接停止之前，ping有1秒的ack。
	PermitWithoutStream: true,             // 在没有激活Stream的情况下client是否会发送ping
}

func TestKeepaliveGrpcClient(t *testing.T) {
	conn, err := grpc.Dial(Address, grpc.WithInsecure(), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	log.Println("Performing unary request")
	res, err := c.SayHello(ctx, &pb.HelloRequest{Name: "keepalive demo"})
	if err != nil {
		log.Fatalf("unexpected error from UnaryEcho: %v", err)
	}
	log.Println("RPC response:", res)
	select {} // 使用GODEBUG=http2debug=2可以观察ping帧和空闲状态。

}
