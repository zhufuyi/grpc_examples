package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	pb "github.com/zhufuyi/grpc_examples/metrics/defaultMetrics/proto/hellopb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func sayHello(client pb.GreeterClient, i int) error {
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "zhangsan " + strconv.Itoa(i)})
	if err != nil {
		return err
	}
	fmt.Println(resp.Message)

	return nil
}

// 启动metrics服务
func defaultDefaultServer(addr string) {
	grpc_prometheus.EnableHandlingTimeHistogram() // 开启了对RPCs处理时间的记录

	http.Handle("/metrics", promhttp.Handler())
	go func() {
		fmt.Printf("metrics server started on %s\n", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic(err)
		}
	}()
}

func getDialOptions() []grpc.DialOption {
	var options []grpc.DialOption

	// 禁用tls
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	// default Metrics
	options = append(options, grpc.WithUnaryInterceptor(
		grpc_prometheus.UnaryClientInterceptor,
	))
	options = append(options, grpc.WithStreamInterceptor(
		grpc_prometheus.StreamClientInterceptor,
	))
	return options
}

func main() {
	grpc_prometheus.EnableHandlingTimeHistogram()
	conn, err := grpc.Dial("127.0.0.1:8080", getDialOptions()...)

	defaultDefaultServer(":9094")

	client := pb.NewGreeterClient(conn)
	i := 0
	for {
		i++
		time.Sleep(time.Millisecond * 500) // qps is 2
		err = sayHello(client, i)
		if err != nil {
			fmt.Println(err)
		}
	}
}
