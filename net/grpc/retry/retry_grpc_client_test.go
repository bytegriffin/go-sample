package retry

import (
	"context"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/retry/proto"
)

var ( //retry策略采用的是backoff算法
	retryPolicy = `{
		"methodConfig": [{
		  "name": [{"service": "net.grpc.retry"}],
		  "waitForReady": true,
		  "retryPolicy": {
			  "MaxAttempts": 4,
			  "InitialBackoff": ".01s",
			  "MaxBackoff": ".01s",
			  "BackoffMultiplier": 1.0,
			  "RetryableStatusCodes": [ "UNAVAILABLE" ]
		  }
		}]}`
)

// 使用 grpc.WithDefaultServiceConfig() 设置 service config
func retryDial() (*grpc.ClientConn, error) {
	return grpc.Dial(Address, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(retryPolicy))
}

/**
  客户端不能重写service config，但是可以在客户端中完全禁用客户端支持。
  运行前需要设置环境变量 GRPC_GO_RETRY=on
  https://github.com/grpc/grpc/blob/master/doc/service_config.md
*/
func TestRetryGrpcClient(t *testing.T) {
	conn, err := retryDial()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		if e := conn.Close(); e != nil {
			log.Printf("failed to close connection: %s", e)
		}
	}()

	c := pb.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	reply, err := c.SayHello(ctx, &pb.HelloRequest{Name: "Try and Success"})
	if err != nil {
		log.Fatalf("SayHello error: %v", err)
	}
	log.Printf("SayHello reply: %v", reply)
}
