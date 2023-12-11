package main

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/zhufuyi/grpc_examples/loadbalance/ipAddr/proto/hellopb"

	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
	addr string
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	message := fmt.Sprintf("hello %s, (from %s)", r.Name, g.addr)
	return &pb.HelloReply{Message: message}, nil
}

func startServer(addr string) {
	fmt.Println("grpc service is running", addr)

	// listening on TCP port
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	pb.RegisterGreeterServer(server, &greeterServer{addr: addr})

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}

func main() {
	addrs := []string{":8282", ":8482", ":8682"}

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
