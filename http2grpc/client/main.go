package main

import (
	"context"
	"fmt"
	"time"

	"github.com/zhufuyi/grpc_examples/http2grpc/proto/pb"

	"github.com/zhufuyi/pkg/krand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createUser(client pb.UserServiceClient) (int64, error) {
	name := krand.String(krand.R_LOWER)
	resp, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{
		Name:  name,
		Email: name + "@bar.com",
	})
	if err != nil {
		return 0, err
	}

	fmt.Println("add user success, id =", resp.Id)
	return resp.Id, nil
}

func getUser(client pb.UserServiceClient, id int64) error {
	resp, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
	if err != nil {
		return err
	}

	fmt.Println("get user success", resp)
	return nil
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close() //nolint

	client := pb.NewUserServiceClient(conn)

	for i := 0; i < 10; i++ {
		id, err := createUser(client)
		if err != nil {
			panic(err)
		}

		err = getUser(client, id)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 2)
	}
}
