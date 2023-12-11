package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	pb "github.com/zhufuyi/grpc_examples/metrics/customMetrics/proto/hellopb"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/zhufuyi/sponge/pkg/grpc/interceptor"
	"github.com/zhufuyi/sponge/pkg/grpc/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	serverNameLabelKey   = "name"
	serverNameLabelValue = "hello"
	envLabelKey          = "env"
	envLabelValue        = "dev"
)

var (
	// 自定义Counter指标
	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "demo_server_count",
		Help: "Total number of RPCs handled on the server.",
	}, []string{serverNameLabelKey, envLabelKey})
)

func counterMetricInc(ctx context.Context, counter *prometheus.CounterVec) {
	//if grpc_ctxtags.Extract(ctx).Has(serverNameLabelKey) {
	//	val := grpc_ctxtags.Extract(ctx).Values()[serverNameLabelKey]
	//	fmt.Println(val)
	//}

	counter.WithLabelValues(serverNameLabelValue, envLabelValue).Inc() // values 一一对应key
}

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (g *greeterServer) SayHello(ctx context.Context, r *pb.HelloRequest) (*pb.HelloReply, error) {
	n := rand.Intn(100)
	switch {
	case n%10 == 0: // 大概10%错误率
		fmt.Println("n =", n, "set deadlineExceeded error")
		return nil, status.Errorf(codes.DeadlineExceeded, "time out")
	case n%5 == 0: // 大概20%延时超过250ms
		time.Sleep(time.Millisecond * 255)
	}

	counterMetricInc(ctx, customizedCounterMetric)

	return &pb.HelloReply{Message: "hello " + r.Name}, nil
}

// UnaryServerLabels 设置标签拦截器
func UnaryServerLabels(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 设置prometheus公共标签
	tag := grpc_ctxtags.NewTags().
		Set(serverNameLabelKey, serverNameLabelValue).
		Set(envLabelKey, envLabelValue)
	newCtx := grpc_ctxtags.SetInContext(ctx, tag)

	return handler(newCtx, req)
}

func getServerOptions() []grpc.ServerOption {
	var options []grpc.ServerOption

	// metrics拦截器
	option := grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		UnaryServerLabels,                // 标签
		interceptor.UnaryServerMetrics(), // 一元rpc的metrics拦截器
	))
	options = append(options, option)

	option = grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		interceptor.StreamServerMetrics(), // 流式rpc的metrics拦截器
	))
	options = append(options, option)

	return options
}

func main() {
	rand.Seed(time.Now().UnixNano())

	addr := ":8282"
	fmt.Println("grpc service is running", addr)

	list, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(getServerOptions()...)
	pb.RegisterGreeterServer(server, &greeterServer{})

	// 启动metrics服务器，默认采集grpc指标，开启、go指标
	metricsAddr := ":8283"
	metrics.ServerHTTPService(metricsAddr, server)
	fmt.Println("start metrics server " + metricsAddr)

	err = server.Serve(list)
	if err != nil {
		panic(err)
	}
}
