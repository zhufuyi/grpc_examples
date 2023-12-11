package main

import (
	"context"
	"fmt"
	"math"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/keepalive/proto/hellopb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var (
	infinity = time.Duration(math.MaxInt64)

	kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     infinity,         // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      infinity,         // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: infinity,         // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  20 * time.Second, // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}
)

// ServerKeepAlive 保持连接设置
func ServerKeepAlive() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
	}
}

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func main() {
	addr := ":8282"
	fmt.Println("grpc service is running", addr)

	var options []grpc.ServerOption
	// 默认是每15秒向客户端发送一次ping，修改为间隔20秒发送一次ping
	options = append(options, ServerKeepAlive()...)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(options...)

	pb.RegisterGreeterServer(server, &greeterServer{})

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
