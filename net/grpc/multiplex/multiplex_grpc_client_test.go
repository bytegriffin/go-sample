package multiplex

import (
	"context"
	hpb "go-sample/net/grpc/multiplex/proto"
	epb "go-sample/net/grpc/multiplex/proto2"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func callHello(conn *grpc.ClientConn) {
	hc := hpb.NewHelloClient(conn)
	req := &hpb.HelloRequest{Name: "grpc-hello"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := hc.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("client.SayHello(_) = _, %v", err)
	}
	log.Printf("HelloClient receive：%v", res.Message)
}

func callEcho(conn *grpc.ClientConn) {
	ec := epb.NewEchoClient(conn)
	req := &epb.EchoRequest{Name: "grpc-echo"}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := ec.Echo(ctx, req)
	if err != nil {
		log.Fatalf("client.Echo(_) = _, %v", err)
	}
	log.Printf("EchoClient receive：%v", res.Message)
}

func TestMultiplexGrpcClient(t *testing.T) {
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %\n", err)
	}
	defer conn.Close()

	callHello(conn)
	log.Println("--------------------------------------")
	callEcho(conn)
}
