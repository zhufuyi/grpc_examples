package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/loadbalance/client_loadbalance/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	timestampFormat = "2006-01-02T15:04:05.000000"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	// return之后创建新的metadata
	defer func() {
		trailer := metadata.Pairs("timestamp", time.Now().Format(timestampFormat)+"S2")
		err := grpc.SetTrailer(ctx, trailer)
		if err != nil {
			fmt.Println(err)
		}
	}()

	// 从client读取metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "SayHello: failed to get metadata")
	}
	printMetadata(md, "timestamp")

	// 创建新的metadata
	header := metadata.New(map[string]string{
		"hello":     "world",
		"timestamp": time.Now().Format(timestampFormat) + "S1",
	})
	_ = grpc.SendHeader(ctx, header)

	// 随机延时0~1000微秒
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Microsecond)

	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func printMetadata(md metadata.MD, key string) {
	if ts, ok := md[key]; ok {
		fmt.Printf("metadata: %s=%v\n", key, ts)
	} else {
		fmt.Printf("not found '%s' in metadata\n", key)
	}
}

func main() {
	addr := ":8282"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	pb.RegisterGreeterServer(server, &greeterServer{})

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
