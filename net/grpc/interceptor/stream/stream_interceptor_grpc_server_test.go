package unary

import (
	pb "go-sample/net/grpc/interceptor/proto"
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
)

const (
	Address  = "127.0.0.1:50052"
	CertFile = "../server.pem"
	keyFile  = "../server.key"
)

type streamServer struct {
	pb.UnimplementedHelloServer
}

func (s *streamServer) SaveAll(stream pb.Hello_SaveAllServer) error {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		log.Printf("Server Receive %v", res)

		stream.Send(&pb.HelloResponse{Message: res.Id + " has received."})
	}
	return nil
}

func streamAuthServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("****** before stream auth server interceptor ******")
	err := handler(srv, ss)
	if err != nil {
		log.Printf("failed to invoke stream RPC: %v\n", err)
	}
	log.Printf("****** after stream auth server interceptor ******")
	return err
}

func streamLogServerInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("====== before stream log server interceptor ======")
	err := handler(srv, ss)
	if err != nil {
		log.Printf("failed to handler stream RPC: %v\n", err)
	}
	log.Printf("====== after stream log server interceptor ======")
	return err
}

func TestStreamInterceptorGrpcServer(t *testing.T) {
	//监听端口
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	creds, err := credentials.NewServerTLSFromFile(CertFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load certificates: %v", err)
	}

	opts := []grpc.ServerOption{
		// 1. TLS Credential
		grpc.Creds(creds),
		// 2.Interceptor
		grpc.ChainStreamInterceptor(streamLogServerInterceptor, streamAuthServerInterceptor),
	}
	s := grpc.NewServer(opts...)

	//将SimpleHello注册到grpc中
	pb.RegisterHelloServer(s, &streamServer{})
	log.Println("Listen on " + Address)

	//开启grpc服务
	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}

}
