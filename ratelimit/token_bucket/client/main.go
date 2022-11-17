package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/zhufuyi/grpc_examples/ratelimit/token_bucket/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(i int, client pb.GreeterClient) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "foo " + strconv.Itoa(i)})
	if err != nil {
		return err
	}

	fmt.Println(resp.Message)
	return nil
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8282", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close() //nolint

	client := pb.NewGreeterClient(conn)

	for i := 0; i < 6; i++ {
		go func(i int) {
			if err := sayHello(i, client); err != nil {
				fmt.Println(i, err)
			}
		}(i)
	}

	time.Sleep(time.Second * 2)
}
