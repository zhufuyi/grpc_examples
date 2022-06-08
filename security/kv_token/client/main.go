package main

import (
	"context"
	"fmt"

	"github.com/zhufuyi/grpc_examples/pkg/gtls"
	"github.com/zhufuyi/grpc_examples/pkg/gtls/certfile"
	pb "github.com/zhufuyi/grpc_examples/security/kv_token/proto/hellopb"
	"google.golang.org/grpc"
)

type MyToken struct {
	AppID  string
	AppKey string
}

func (m *MyToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_id":  m.AppID,
		"app_key": m.AppKey,
	}, nil
}

func (m *MyToken) RequireTransportSecurity() bool {
	return true
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// tls加密
	credentials, err := gtls.GetClientTLSCredentials("localhost", certfile.Path("one-way/server.crt"))
	if err != nil {
		panic(err)
	}
	options = append(options, grpc.WithTransportCredentials(credentials))

	// token
	options = append(options, grpc.WithPerRPCCredentials(&MyToken{"grpc", "123456"}))

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhangsan"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Message)
}
