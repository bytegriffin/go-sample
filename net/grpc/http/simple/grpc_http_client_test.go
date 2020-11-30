package http

import (
	"context"
	pb "go-sample/net/grpc/http/simple/proto"
	"log"
	"testing"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"google.golang.org/grpc"
)

/**
  Grpc开启Https接口
*/
func TestGrpcHttpClient(t *testing.T) {
	creds, err := credentials.NewClientTLSFromFile("server.pem", "127.0.0.1")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	conn, err := grpc.Dial(port, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Connection Failed：%v", err)
	}
	defer conn.Close()

	client := pb.NewHelloClient(conn)
	res, err1 := client.SayHello(context.Background(), &pb.HelloRequest{Name: "test"})
	if err1 != nil {
		log.Fatalf("request failed：%v", err)
	}
	log.Println(res.Message)
}
