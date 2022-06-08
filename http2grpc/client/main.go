package main

import (
	"context"
	"fmt"

	pb "github.com/zhufuyi/grpc_examples/http2grpc/proto/accountpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func addUser(client pb.AccountClient) error {
	resp, err := client.AddUser(context.Background(), &pb.User{Id: 1, Name: "zhangsan", Email: "zhangsan@126.com"})
	if err != nil {
		return err
	}

	fmt.Println("add user success", resp.Id)
	return nil
}

func getUser(client pb.AccountClient) error {
	resp, err := client.GetUser(context.Background(), &pb.ID{Id: 1})
	if err != nil {
		return err
	}

	fmt.Println("get user success", resp)
	return nil
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewAccountClient(conn)

	err = addUser(client)
	if err != nil {
		panic(err)
	}

	err = getUser(client)
	if err != nil {
		panic(err)
	}
}
