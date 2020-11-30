package serverside

import (
	"context"
	pb "go-sample/net/grpc/clientside/proto"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

const (
	serverAddress = "127.0.0.1:50053"
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
	imageFile, err2 := os.Open("logo.jpg")
	if err2 != nil {
		log.Fatalln(err.Error())
	}
	defer imageFile.Close()

	// 设置metadata数据
	md := metadata.New(map[string]string{"no": "123"})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, md)

	stream, err3 := c.UploadPhoto(ctx)
	if err3 != nil {
		log.Fatalln(err)
	}

	for {
		chunk := make([]byte, 8*1024)
		chunkSize, err := imageFile.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err.Error())
		}
		if chunkSize < len(chunk) {
			chunk = chunk[:chunkSize]
		}
		stream.Send(&pb.HelloRequest{Data: chunk})
		// 为了突出效果，可以设置间隔时间来分块发送
		time.Sleep(1 * time.Second)
	}
	res, err4 := stream.CloseAndRecv()
	if err4 != nil {
		log.Fatal(err.Error())
	}
	log.Println(res.Message)

}
