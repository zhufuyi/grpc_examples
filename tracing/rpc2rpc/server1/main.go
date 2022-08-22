package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/zhufuyi/grpc_examples/tracing"
	pb "github.com/zhufuyi/grpc_examples/tracing/http2rpc/proto/hellopb"
	"github.com/zhufuyi/pkg/grpc/middleware"
	"github.com/zhufuyi/pkg/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var helloClient pb.GreeterClient

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello 一元RPC
func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	time.Sleep(time.Millisecond * 15)
	resp, err := helloClient.SayHi(ctx, &pb.HelloRequest{Name: "foo"})
	if err != nil {
		return nil, err
	}

	resp2 := &pb.HelloReply{Message: resp.Message + ", hello " + r.Name}
	fmt.Printf("resp: %s\n", resp2.Message)

	tracing.SpanDemo("sayHello", ctx) // 模拟创建一个span

	return resp2, nil
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		middleware.UnaryClientTracing(),
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
		middleware.UnaryServerTracing(),
	))

	return options
}

func main() {
	tracing.InitTrace("hello-server1")
	defer tracer.Close(context.Background())

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
