package metadata

import (
	"context"
	pb "go-sample/net/grpc/metadata/proto"
	"log"
	"testing"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

/**
  有时不想让数据全部通过body来传输，可以使用metadata，
  metadata是header中的一员，metadata的key值永远是string，
  value可以是二进制或string，如果需要存储二进制的值，那么需要使用“-bin"来做key值的后缀进行标注，
  这样被标注过的二进制值将会在发送前自动被base64编码，直到被接收后才会进行解码。
*/
func TestMetadataUnaryGrpcClient(t *testing.T) {
	conn, err := grpc.Dial(ServerAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	//metadata生成方式一：会覆盖同名key的value值
	//md := metadata.New(map[string]string{"key1": "val1", "key2": "val2"})
	//metadata生成方式二：不会覆盖同名key的value值
	//md := metadata.Pairs(
	//	"key1", "val1",
	//	"key2", "val2-1", // 相同的key值会被组合成slice
	//	"key2", "val2-2",
	//	"key-bin", string([]byte{96, 102}),
	//)
	//ctx := metadata.NewOutgoingContext(ctx, md)
	//metadata生成方式三：不会覆盖同名key的value值，推荐使用该方式
	ctx := metadata.AppendToOutgoingContext(context.Background(), "key1", "v1", "key2", "v2", "key2", "v3")

	req := &pb.HelloRequest{Name: "grpc-metadata"}
	c := pb.NewHelloClient(conn)
	//Client端接收Server端发来的Metadata
	var header, trailer metadata.MD
	res, err1 := c.SayHello(ctx, req, grpc.Header(&header), grpc.Trailer(&trailer))
	if err1 != nil {
		log.Fatalln(err)
	}
	log.Printf("Client receive message: %v，header ：%v，trailer ：%v", res.GetMessage(), header["server-header"], trailer["server-trailer"])
}
