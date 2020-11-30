package cancellation

import (
	"context"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "go-sample/net/grpc/cancellation/proto"
)

func sendMessage(stream pb.HelloService_SaveAllClient, msg string) error {
	log.Printf("sending message %q\n", msg)
	return stream.Send(&pb.HelloRequest{Name: msg})
}

func recvMessage(stream pb.HelloService_SaveAllClient, wantErrCode codes.Code) {
	res, err := stream.Recv()
	if status.Code(err) != wantErrCode {
		log.Fatalf("stream.Recv() = %v, %v; want _, status.Code(err)=%v", res, err, wantErrCode)
	}
	if err != nil {
		log.Printf("stream.Recv() returned expected error %v\n", err)
		return
	}
	log.Printf("received message %q\n", res.GetMessage())
}

func TestCancelGrpcClient(t *testing.T) {
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewHelloServiceClient(conn)

	// 初始化发送超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	stream, err := c.SaveAll(ctx)
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	if err := sendMessage(stream, "hello"); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}
	if err := sendMessage(stream, "world"); err != nil {
		log.Fatalf("error sending on stream: %v", err)
	}

	recvMessage(stream, codes.OK)
	recvMessage(stream, codes.OK)

	//取消context
	cancel()

	// 此次信息Client端是高数Server端口要取消context
	sendMessage(stream, "closed")

	// 此次接收将会失败
	recvMessage(stream, codes.Canceled)
}
