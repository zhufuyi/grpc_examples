package main

import (
	"context"
	"fmt"
	"io"
	"time"

	pb "github.com/zhufuyi/grpc_examples/helloworld/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func unarySayHello(client pb.GreeterClient) error {
	resp, err := client.UnarySayHello(context.Background(), &pb.HelloRequest{Name: "foo"})
	if err != nil {
		return err
	}

	fmt.Println("resp:", resp.Message)
	return nil
}

func serverStreamingSayHello(client pb.GreeterClient) error {
	stream, err := client.ServerStreamingSayHello(context.Background(), &pb.HelloRequest{Name: "foo"})
	if err != nil {
		return err
	}

	for {
		// 阻塞等待接收流数据，当结束时会受到EOF表示结束，当出现错误会返回rpc错误信息
		// 默认的MaxReceiveMessageSize值为1024x1024x4字节，如果有特殊需求可以调整
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // 判断是否数据流结束
				break
			}
			return err
		}

		fmt.Println("resp:", resp.Message)
	}

	return nil
}

func clientStreamingSayHello(client pb.GreeterClient) error {
	stream, err := client.ClientStreamingSayHello(context.Background())
	if err != nil {
		return err
	}

	names := []string{"foo1", "foo2", "foo3"}
	for _, name := range names {
		err = stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Println("resp:", resp.Message)

	return nil
}

func bidirectionalStreamingSayHello(client pb.GreeterClient) error {
	stream, err := client.BidirectionalStreamingSayHello(context.Background())
	if err != nil {
		return err
	}

	names := []string{"foo1", "foo2", "foo3"}
	var resp *pb.HelloReply
	for _, name := range names {
		err = stream.Send(&pb.HelloRequest{Name: name})
		if err != nil {
			return err
		}

		resp, err = stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fmt.Println("resp:", resp.Message)
	}

	time.Sleep(10 * time.Millisecond)
	err = stream.CloseSend()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8282", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client := pb.NewGreeterClient(conn)

	fmt.Println("\n一元RPC调用示例：unarySayHello")
	if err := unarySayHello(client); err != nil {
		panic(err)
	}

	fmt.Println("\n服务端流式RPC：serverStreamingSayHello")
	if err := serverStreamingSayHello(client); err != nil {
		panic(err)
	}

	fmt.Println("\n客户端流式RPC：clientStreamingSayHello")
	if err := clientStreamingSayHello(client); err != nil {
		panic(err)
	}

	fmt.Println("\n双向流式RPC：bidirectionalStreamingSayHello")
	if err := bidirectionalStreamingSayHello(client); err != nil {
		panic(err)
	}
}
