package main

import (
	"fmt"

	"github.com/zhufuyi/pkg/grpc/interceptor"

	pb "github.com/zhufuyi/grpc_examples/tracing/http2rpc/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var helloClient pb.GreeterClient

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 使用不安全传输
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		interceptor.UnaryClientTracing(),
	))

	return options
}

func connectRPCServer(rpcAddr string) {
	conn, err := grpc.Dial(rpcAddr, getDialOptions()...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("connect RPC server(%s) success.\n", rpcAddr)
	helloClient = pb.NewGreeterClient(conn)
}
