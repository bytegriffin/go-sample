package jaeger

import (
	"context"
	"log"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"

	olog "github.com/opentracing/opentracing-go/log"

	pb "go-sample/net/grpc/opentracing/proto"
)

func ClientInterceptor(tracer opentracing.Tracer) grpc.UnaryClientInterceptor {
	log.Println("==============Client Interceptor==============")
	return func(ctx context.Context, method string,
		req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {

		// 从Context获取SpanContext，如果上层没有开启追踪，则新建一个追踪
		// 如果上层已经有了，则创建一个子Span。即：一个RPC调用的服务端的span，
		// 和RPC服务客户端的span构成ChildOf关系
		var parentCtx opentracing.SpanContext
		parentSpan := opentracing.SpanFromContext(ctx)
		if parentSpan != nil {
			parentCtx = parentSpan.Context()
		}
		span := tracer.StartSpan(
			method,
			opentracing.ChildOf(parentCtx),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC Client"},
			ext.SpanKindRPCClient,
		)
		defer span.Finish()

		//从context中取出metadata数据，如果context中没有就创建一个新的metadata。
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		} else {
			md = md.Copy()
		}

		//将追踪数据注入到metadata中
		err := tracer.Inject(span.Context(), opentracing.TextMap, MDCarrier{md})
		if err != nil {
			span.LogFields(olog.String("inject-error", err.Error()))
		}
		//将metadata装入context中
		newCtx := metadata.NewOutgoingContext(ctx, md)
		//使用带有追踪数据的context进行gRPC调用
		err = invoker(newCtx, method, req, reply, cc, opts...)
		if err != nil {
			span.LogFields(olog.String("call-error", err.Error()))
		}
		return err
	}
}

func TestJaegerGrpcClient(t *testing.T) {
	dialOpts := []grpc.DialOption{grpc.WithInsecure()}
	tracer, _, err := NewJaegerTracer(jaegerServiceName, jagentHost)
	if err != nil {
		log.Printf("new tracer err: %+v\n", err)
	}
	if tracer != nil {
		dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(ClientInterceptor(tracer)))
	}

	conn, err := grpc.Dial(ServerAddress, dialOpts...)
	if err != nil {
		log.Fatalf("did not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)
	req := &pb.HelloRequest{Name: "grpc-jaeger"}
	res, err := c.SayHello(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Client receive：%v", res.Message)

}
