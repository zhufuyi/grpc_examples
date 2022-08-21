package main

import (
	"context"
	"fmt"
	"net"

	pb "github.com/zhufuyi/grpc_examples/keepalive/proto/hellopb"
	"github.com/zhufuyi/pkg/grpc/keepalive"
	"google.golang.org/grpc"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func main() {
	addr := ":8080"
	fmt.Println("start rpc server", addr)

	var options []grpc.ServerOption
	// 默认是每15秒向客户端发送一次ping，修改为间隔20秒发送一次ping
	options = append(options, keepalive.ServerKeepAlive()...)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(options...)

	pb.RegisterGreeterServer(server, &GreeterServer{})

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
