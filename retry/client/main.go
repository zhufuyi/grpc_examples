package main

import (
	"context"
	"fmt"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	pb "github.com/zhufuyi/grpc_examples/retry/proto/hellopb"
	"github.com/zhufuyi/pkg/grpc/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(client pb.GreeterClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "zhangsan"})
	if err != nil {
		return err
	}
	fmt.Println(resp.Message)

	return nil
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 重试
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			middleware.UnaryClientRetry(), // 可以修改默认重试次数、重试时间间隔、触发重试错误码
		),
	)
	options = append(options, option)

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)

	client := pb.NewGreeterClient(conn)

	err = sayHello(client)
	if err != nil {
		fmt.Println(err)
	}
}
