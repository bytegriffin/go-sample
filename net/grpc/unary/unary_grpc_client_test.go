package unary

import (
	"context"
	pb "go-sample/net/grpc/unary/proto"
	"log"
	"testing"

	"google.golang.org/grpc"
)

func TestUnaryGrpcClient(t *testing.T) {

	// 使用证书连结grpc服务
	//creds, _ := credentials.NewClientTLSFromFile("server.pem", "")
	//options := []grpc.DialOption{grpc.WithBlock(),grpc.WithTransportCredentials(creds)}
	//conn, err := grpc.Dial(serverAddress, options ...)
	// 不使用证书连结grpc服务
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用通讯方法
	req := &pb.HelloRequest{Name: "go-grpc"}
	res, err := c.SayHello(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Client receive：%v", res.Message)
}
