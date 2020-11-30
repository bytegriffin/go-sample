package debug

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/channelz/service"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	pb "go-sample/net/grpc/debug/proto"
)

const (
	defaultName        = "world"
	clientChannelzPort = ":50052"
)

func TestDebugGrpcClient(t *testing.T) {
	// 1.创建client端的channelz服务
	listen, err := net.Listen("tcp", clientChannelzPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listen.Close()
	s := grpc.NewServer()
	service.RegisterChannelzServiceToServer(s)
	go s.Serve(listen)
	defer s.Stop()

	// 2.初始化 manual resolver 和 Dial
	//r := manual.NewBuilderWithScheme("whatever")

	r, rcleanup := manual.GenerateAndRegisterManualResolver()
	defer rcleanup()

	// 连接Server
	conn, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithResolvers(r), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// 设置目的地址
	r.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: Addresses[0]}, {Addr: Addresses[1]}, {Addr: Addresses[2]}}})

	c := pb.NewHelloClient(conn)

	// 连接到服务并且输出返回信息
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	// 3.发送100次SayHello
	for i := 0; i < 100; i++ {
		// 设置150ms超时
		ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
		if err != nil {
			log.Printf("could not hello: %v", err)
		} else {
			log.Printf("Hello: %s", r.Message)
		}
		cancel()
	}

	// 4.等待用户退出程序
	//除非使用CTRL+C，否则channelz数据将一直可用
	select {}
}
