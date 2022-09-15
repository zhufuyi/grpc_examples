package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/zhufuyi/grpc_examples/registerDiscovery/etcd/proto/hellopb"

	"github.com/zhufuyi/pkg/registry"
	"github.com/zhufuyi/pkg/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serverName = "hello-demo"
)

var (
	grpcPort  = "8080"
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

func getETCDRegistry(etcdEndpoints []string, instanceName string, instanceEndpoints []string) (registry.Registry, *registry.ServiceInstance) {
	serviceInstance := registry.NewServiceInstance(instanceName, instanceEndpoints)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdEndpoints,
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	})
	if err != nil {
		panic(err)
	}
	iRegistry := etcd.New(cli)

	return iRegistry, serviceInstance
}

func main() {
	flag.StringVar(&grpcPort, "p", "8080", "grpc listen port")
	flag.Parse()
	grpcAddr := "127.0.0.1:" + grpcPort

	// 运行rpc服务
	go startServer(grpcAddr)

	// 注册服务到etcd
	iRegistry, serviceInstance := getETCDRegistry(etcdAddrs, serverName, []string{"grpc://" + grpcAddr})
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second) //nolint
	if err := iRegistry.Register(ctx, serviceInstance); err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
	ctx, _ = context.WithTimeout(context.Background(), 3*time.Second) //nolint
	_ = iRegistry.Deregister(ctx, serviceInstance)
}
