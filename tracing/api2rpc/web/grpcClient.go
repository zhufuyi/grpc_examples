package main

import (
	"fmt"

	"grpc_examples/pkg/tracer/otgrpc"
	pb "grpc_examples/tracing/api2rpc/proto/hellopb"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var helloClient pb.GreeterClient

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
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
