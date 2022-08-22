package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/zhufuyi/grpc_examples/security/tls/proto/hellopb"
	"github.com/zhufuyi/pkg/grpc/gtls"
	"github.com/zhufuyi/pkg/grpc/gtls/certfile"
	"google.golang.org/grpc"
)

func sayHello(client pb.GreeterClient) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Second))
	defer cancel()

	resp, err := client.SayHello(ctx, &pb.HelloRequest{Name: "foo"})
	if err != nil {
		return err
	}

	fmt.Println(resp.Message)
	return nil
}

func main() {
	// 单向认证
	//credentials, err := gtls.GetClientTLSCredentials("localhost", certfile.Path("/one-way/server.crt"))
	// 双向认证
	credentials, err := gtls.GetClientTLSCredentialsByCA(
		"localhost",
		certfile.Path("two-way/ca.pem"),
		certfile.Path("two-way/client/client.pem"),
		certfile.Path("two-way/client/client.key"),
	)
	if err != nil {
		panic(err)
	}

	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(credentials))
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
