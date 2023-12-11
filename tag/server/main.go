package main

import (
	"context"
	"fmt"
	"net"
	"reflect"

	pb "github.com/zhufuyi/grpc_examples/tag/proto/hellopb"

	"google.golang.org/grpc"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	t := reflect.TypeOf(*r)

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fmt.Println(sf.Tag.Get("json"), sf.Tag.Get("gorm"), sf.Tag.Get("bson"))
	}

	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	return options
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
	server := grpc.NewServer(getServerOptions()...)

	// register greeterServer to the server
	pb.RegisterGreeterServer(server, &greeterServer{})

	// start the server
	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
