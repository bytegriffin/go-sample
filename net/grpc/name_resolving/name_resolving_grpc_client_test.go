package name_resolving

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc/resolver"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/name_resolving/proto"
)

const (
	exampleScheme      = "example"
	exampleServiceName = "resolver.example.grpc.io"

	backendAddr = "localhost:51051"
)

func callUnary(c pb.HelloClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: message})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println(r.Message)
}

func makeRPCs(cc *grpc.ClientConn, n int, msg string) {
	hwc := pb.NewHelloClient(cc)
	for i := 0; i < n; i++ {
		callUnary(hwc, fmt.Sprintf("this is %v name_resolving", msg))
	}
}

/**
  本例子中示范了ClientConn如何选择不同的名称解析器。
  server端是在70051端口上工作，创建了两个客户端：
  第一个是通过passthrough:///localhost:51051进行连接；
  第二个是通过example:///resolver.example.grpc.io进行连接，
  它们最终都会连接到server端。也就是说：第二个需要名称解析器将域名进行解析后才能正确连接。
*/
func TestNameResolvingGrpcClient(t *testing.T) {
	passthroughConn, err := grpc.Dial(
		fmt.Sprintf("passthrough:///%s", backendAddr), // Dial to "passthrough:///localhost:51051"
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer passthroughConn.Close()

	log.Printf("--- calling Hello/SayHello to \"passthrough:///%s\"\n", backendAddr)
	makeRPCs(passthroughConn, 10, "passthrough")

	log.Println("==================================================================================")

	exampleConn, err := grpc.Dial(
		fmt.Sprintf("%s:///%s", exampleScheme, exampleServiceName), // Dial to "example:///resolver.example.grpc.io"
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer exampleConn.Close()

	log.Printf("--- calling Hello/SayHello to \"%s:///%s\"\n", exampleScheme, exampleServiceName)
	makeRPCs(exampleConn, 10, "example")
}

// 以下代码是专门为第二个域名的名称解析器
// 即：将 resolver.example.grpc.io => localhost:51051
type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: {backendAddr},
		},
	}
	r.start()
	return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return exampleScheme }

// Resolver(https://godoc.org/google.golang.org/grpc/resolver#Resolver).
type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}

func init() {
	// 注册 ResolverBuilder
	resolver.Register(&exampleResolverBuilder{})
}
