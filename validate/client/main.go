package main

import (
	"context"
	"fmt"

	pb "github.com/zhufuyi/grpc_examples/validate/proto/accountpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func login(client pb.AccountClient, req *pb.LoginRequest) error {
	resp, err := client.Login(context.Background(), req)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	fmt.Printf("login success, token=%s\n\n", resp.Token)
	return nil
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close() //nolint

	client := pb.NewAccountClient(conn)

	// 所有字段正确
	req := &pb.LoginRequest{
		Id:       10001,
		Email:    "foo@bar.com",
		Password: "123456",
		Phone:    "13566666666",
	}
	if err := login(client, req); err != nil {
		fmt.Println(err)
	}

	// email 字段错误
	req = &pb.LoginRequest{
		Id:       10001,
		Email:    "foo",
		Password: "123456",
		Phone:    "13566666666",
	}
	if err := login(client, req); err != nil {
		fmt.Println(err)
	}

	// password 字段错误
	req = &pb.LoginRequest{
		Id:       10001,
		Email:    "foo@bar.com",
		Password: "apcdef",
		Phone:    "13566666666",
	}
	if err := login(client, req); err != nil {
		fmt.Println(err)
	}

	// phone 字段错误
	req = &pb.LoginRequest{
		Id:       10001,
		Email:    "foo@bar.com",
		Password: "123456",
		Phone:    "1234567890",
	}
	if err := login(client, req); err != nil {
		fmt.Println(err)
	}
}
