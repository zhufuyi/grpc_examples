package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/hystrix/withMetrics/proto/hellopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	n := rand.Intn(100)
	if n%3 == 0 { // 大约30%概率出错
		time.Sleep(time.Millisecond * 3)
		return nil, status.Errorf(codes.DeadlineExceeded, "name: %s", r.Name)
	}
	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	addr := ":8080"
	fmt.Println("start rpc server", addr)

	// 监听TCP端口
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// 创建grpc server对象，拦截器可以在这里注入
	server := grpc.NewServer()

	// grpc的server内部服务和路由
	pb.RegisterGreeterServer(server, &GreeterServer{})

	// 调用服务器执行阻塞等待客户端
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
