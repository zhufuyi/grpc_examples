package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/zhufuyi/grpc_examples/tracing"
	pb "github.com/zhufuyi/grpc_examples/tracing/http2rpc/proto/hellopb"

	"github.com/zhufuyi/pkg/grpc/interceptor"
	"github.com/zhufuyi/pkg/tracer"
	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello 一元RPC
func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	time.Sleep(time.Millisecond * 15)
	resp := &pb.HelloReply{Message: "hello " + r.Name}
	fmt.Printf("resp: %s\n", resp.Message)

	tracing.SpanDemo(ctx, "sayHello") // 模拟创建一个span

	return resp, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// 链路跟踪拦截器
	options = append(options, grpc.UnaryInterceptor(
		interceptor.UnaryServerTracing(),
	))

	return options
}

func main() {
	tracing.InitTrace("hello-server")
	defer tracer.Close(context.Background()) //nolint

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
	pb.RegisterGreeterServer(server, &greeterServer{})

	// 调用服务器执行阻塞等待客户端
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
