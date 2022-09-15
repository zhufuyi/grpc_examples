package main

import (
	"context"
	"fmt"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/timeout/proto/hellopb"

	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	time.Sleep(time.Second) // 等待时间，让调用者超时错误
	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func main() {
	addr := ":8080"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	pb.RegisterGreeterServer(server, &greeterServer{})

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
