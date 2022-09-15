package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/interceptor/proto/hellopb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(client pb.GreeterClient, name string) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		return err
	}

	fmt.Println(resp.Message)
	return nil
}

// Timeout 超时拦截器
func Timeout(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	}
	if cancel != nil {
		defer cancel()
	}

	return invoker(ctx, method, req, resp, cc, opts...)
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁止tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 超时拦截器
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			Timeout,
		),
	)
	options = append(options, option)

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)
	if err != nil {
		panic(err)
	}
	defer conn.Close() //nolint

	client := pb.NewGreeterClient(conn)

	err = sayHello(client, "foo")
	if err != nil {
		fmt.Println(err)
	}

	err = sayHello(client, "")
	if err != nil {
		fmt.Println(err)
	}
}
