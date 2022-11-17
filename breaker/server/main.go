package main

import (
	"context"
	"fmt"
	"net"
	"reflect"

	pb "github.com/zhufuyi/grpc_examples/breaker/proto/hellopb"

	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	t := reflect.TypeOf(*r)

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fmt.Println(sf.Tag.Get("json"), sf.Tag.Get("gorm"), sf.Tag.Get("bson"))
	}

	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

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
