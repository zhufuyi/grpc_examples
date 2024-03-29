package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	pb "github.com/zhufuyi/grpc_examples/helloworld/proto/hellopb"

	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

// 一元RPC
func (g *greeterServer) UnarySayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fmt.Println("unarySayHello receive: " + req.Name)
	msg := "hello " + req.Name
	fmt.Println("unarySayHello send   : " + msg)
	return &pb.HelloReply{Message: msg}, nil
}

// 服务端流式RPC
func (g *greeterServer) ServerStreamingSayHello(req *pb.HelloRequest, stream pb.Greeter_ServerStreamingSayHelloServer) error {
	recValues := req.Name
	sendValues := []string{}

	defer func() {
		fmt.Println("\nserverStreamingSayHello receive: ", recValues)
		fmt.Println("serverStreamingSayHello send   : ", sendValues)
	}()

	for i := 0; i < 3; i++ {
		sendMsg := fmt.Sprintf("hello %s %d", req.Name, i)
		sendValues = append(sendValues, sendMsg)
		err := stream.Send(&pb.HelloReply{Message: sendMsg})
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}

	return nil
}

// 客户端流式RPC
func (g *greeterServer) ClientStreamingSayHello(stream pb.Greeter_ClientStreamingSayHelloServer) error {
	recValues := []string{}
	sendValues := ""

	defer func() {
		fmt.Println("\nclientStreamingSayHello receive: ", recValues)
		fmt.Println("clientStreamingSayHello send   : ", sendValues)
	}()

	for {
		// 阻塞等待接收流数据，当结束时会受到EOF表示结束，当出现错误会返回rpc错误信息
		// 默认的MaxReceiveMessageSize值为1024x1024x4字节，如果有特殊需求可以调整
		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF { // 判断是否数据流结束
				sendValues = "hello " + strings.Join(recValues, ",")
				return stream.SendAndClose(&pb.HelloReply{
					Message: sendValues,
				})
			}
			return err
		}

		recValues = append(recValues, resp.Name)
	}
}

// 双向流式RPC
func (g *greeterServer) BidirectionalStreamingSayHello(stream pb.Greeter_BidirectionalStreamingSayHelloServer) error {
	recValues := []string{}
	sendValues := []string{}

	defer func() {
		fmt.Println("\nbidirectionalStreamingSayHello receive: ", recValues)
		fmt.Println("bidirectionalStreamingSayHello send   : ", sendValues)
	}()

	var resp *pb.HelloRequest
	var err error
	for {
		resp, err = stream.Recv()
		if err != nil {
			if err == io.EOF { // 判断是否数据流结束
				return nil
			}
			return err
		}
		recValues = append(recValues, resp.Name)

		sendMsg := "hello " + resp.Name
		err = stream.Send(&pb.HelloReply{Message: sendMsg})
		if err != nil {
			return err
		}
		sendValues = append(sendValues, sendMsg)
	}
}

func main() {
	addr := ":8282"
	fmt.Println("grpc service is running", addr)

	// listening on TCP port
	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	// create a grpc server object where interceptors can be injected
	server := grpc.NewServer()

	// register greeterServer to the server
	pb.RegisterGreeterServer(server, &greeterServer{})

	// start the server
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
