package grpcgateway

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"golang.org/x/net/context"
	"golang.org/x/net/http2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "go-sample/net/grpc/http/grpc_gateway/proto"
)

func TestGrpcGatewayHttpClient(t *testing.T) {
	startH2Client()
}

func startGrpcClient() {
	creds, err := credentials.NewClientTLSFromFile(CertFile, "localhost")
	if err != nil {
		log.Printf("Failed to create TLS credentials %v", err)
	}
	conn, err1 := grpc.Dial(GrpcAddress, grpc.WithTransportCredentials(creds))
	if err1 != nil {
		log.Println(err1)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)
	body := &pb.HelloRequest{Name: "Grpc client test"}

	res, err2 := c.GetHello(context.Background(), body)
	if err2 != nil {
		log.Println(err2)
	}

	log.Printf("Client receive：%v", res.Message)
}

// h2 client
func startH2Client() {
	crt, err := ioutil.ReadFile(CertFile)
	if err != nil {
		log.Fatal(err)
	}

	rootCAs := x509.NewCertPool()
	rootCAs.AppendCertsFromPEM(crt)

	client := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            rootCAs,
				InsecureSkipVerify: false,
				ServerName:         "localhost",
			},
			DisableCompression: true,
			AllowHTTP:          false,
		},
	}

	// Get请求
	resp, err := client.Get("https://localhost:8080/get/test")
	// Post请求
	//resp, err := client.Post("https://localhost:8080/post",
	//	"application/x-www-form-urlencoded", strings.NewReader("name=test"))
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	resp.Body.Close()

	log.Printf("Response Http Proto: %s Content: %s", resp.Proto, string(bytes))
}
