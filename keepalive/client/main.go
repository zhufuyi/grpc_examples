package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/keepalive/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ClientKeepAlive 保持连接设置
func ClientKeepAlive() grpc.DialOption {
	return grpc.WithKeepaliveParams(
		keepalive.ClientParameters{
			Time:                20 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout:             1 * time.Second,  // wait 1 second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		},
	)
}

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

	// keepalive option
	options = append(options, ClientKeepAlive())

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8282", getDialOptions()...)
	if err != nil {
		panic(err)
	}
	defer conn.Close() //nolint

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
