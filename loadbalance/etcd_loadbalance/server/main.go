package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/zhufuyi/grpc_examples/loadbalance/etcd_loadbalance/proto/hellopb"

	"github.com/zhufuyi/pkg/servicerd/registry"
	"github.com/zhufuyi/pkg/servicerd/registry/etcd"
	"google.golang.org/grpc"
)

const (
	serverName = "hello-demo"
)

var (
	etcdAddrs = []string{"192.168.3.37:2379"}
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
	fmt.Println("start rpc server", addr)

	// 监听TCP端口
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

func registryEtcd(rpcAddr string) (registry.Registry, *registry.ServiceInstance) {
	iRegistry, instance, err := etcd.NewRegistry(
		etcdAddrs,
		serverName+"_grpc_"+rpcAddr,
		serverName,
		[]string{rpcAddr},
	)
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second) //nolint
	if err = iRegistry.Register(ctx, instance); err != nil {
		panic(err)
	}

	return iRegistry, instance
}

func main() {
	grpcAddrs := []string{"grpc://127.0.0.1:8282", "grpc://127.0.0.1:8482", "grpc://127.0.0.1:8682"}

	for i := 0; i < len(grpcAddrs); i++ {
		// 启动rpc服务
		go startServer(grpcAddrs[i])

		iRegistry, instance := registryEtcd(grpcAddrs[i])
		defer func() {
			_ = iRegistry.Deregister(context.Background(), instance)
		}()
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
}
