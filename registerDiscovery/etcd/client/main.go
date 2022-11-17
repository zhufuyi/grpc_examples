package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/registerDiscovery/etcd/proto/hellopb"

	"github.com/zhufuyi/pkg/etcdcli"
	"github.com/zhufuyi/pkg/servicerd/discovery"
	"github.com/zhufuyi/pkg/servicerd/registry/etcd"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverName = "helloDemo"

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

	// 使用etcd服务发现
	options = append(options, discoveryEtcd())

	// 使用不安全传输
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	return options
}

func discoveryEtcd() grpc.DialOption {
	cli, err := etcdcli.Init(etcdAddrs, etcdcli.WithDialTimeout(time.Second*5))
	if err != nil {
		panic(fmt.Sprintf("etcdcli.Init error: %v, addr: %v", err, etcdAddrs))
	}
	iDiscovery := etcd.New(cli)

	return grpc.WithResolvers(
		discovery.NewBuilder(
			iDiscovery,
			discovery.WithInsecure(true),
		))
}

func main() {
	endpoint := "discovery:///" + serverName // 通过服务名称连接grpc服务
	//endpoint := "127.0.0.1:8282"

	conn, err := grpc.Dial(endpoint, getDialOptions()...)
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	if err := sayHello(client); err != nil {
		fmt.Println(err)
	}
}
