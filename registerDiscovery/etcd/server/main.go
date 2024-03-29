package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	pb "github.com/zhufuyi/grpc_examples/registerDiscovery/etcd/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"
	"google.golang.org/grpc"
)

const (
	serverName = "helloDemo"
)

var (
	grpcPort  = "8282"
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

func registryEtcd(rpcAddr string) (registry.Registry, *registry.ServiceInstance) {
	iRegistry, instance, err := etcd.NewRegistry(
		etcdAddrs,
		serverName+"_grpc"+"_"+parseAddr(rpcAddr),
		serverName,
		[]string{"grpc://" + rpcAddr},
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
	var ip string
	flag.StringVar(&grpcPort, "p", "8282", "grpc listen port")
	flag.StringVar(&ip, "i", "127.0.0.1", "host ip")
	flag.Parse()
	grpcAddr := ip + ":" + grpcPort

	// 注册服务到etcd
	iRegistry, instance := registryEtcd(grpcAddr)

	// 运行grpc服务
	go startServer(grpcAddr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
	_ = iRegistry.Deregister(context.Background(), instance)
}

func parseAddr(addr string) string {
	ss := strings.Split(addr, ":")
	if len(ss) == 1 {
		return ss[0] + "_unknown"
	}
	return ss[0] + "_" + ss[1]
}
