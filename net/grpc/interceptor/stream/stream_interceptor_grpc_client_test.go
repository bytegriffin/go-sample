package unary

import (
	"context"
	pb "go-sample/net/grpc/interceptor/proto"
	"io"
	"log"
	"strconv"
	"testing"

	"google.golang.org/grpc/grpclog"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

func streamAuthClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("****** before stream auth client interceptor ******")
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	log.Printf("****** after stream auth client interceptor ******")
	return s, nil
}

func streamLogClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	log.Printf("====== before stream log client interceptor ======")
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	log.Printf("====== after stream log client interceptor ======")
	return s, nil
}

func TestStreamInterceptorGrpcClient(t *testing.T) {
	creds, err := credentials.NewClientTLSFromFile(CertFile, "127.0.0.1")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	opts := []grpc.DialOption{
		// 1. TLS Credential
		grpc.WithTransportCredentials(creds),
		// 2. Client Unary Interceptors
		grpc.WithChainStreamInterceptor(
			streamAuthClientInterceptor,
			streamLogClientInterceptor,
		),
	}

	// 获取连结
	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	// 初始化客户端
	c := pb.NewHelloClient(conn)

	// 调用通讯方法
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
		req := &pb.HelloRequest{Name: "asdf", Id: strconv.Itoa(n)}
		err := stream.Send(req)
		if err != nil {
			log.Fatalln(err.Error())
		}
		log.Printf("Client Send： %v\n", req)
	}
	stream.CloseSend()
	<-finishChannel
}
