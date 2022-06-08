package main

import (
	"context"
	"fmt"

	"github.com/zhufuyi/grpc_examples/pkg/tracer"
	pb "github.com/zhufuyi/grpc_examples/tracing/api2rpc/proto/hellopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

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

	// 禁用tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		tracer.UnaryClientTracing(),
	))

	return options
}

func main() {
	// 连接jaeger服务端
	_, err := tracer.InitJaeger("hello_server", "192.168.3.36:6831")
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	if err := sayHello(client); err != nil {
		panic(err)
	}

}
