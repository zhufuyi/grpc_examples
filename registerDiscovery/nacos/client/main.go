package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	pb "github.com/zhufuyi/grpc_examples/registerDiscovery/nacos/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/nacoscli"
	"github.com/zhufuyi/sponge/pkg/servicerd/discovery"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/nacos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serverName = "helloDemo"

var nacosAddr = "192.168.3.37:8848"

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

	// 使用nacos服务发现
	options = append(options, discoveryFromNacos())

	// 使用不安全传输
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 负载均衡策略，轮询方式，文档说明
	// https://github.com/grpc/grpc-go/tree/master/examples/features/load_balancing
	// https://github.com/grpc/grpc/blob/master/doc/service_config.md
	options = append(options, grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))

	return options
}

func discoveryFromNacos() grpc.DialOption {
	ip, port := parseNacosAddr(nacosAddr)
	cli, err := nacoscli.NewNamingClient(ip, port, "public")
	if err != nil {
		panic(fmt.Sprintf("nacoscli.Init error: %v, addr: %v", err, nacosAddr))
	}
	iDiscovery := nacos.New(cli)

	return grpc.WithResolvers(
		discovery.NewBuilder(
			iDiscovery,
			discovery.WithInsecure(true),
		))
}

func parseNacosAddr(addr string) (string, int) {
	ss := strings.Split(addr, ":")
	if len(ss) != 2 {
		panic(addr + " nacos addr error")
	}
	port, _ := strconv.Atoi(ss[1])
	return ss[0], port
}

func main() {
	endpoint := "discovery:///" + serverName + ".grpc" // 通过服务名称连接grpc服务
	//endpoint := "127.0.0.1:8282"

	conn, err := grpc.Dial(endpoint, getDialOptions()...)
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	for {
		if err = sayHello(client); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 10)
	}
}
