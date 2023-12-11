package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	pb "github.com/zhufuyi/grpc_examples/registerDiscovery/nacos/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/servicerd/registry"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"
	"google.golang.org/grpc"
)

const (
	serverName = "helloDemo"
)

var (
	grpcPort  = "8282"
	nacosAddr = "192.168.3.37:8848"
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
	fmt.Println("start grpc server", addr)

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

func registryNacos(rpcAddr string) (registry.Registry, *registry.ServiceInstance) {
	ip, port := parseNacosAddr(nacosAddr)
	localIP, localPort := parseNacosAddr(rpcAddr)
	iRegistry, instance, err := nacos.NewRegistry(
		ip,
		port,
		"public",
		fmt.Sprintf("%s_grpc_%s_%d", serverName, localIP, localPort),
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

	// 注册服务到nacos
	iRegistry, instance := registryNacos(grpcAddr)

	// 运行grpc服务
	go startServer(grpcAddr)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	<-c
	_ = iRegistry.Deregister(context.Background(), instance)
}

func parseNacosAddr(addr string) (string, int) {
	ss := strings.Split(addr, ":")
	if len(ss) != 2 {
		panic(addr + " nacos addr error")
	}
	port, _ := strconv.Atoi(ss[1])
	return ss[0], port
}
