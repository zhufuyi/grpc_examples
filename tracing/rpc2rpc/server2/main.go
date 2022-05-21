package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"grpc_examples/pkg/tracer"
	"grpc_examples/pkg/tracer/otgrpc"
	pb "grpc_examples/tracing/api2rpc/proto/hellopb"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

// 一元RPC
func (g *GreeterServer) SayHi(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	parentSpan := opentracing.SpanFromContext(ctx)

	// 模拟kafka推送
	kafkaPushSpan := opentracing.StartSpan("kafka push", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(10 * time.Millisecond)
	kafkaPushSpan.Finish()

	// 模拟mongodb查询
	mongoQuerySpan := opentracing.StartSpan("mongodb query", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(30 * time.Millisecond)
	mongoQuerySpan.Finish()

	return &pb.HelloReply{Message: "hi " + r.Name}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 链路跟踪拦截器
	options = append(options, grpc.UnaryInterceptor(
		otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer()),
	))

	return options
}

func main() {
	// 连接jaeger服务端
	closer, err := tracer.InitJaeger("tracing_demo", "192.168.3.36:6831")
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	addr := ":8081"
	fmt.Println("start rpc server", addr)

	// 监听TCP端口
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer(getServerOptions()...)

	// grpc的server内部服务和路由
	pb.RegisterGreeterServer(server, &GreeterServer{})

	// 调用服务器执行阻塞等待客户端
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
