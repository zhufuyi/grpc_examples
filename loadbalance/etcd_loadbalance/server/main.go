package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/zhufuyi/grpc_examples/loadbalance/etcd_loadbalance/proto/hellopb"
	"github.com/zhufuyi/pkg/grpc/etcd/discovery"
	"google.golang.org/grpc"
)

const (
	serverName = "hello-demo"
)

var (
	etcdAddrs = []string{"192.168.3.36:2379"}
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
	grpcAddrs := []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"}

	for i := 0; i < len(grpcAddrs); i++ {
		// 启动rpc服务
		go startServer(grpcAddrs[i])

		// 注册服务到etcd
		etcdRegister := discovery.RegisterRPCAddr(serverName, grpcAddrs[i], etcdAddrs)
		defer etcdRegister.Stop()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
}
