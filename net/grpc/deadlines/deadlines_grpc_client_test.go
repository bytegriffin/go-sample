package unary

import (
	"context"
	pb "go-sample/net/grpc/unary/proto"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

func TestDeadlinesGrpcClient(t *testing.T) {

	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)
	//时间间隔要超过Server端的处理时间才能执行成功
	//成功
	callWithDeadlines(c, 10*time.Second)
	log.Println("=================")
	//失败
	callWithDeadlines(c, 4*time.Second)
}

func callWithDeadlines(c pb.HelloClient, timeout time.Duration) {
	//设置客户端超时时间
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req := &pb.HelloRequest{Name: "grpc-deadlines"}
	res, err := c.SayHello(ctx, req)
	//获取Server端返回的错误
	if err != nil {
		//获取错误状态
		status, ok := status.FromError(err)
		if ok {
			//判断是否为调用超时
			if status.Code() == codes.DeadlineExceeded {
				log.Fatalln("Error：超时!")
			}
		}
		log.Fatalf("Call SayHello err: %v，Error Code：%v，Error Message：%v", err, status.Code(), status.Message())
	}
	log.Printf("Client receive：%v", res.GetMessage())
}
