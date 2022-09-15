package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/registerDiscovery/etcd/proto/hellopb"

	"github.com/zhufuyi/pkg/grpc/grpccli"
	"github.com/zhufuyi/pkg/logger"
	"github.com/zhufuyi/pkg/registry"
	"github.com/zhufuyi/pkg/registry/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

func discoveryETCD(endpoints []string) registry.Discovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 10 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		},
	})
	if err != nil {
		panic(err)
	}

	return etcd.New(cli)
}

func main() {
	endpoint := "discovery:///" + serverName
	iDiscovery := discoveryETCD(etcdAddrs)
	ctx, _ := context.WithTimeout(context.Background(), time.Second*3) //nolint
	conn, err := grpccli.DialInsecure(ctx, endpoint,
		grpccli.WithEnableLog(logger.Get()),
		grpccli.WithDiscovery(iDiscovery),
	)
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	if err := sayHello(client); err != nil {
		fmt.Println(err)
	}
}
