package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "grpc_examples/metrics/customizedMetrics/proto/hellopb"
	"grpc_examples/pkg/metrics/clientMetrics"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(client pb.GreeterClient, i int) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhangsan " + strconv.Itoa(i)})
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

	// Metrics
	options = append(options, grpc.WithUnaryInterceptor(clientMetrics.UnaryClientMetrics()))
	options = append(options, grpc.WithStreamInterceptor(clientMetrics.StreamClientMetrics()))
	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)

	clientMetrics.Serve(":9094")
	fmt.Println("start metrics server", ":9094")

	client := pb.NewGreeterClient(conn)
	i := 0
	for {
		i++
		time.Sleep(time.Millisecond * 500) // qps is 2
		err = sayHello(client, i)
		if err != nil {
			fmt.Println(err)
		}
	}
}
