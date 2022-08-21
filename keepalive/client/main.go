package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/keepalive/proto/hellopb"
	"github.com/zhufuyi/pkg/grpc/keepalive"
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

	// 禁止tls加密
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// keepalive option
	options = append(options, keepalive.ClientKeepAlive())

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)
	err = sayHello(client)
	if err != nil {
		panic(err)
	}

	// 保持连接查看连接状态
	for {
		fmt.Println("connectStat:", conn.GetState())
		time.Sleep(time.Second * 5)
	}
}
