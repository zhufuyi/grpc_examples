package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"grpc_examples/pkg/etcd/discovery"
	pb "grpc_examples/registerDiscover/etcd/proto/hellopb"

	"google.golang.org/grpc"
)

const (
	serverName = "hello-demo"
)

var (
	grpcPort  = "8080"
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
	flag.StringVar(&grpcPort, "p", "8080", "grpc listen port")
	flag.Parse()
	grpcAddr := "127.0.0.1:" + grpcPort

	// 运行rpc服务
	go startServer(grpcAddr)

	// 注册服务到etcd
	etcdRegister := discovery.RegisterRPCAddr(serverName, grpcAddr, etcdAddrs)
	//etcdRegister := discovery.RegisterRPCAddr( // 自定义设置方式
	//	serverName,
	//	grpcAddr,
	//	etcdAddrs,
	//	discovery.WithTTLSeconds(5),
	//	discovery.WithLogger(zap.NewExample()),
	//	discovery.WithWeight(5),
	//)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
	etcdRegister.Stop()
}
