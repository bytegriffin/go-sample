package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/xds/proto"
)

var clientHelp = flag.Bool("help", false, "Print usage information")

const (
	defaultTarget = "localhost:50051"
	defaultName   = "world"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `
Usage: client [name [target]]
  name
        The name you wish to be greeted by. Defaults to %q
  target
        The URI of the server, e.g. "xds:///helloworld-service". Defaults to %q
`, defaultName, defaultTarget)

		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if *clientHelp {
		flag.Usage()
		return
	}
	args := flag.Args()

	if len(args) > 2 {
		flag.Usage()
		return
	}

	name := defaultName
	if len(args) > 0 {
		name = args[0]
	}

	target := defaultTarget
	if len(args) > 1 {
		target = args[1]
	}

	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewHelloClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMessage())
}
