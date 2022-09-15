package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/loadbalance/etcd_loadbalance/proto/hellopb"

	"github.com/zhufuyi/pkg/grpc/etcd/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

const serverName = "hello-demo"

var etcdAddrs = []string{"192.168.3.37:2379"}

func sayHello(client pb.GreeterClient) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "foo"})
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

	// 负载均衡策略，轮询，https://github.com/grpc/grpc-go/tree/master/examples/features/load_balancing 和 https://github.com/grpc/grpc/blob/master/doc/service_config.md
	options = append(options, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))

	return options
}

func main() {
	// 使用etcd服务发现
	r := discovery.NewResolver(etcdAddrs)
	resolver.Register(r)

	conn, err := grpc.Dial("etcd:///"+serverName, getDialOptions()...)
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)
	for {
		err = sayHello(client)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 2)
	}
}
