package main

import (
	"context"
	"fmt"

	"github.com/zhufuyi/grpc_examples/pkg/gtls"
	"github.com/zhufuyi/grpc_examples/pkg/gtls/certfile"
	pb "github.com/zhufuyi/grpc_examples/security/jwt_token/proto/accountpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var isUseTLS bool // 是否开启TLS加密

func registerUser(client pb.AccountClient) (*pb.RegisterReply, error) {
	resp, err := client.Register(context.Background(), &pb.RegisterRequest{
		Name:     "lisi",
		Password: "123456"},
	)
	if err != nil {
		return nil, err
	}

	fmt.Printf("register success %v\n\n", resp)
	return resp, nil
}

func getUser(client pb.AccountClient, req *pb.RegisterReply) error {
	md := metadata.Pairs("authorization", req.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Id: req.Id})
	if err != nil {
		return err
	}

	fmt.Println("get user success", resp)
	return nil
}

func getDialOptions(isUseTLS bool) []grpc.DialOption {
	var options []grpc.DialOption

	if !isUseTLS {
		// 不使用加密传输
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		// tls加密
		credentials, err := gtls.GetClientTLSCredentials("localhost", certfile.Path("one-way/server.crt"))
		if err != nil {
			panic(err)
		}
		options = append(options, grpc.WithTransportCredentials(credentials))
	}

	return options
}

func main() {
	isUseTLS = true // 设置是否需要TLS加密

	conn, err := grpc.Dial("127.0.0.1:9090", getDialOptions(isUseTLS)...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewAccountClient(conn)

	resp, err := registerUser(client)
	if err != nil {
		panic(err)
	}

	err = getUser(client, resp)
	if err != nil {
		panic(err)
	}
}
