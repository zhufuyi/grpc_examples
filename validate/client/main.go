package main

import (
	"context"
	"fmt"

	pb "grpc_examples/validate/proto/accountpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func login(client pb.AccountClient, email string, password string) error {
	resp, err := client.Login(context.Background(), &pb.LoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return fmt.Errorf("%v\n", err)
	}

	fmt.Printf("login success, token=%s\n\n", resp.Token)
	return nil
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewAccountClient(conn)

	var email, password string

	email = "zhangsan@126.com"
	password = "123456"
	if err := login(client, email, password); err != nil {
		fmt.Println(err)
	}

	email = "zhangsan"
	password = "123456"
	if err := login(client, email, password); err != nil {
		fmt.Println(err)
	}

	email = "zhangsan@126.com"
	password = "abcdef"
	if err := login(client, email, password); err != nil {
		fmt.Println(err)
	}
}
