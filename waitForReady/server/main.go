package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/waitForReady/proto/hellopb"
	"google.golang.org/grpc"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	// 随机延时0~1000毫秒
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

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

	pb.RegisterGreeterServer(server, &GreeterServer{})

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
