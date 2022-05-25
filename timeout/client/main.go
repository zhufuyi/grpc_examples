package main

import (
	"context"
	"fmt"
	"time"

	pb "grpc_examples/timeout/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(client pb.GreeterClient) error {
	to := time.Millisecond * 200
	ctx, _ := context.WithTimeout(context.Background(), to)

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "zhangsan"})
	if err != nil {
		return fmt.Errorf("%v, timeout=%v", err, to)
	}

	fmt.Println(resp.Message)
	return nil
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

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
