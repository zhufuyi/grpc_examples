package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "grpc_examples/loadbalance/client_loadbalance/proto/hellopb"

	"google.golang.org/grpc"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
	addr string
}

func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	message := fmt.Sprintf("hello %s, (from %s)", r.Name, g.addr)
	return &pb.HelloReply{Message: message}, nil
}

func startServer(addr string) {
	fmt.Println("start rpc server", addr)

	// 监听TCP端口
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &GreeterServer{addr: addr})

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}

func main() {
	addrs := []string{":8080", ":8081", ":8082"}

	var wg sync.WaitGroup
	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			startServer(addr)
		}(addr)
	}
	wg.Wait()
}
