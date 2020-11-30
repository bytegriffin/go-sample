package http

import (
	"context"
	pb "go-sample/net/grpc/http/simple/proto"
	"log"
	"net/http"
	"testing"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type testServer struct {
	pb.UnimplementedHelloServer
}

func (s *testServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", req.GetName())
	return &pb.HelloResponse{Message: "hello"}, nil
}

/**
  Grpc开启Https接口
*/
func TestGrpcHttpServer(t *testing.T) {
	creds, err := credentials.NewServerTLSFromFile("server.pem", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterHelloServer(s, &testServer{})

	httpServer := &http.Server{
		Addr: port,
		// Handler: mux,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Server receive %v", r)
			s.ServeHTTP(w, r)
			w.Write([]byte("Hello World, Http Protocol: " + r.Proto))
		}),
	}

	if err := httpServer.ListenAndServeTLS("server.pem", "server.key"); err != nil {
		log.Fatalf("failed to serve：%v", err)
	}

}
