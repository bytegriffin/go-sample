package serverside

import (
	"context"
	pb "go-sample/net/grpc/serverside/proto"
	"io"
	"log"
	"testing"

	"google.golang.org/grpc"
)

const (
	serverAddress = "127.0.0.1:50052"
)

func TestServerSideGrpcClient(t *testing.T) {

	// 使用证书连结grpc服务
	//creds, _ := credentials.NewClientTLSFromFile("server.pem", "")
	//options := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock(),grpc.WithTransportCredentials(creds)}
	//conn, err := grpc.Dial(serverAddress, options ...)
	// 不使用证书连结grpc服务
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloServiceClient(conn)

	// 调用通讯方法
	req := &pb.HelloRequest{Name: "go-grpc"}
	res, err2 := c.GetAll(context.Background(), req)

	if err2 != nil {
		log.Fatalln(err2)
	}
	for {
		res, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		log.Printf("Client receive：%v", res)
	}

}
