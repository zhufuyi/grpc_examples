package main

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/recovery/proto/hellopb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/zhufuyi/pkg/grpc/interceptor"
	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	var data []int
	fmt.Println(data[5]) // 下标越界，触发panic
	return &pb.HelloReply{Message: time.Now().Format("2006-01-02T15:04:05.000000") + " hello " + r.Name}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	recoveryOption := grpc_middleware.WithUnaryServerChain(
		interceptor.UnaryServerRecovery(),
	)
	options = append(options, recoveryOption)

	return options
}

func main() {
	addr := ":8282"
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
