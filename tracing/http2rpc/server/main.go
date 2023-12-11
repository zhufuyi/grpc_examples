package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/zhufuyi/grpc_examples/tracing/http2rpc"
	pb "github.com/zhufuyi/grpc_examples/tracing/http2rpc/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/tracer"
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

	addr := ":8282"
	fmt.Println("grpc service is running", addr)

	// listening on TCP port
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// create a grpc server object where interceptors can be injected
	server := grpc.NewServer(getServerOptions()...)

	// register greeterServer to the server
	pb.RegisterGreeterServer(server, &greeterServer{})

	// start the server
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
