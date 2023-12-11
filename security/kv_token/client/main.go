package main

import (
	"context"
	"fmt"

	pb "github.com/zhufuyi/grpc_examples/security/kv_token/proto/hellopb"

	"github.com/zhufuyi/sponge/pkg/grpc/gtls"
	"github.com/zhufuyi/sponge/pkg/grpc/gtls/certfile"
	"google.golang.org/grpc"
)

type myToken struct {
	AppID  string
	AppKey string
}

// GetRequestMetadata 获取请求元数据
func (m *myToken) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"app_id":  m.AppID,
		"app_key": m.AppKey,
	}, nil
}

// RequireTransportSecurity 是否安全传输
func (m *myToken) RequireTransportSecurity() bool {
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
	options = append(options, grpc.WithPerRPCCredentials(&myToken{"grpc", "123456"}))

	return options
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8282", getDialOptions()...)
	if err != nil {
		panic(err)
	}
	defer conn.Close() //nolint

	client := pb.NewGreeterClient(conn)

	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "foo"})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Message)
}
