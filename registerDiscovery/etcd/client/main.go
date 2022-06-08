package main

import (
	"context"
	"fmt"

	"github.com/zhufuyi/grpc_examples/pkg/etcd/discovery"
	pb "github.com/zhufuyi/grpc_examples/registerDiscovery/etcd/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const serverName = "hello-demo"

var etcdAddrs = []string{"192.168.3.36:2379"}

func sayHello(client pb.GreeterClient) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhangsan"})
	if err != nil {
		return err
	}

	fmt.Println(resp.Message)
	return nil
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁止tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return options
}

func main() {
	// 使用etcd服务发现
	r := discovery.NewResolver(etcdAddrs)
	//r := discovery.NewResolver( // 自定义设置方式
	//	etcdAddrs,
	//	discovery.WithDialTimeout(5),
	//	discovery.WithLogger(zap.NewExample()),
	//)
	resolver.Register(r)

	conn, err := grpc.Dial("etcd:///"+serverName, getDialOptions()...)
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	if err := sayHello(client); err != nil {
		fmt.Println(err)
	}
}
