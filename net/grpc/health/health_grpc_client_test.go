package health

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"

	pb "go-sample/net/grpc/health/proto"

	_ "google.golang.org/grpc/health"
)

var serviceConfig = `{
	"loadBalancingPolicy": "round_robin",
	"healthCheckConfig": {
		"serviceName": ""
	}
}`

func callUnaryEcho(c pb.HelloClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{})
	if err != nil {
		log.Printf("UnaryEcho: _, %v ", err)
	} else {
		log.Printf("UnaryEcho: %v", r.GetMessage())
	}
}

func TestHealthGrpcClient(t *testing.T) {

	r := manual.NewBuilderWithScheme("whatever")
	r.InitialState(resolver.State{
		Addresses: []resolver.Address{
			{Addr: Address},
			{Addr: "localhost:50552"},
		},
	})

	address := fmt.Sprintf("%s:///unused", r.Scheme())

	options := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithResolvers(r),
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	conn, err := grpc.Dial(address, options...)
	if err != nil {
		log.Fatalf("did not connect %v", err)
	}
	defer conn.Close()

	echoClient := pb.NewHelloClient(conn)

	for {
		callUnaryEcho(echoClient)
		time.Sleep(time.Second)
	}
}
