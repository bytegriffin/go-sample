package cancellation

import (
	"io"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/cancellation/proto"
)

const (
	ServerAddress = "127.0.0.1:50191"
)

type server struct {
	pb.UnimplementedHelloServiceServer
}

func (s *server) SaveAll(stream pb.HelloService_SaveAllServer) error {
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Printf("server: error receiving from stream: %v\n", err)
			if err == io.EOF {
				return nil
			}
			return err
		}
		log.Printf("server receive message %q\n", in.GetName())
		stream.Send(&pb.HelloResponse{Message: in.GetName() + " has received."})
	}
}

func TestCancelGrpcServer(t *testing.T) {

	lis, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("server listening at port %v\n", lis.Addr())
	s := grpc.NewServer()
	pb.RegisterHelloServiceServer(s, &server{})
	s.Serve(lis)

}
