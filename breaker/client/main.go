package main

import (
	"context"
	"fmt"

	pb "github.com/zhufuyi/grpc_examples/breaker/proto/hellopb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/zhufuyi/pkg/grpc/interceptor"
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

	// use insecure transfer
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// circuit breaker
	option := grpc.WithUnaryInterceptor(
		grpc_middleware.ChainUnaryClient(
			interceptor.UnaryClientCircuitBreaker(),
		),
	)
	options = append(options, option)

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8282", getDialOptions()...)
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	if err := sayHello(client); err != nil {
		fmt.Println(err)
	}
}
