package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/zhufuyi/grpc_examples/tracing/rpc2rpc"
	pb "github.com/zhufuyi/grpc_examples/tracing/rpc2rpc/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/tracer"
	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHi 一元RPC
func (g *greeterServer) SayHi(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	time.Sleep(time.Millisecond * 15)
	resp := &pb.HelloReply{Message: "hi " + r.Name}
	fmt.Printf("resp: %s\n", resp.Message)

	rpc2rpc.SpanDemo(ctx, "sayHi") // 模拟创建一个span

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
	rpc2rpc.InitTrace("hello-server2")
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
