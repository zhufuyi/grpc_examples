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
	"google.golang.org/grpc/credentials/insecure"
)

var helloClient pb.GreeterClient

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello 一元RPC
func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	time.Sleep(time.Millisecond * 15)
	resp, err := helloClient.SayHi(ctx, &pb.HelloRequest{Name: "foo"})
	if err != nil {
		return nil, err
	}

	resp2 := &pb.HelloReply{Message: resp.Message + ", hello " + r.Name}
	fmt.Printf("resp: %s\n", resp2.Message)

	rpc2rpc.SpanDemo(ctx, "sayHello") // 模拟创建一个span

	return resp2, nil
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 使用不安全传输
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		interceptor.UnaryClientTracing(),
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
		interceptor.UnaryServerTracing(),
	))

	return options
}

func main() {
	rpc2rpc.InitTrace("hello-server1")
	defer tracer.Close(context.Background()) //nolint

	// 连接server2
	connectRPCServer("127.0.0.1:8282")

	addr := ":8482"
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
