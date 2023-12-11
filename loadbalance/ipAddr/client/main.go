package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/loadbalance/ipAddr/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/grpc/resolve"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

	// 使用不安全传输
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 负载均衡策略，轮询方式，文档说明
	// https://github.com/grpc/grpc-go/tree/master/examples/features/load_balancing
	// https://github.com/grpc/grpc/blob/master/doc/service_config.md
	options = append(options, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))

	return options
}

func main() {
	endpoint := resolve.Register("grpc", "hello.grpc.io", []string{"127.0.0.1:8282", "127.0.0.1:8482", "127.0.0.1:8682"})

	roundRobinConn, err := grpc.Dial(endpoint, getDialOptions()...)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = roundRobinConn.Close()
	}()

	client := pb.NewGreeterClient(roundRobinConn)
	for {
		err = sayHello(client)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 2)
	}
}
