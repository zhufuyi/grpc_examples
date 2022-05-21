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
	"google.golang.org/grpc/credentials/insecure"
)

var helloClient pb.GreeterClient

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

// 一元RPC
func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	parentSpan := opentracing.SpanFromContext(ctx)

	// 模拟redis
	redisQuerySpan := opentracing.StartSpan("redis query", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(10 * time.Millisecond)
	redisQuerySpan.Finish()

	// 模拟mysql查询
	mysqlQuerySpan := opentracing.StartSpan("mysql query", opentracing.ChildOf(parentSpan.Context()))
	time.Sleep(30 * time.Millisecond)
	mysqlQuerySpan.Finish()

	// 调用server2
	uploadFileSpan := opentracing.StartSpan("invoker rpc 2", opentracing.ChildOf(parentSpan.Context()))
	resp, err := helloClient.SayHi(ctx, &pb.HelloRequest{Name: "lisi"})
	if err != nil {
		return nil, err
	}
	uploadFileSpan.Finish()

	return &pb.HelloReply{Message: resp.Message + ", hello " + r.Name}, nil
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
	))

	return options
}

func connectRPCServer(rpcAddr string) {
	conn, err := grpc.Dial(rpcAddr, getDialOptions()...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("connect RPC server(%s) success.\n", rpcAddr)
	helloClient = pb.NewGreeterClient(conn)
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

	// 连接server2
	connectRPCServer("127.0.0.1:8081")

	addr := ":8080"
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
