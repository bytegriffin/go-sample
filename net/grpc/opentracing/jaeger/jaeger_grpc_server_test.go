package jaeger

import (
	"context"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc/reflection"

	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	pb "go-sample/net/grpc/opentracing/proto"

	opentracing "github.com/opentracing/opentracing-go"
)

const (
	ServerAddress     = "127.0.0.1:50051"
	jaegerServiceName = "test"
	jagentHost        = "127.0.0.1:6831"
)

type unaryServer struct {
	pb.UnimplementedHelloServer
}

func (s *unaryServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Server received：%v", in.GetName())

	//开启一个Span用来追踪SQL
	//if parent := opentracing.SpanFromContext(ctx); parent != nil {
	//	pctx := parent.Context()
	//	if tracer := opentracing.GlobalTracer(); tracer != nil {
	//		mysqlSpan := tracer.StartSpan("FindUserTable", opentracing.ChildOf(pctx))
	//		//模拟mysql操作
	//		time.Sleep(time.Millisecond * 100)
	//		defer mysqlSpan.Finish()
	//	}
	//}

	return &pb.HelloResponse{Code: 200, Message: "hello，" + in.GetName()}, nil
}

func ServerInterceptor(tracer opentracing.Tracer) grpc.UnaryServerInterceptor {
	log.Println("==============Server Interceptor==============")
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp interface{}, err error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		//SpanContext分为三种传输方式：Binary、TextMap。HTTPHeaders
		spanContext, err := tracer.Extract(opentracing.TextMap, MDCarrier{md})
		if err != nil && err != opentracing.ErrSpanContextNotFound {
			log.Printf("extract from metadata err: %v", err)
		} else {
			span := tracer.StartSpan(
				info.FullMethod,
				ext.RPCServerOption(spanContext),
				opentracing.Tag{Key: string(ext.Component), Value: "gRPC Server"},
				ext.SpanKindRPCServer,
			)
			defer span.Finish()

			ctx = opentracing.ContextWithSpan(ctx, span)
		}

		return handler(ctx, req)
	}
}

/**
  OpenTracing是一个分布式追踪协议，其中Traces（调用链）由一组Span（执行过程中个记录的信息）
  构成的有向无环图（DAG），Span表示Jaeger中的逻辑工作单元，具有操作名称、操作的开始时间和持续时间，
  Span之间被称为References。
  Jaeger是由Uber Technologies开发的开源发布的分布式跟踪系统，是OpenTracing API的特定于语言的实现。

  https://opentracing.io/specification/
  https://www.jaegertracing.io/docs/1.21/getting-started/
*/
func TestJaegerGrpcServer(t *testing.T) {
	listen, err := net.Listen("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var servOpts []grpc.ServerOption
	tracer, closer, err := NewJaegerTracer(jaegerServiceName, jagentHost)
	defer closer.Close()
	if err != nil {
		log.Printf("new tracer err: %+v\n", err)
	}
	if tracer != nil {
		servOpts = append(servOpts, grpc.UnaryInterceptor(ServerInterceptor(tracer)))
	}

	s := grpc.NewServer(servOpts...)

	pb.RegisterHelloServer(s, &unaryServer{})

	reflection.Register(s)
	log.Println("Listen on " + ServerAddress)

	if err := s.Serve(listen); err != nil {
		log.Fatalln("Grpc Server is failed.")
	}
}
