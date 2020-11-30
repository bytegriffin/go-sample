package serverside

import (
	"context"
	pb "go-sample/net/grpc/bidirection/proto"
	"io"
	"log"
	"testing"

	"google.golang.org/grpc"
)

const (
	serverAddress = "127.0.0.1:50054"
)

func TestBidirectionalGrpcClient(t *testing.T) {

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
	//req := &pb.HelloRequest{Name: "go-grpc"}
	stream, err2 := c.SaveAll(context.Background())
	if err2 != nil {
		log.Fatalln(err2)
	}
	finishChannel := make(chan struct{})

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				finishChannel <- struct{}{}
				break
			}
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Printf("Client receive：%v", res)
		}
	}()

	for n := 1; n < 6; n++ {
		req := &pb.HelloRequest{Name: "asdf", Id: int32(n)}
		err := stream.Send(req)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("Client Send： %v\n", req)
	}
	stream.CloseSend()
	<-finishChannel
}
