package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zhufuyi/grpc_examples/tracing/rpc2rpc"
	pb "github.com/zhufuyi/grpc_examples/tracing/rpc2rpc/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/tracer"
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

	// tracing跟踪
	options = append(options, grpc.WithUnaryInterceptor(
		interceptor.UnaryClientTracing(),
	))

	return options
}

func main() {
	rpc2rpc.InitTrace("hello-client")
	defer tracer.Close(context.Background()) //nolint

	conn, err := grpc.Dial("127.0.0.1:8482", getDialOptions()...)
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	for i := 0; i < 3; i++ {
		if err := sayHello(client); err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 3)
	}
	time.Sleep(time.Second * 5)
}
