package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pb "github.com/zhufuyi/grpc_examples/metrics/defaultMetrics/proto/hellopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *GreeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	n := rand.Intn(100)
	switch {
	case n%10 == 0: // 大概10%错误率
		fmt.Println("n =", n, "set deadlineExceeded error")
		return nil, status.Errorf(codes.DeadlineExceeded, "time out")
	case n%5 == 0: // 大概20%延时超过250ms
		time.Sleep(time.Millisecond * 255)
	}

	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

// 启动metrics服务
func defaultDefaultServer(addr string, grpcServer *grpc.Server) {
	grpc_prometheus.EnableHandlingTimeHistogram() // 开启了对RPCs处理时间的记录
	grpc_prometheus.Register(grpcServer)          // 注册rpc服务到metrics

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		fmt.Printf("metrics server started on %s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}()
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	option := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_prometheus.UnaryServerInterceptor, // 一元rpc的metrics拦截器
	))
	options = append(options, option)

	option = grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_prometheus.StreamServerInterceptor, // 流式rpc的metrics拦截器
	))
	options = append(options, option)

	return options
}

func main() {
	rand.Seed(time.Now().UnixNano())

	addr := ":8080"
	fmt.Println("start rpc server", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions()...)
	pb.RegisterGreeterServer(server, &GreeterServer{})

	// 启动metrics服务器，默认采集grpc指标，开启RPCs处理时间的记录、go指标
	defaultDefaultServer(":9092", server)

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
