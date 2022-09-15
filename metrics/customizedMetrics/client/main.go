package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/zhufuyi/grpc_examples/metrics/customizedMetrics/proto/hellopb"

	"github.com/zhufuyi/pkg/grpc/interceptor"
	"github.com/zhufuyi/pkg/grpc/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(client pb.GreeterClient, i int) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "foo " + strconv.Itoa(i)})
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
	options = append(options, grpc.WithUnaryInterceptor(interceptor.UnaryClientMetrics()))
	options = append(options, grpc.WithStreamInterceptor(interceptor.StreamClientMetrics()))
	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)
	if err != nil {
		panic(err)
	}

	metrics.ClientHTTPService(":9094")
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
