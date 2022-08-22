package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/loadbalance/client_loadbalance/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	timestampFormat = "2006-01-02T15:04:05.000000"
)

func sayHello(client pb.GreeterClient) error {
	time.Now().Nanosecond()
	// 创建 metadata 和 context
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat)+"C")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// 使用带有元数据的上下文进行RPC调用。
	var header, trailer metadata.MD
	options := []grpc.CallOption{
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	}

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "foo"}, options...)
	if err != nil {
		return err
	}
	fmt.Println(resp.Message)

	printMetadata(header, "hello")
	printMetadata(header, "foo")
	printMetadata(header, "timestamp")
	printMetadata(trailer, "timestamp")

	mergerMD := metadata.Join(md, header, trailer)
	printMetadata(mergerMD, "timestamp")

	return nil
}

func printMetadata(md metadata.MD, key string) {
	if ts, ok := md[key]; ok {
		fmt.Printf("metadata: %s=%v\n", key, ts)
	} else {
		fmt.Printf("not found '%s' in metadata\n", key)
	}
}

func main() {
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.Dial("127.0.0.1:8080", options...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	err = sayHello(client)
	if err != nil {
		panic(err)
	}
}
