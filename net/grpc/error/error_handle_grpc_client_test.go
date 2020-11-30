package error

import (
	"context"
	pb "go-sample/net/grpc/error/proto"
	"log"
	"testing"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func TestErrorHandleGrpcClient(t *testing.T) {
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)

	req := &pb.HelloRequest{Id: "1231123", Name: "go-grpc"}
	res, resError := c.SayHello(context.Background(), req)
	// 这里进行Server端的错误处理
	if resError != nil {
		errStatus, _ := status.FromError(resError)
		log.Printf("Error Code：%v，Error Message：%v", errStatus.Code(), errStatus.Message())
		if codes.InvalidArgument == errStatus.Code() {
			log.Println("=================InvalidArgument Error===============")
		}
	}
	log.Printf("Client receive：%v", res.Message)
}
