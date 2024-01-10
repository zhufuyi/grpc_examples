package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/usage/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/etcdcli"
	"github.com/zhufuyi/sponge/pkg/grpc/client"
	"github.com/zhufuyi/sponge/pkg/grpc/gtls"
	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/servicerd/discovery"
	"github.com/zhufuyi/sponge/pkg/servicerd/registry/etcd"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
)

var newUnaryInterceptors = []grpc.UnaryClientInterceptor{
	interceptor.UnaryClientRecovery(),
	interceptor.UnaryClientLog(logger.Get()),
}

func newBuilder() resolver.Builder {
	cli, err := etcdcli.Init([]string{"192.168.3.37:2379"}, etcdcli.WithDialTimeout(time.Second*5))
	if err != nil {
		panic(fmt.Sprintf("etcdcli.Init error: %v", err))
	}
	iDiscovery := etcd.New(cli)

	return discovery.NewBuilder(
		iDiscovery,
		discovery.WithInsecure(true),
	)
}

func newCredential() credentials.TransportCredentials {
	var (
		credential credentials.TransportCredentials
		err        error
	)

	// server-side authentication only
	credential, err = gtls.GetClientTLSCredentials(
		"localhost",
		certfile.Path("one-way/server.crt"),
	)

	// two-way authentication between client and server
	//credential, err = gtls.GetClientTLSCredentialsByCA(
	//	"localhost",
	//	certfile.Path("two-way/ca.pem"),
	//	certfile.Path("two-way/client/client.pem"),
	//	certfile.Path("two-way/client/client.key"),
	//)
	if err != nil {
		panic(err)
	}
	return credential
}

func main() {
	endpoint := "127.0.0.1:8282" // ip direct connection
	//endpoint = "discovery:///hello-demo" // used when using service discovery

	conn, err := client.Dial(context.Background(), endpoint,
		//client.WithServiceDiscover(newBuilder()),
		client.WithLoadBalance(),
		client.WithSecure(newCredential()),
		client.WithUnaryInterceptor(newUnaryInterceptors...),
		//WithStreamInterceptor(streamInterceptors...),
	)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	greeterClient := pb.NewGreeterClient(conn)

	for i := 0; i < 100; i++ {
		time.Sleep(time.Second * 3)
		reply, err := greeterClient.SayHello(context.Background(), &pb.HelloRequest{
			Name: fmt.Sprintf("Tom-%d", i),
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(reply.Message)
	}
}
