package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "go-sample/net/grpc/xds/proto"
)

const (
	defaultPort = 50051
)

var serverHelp = flag.Bool("help", false, "Print usage information")

type server struct {
	pb.UnimplementedHelloServer
	serverName string
}

func newServer(serverName string) *server {
	return &server{
		serverName: serverName,
	}
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received: %v", in.GetName())
	return &pb.HelloResponse{Message: "Hello " + in.GetName() + ", from " + s.serverName}, nil
}

func determineHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Failed to get hostname: %v, will generate one", err)
		rand.Seed(time.Now().UnixNano())
		return fmt.Sprintf("generated-%03d", rand.Int()%100)
	}
	return hostname
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `
Usage: server [port [hostname]]
  port
        The listen port. Defaults to %d
  hostname
        The name clients will see in greet responses. Defaults to the machine's hostname
`, defaultPort)

		flag.PrintDefaults()
	}
}

/**
  xDS最初是Envoy获取配置信息的传输协议，xDS是查询文件或管理服务器来动态发现资源的总称，
  包含LDS、RDS、CDS、EDS以及SDS。它正在进化为Service Mesh的通用数据计划API。
*/
func main() {
	flag.Parse()
	if *serverHelp {
		flag.Usage()
		return
	}
	args := flag.Args()

	if len(args) > 2 {
		flag.Usage()
		return
	}

	port := defaultPort
	if len(args) > 0 {
		var err error
		port, err = strconv.Atoi(args[0])
		if err != nil {
			log.Printf("Invalid port number: %v", err)
			flag.Usage()
			return
		}
	}

	var hostname string
	if len(args) > 1 {
		hostname = args[1]
	}
	if hostname == "" {
		hostname = determineHostname()
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterHelloServer(s, newServer(hostname))

	reflection.Register(s)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(s, healthServer)

	log.Printf("serving on %s, hostname %s", lis.Addr(), hostname)
	s.Serve(lis)
}
