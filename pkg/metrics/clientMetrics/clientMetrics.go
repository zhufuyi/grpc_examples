package clientMetrics

import (
	"net/http"
	"sync"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// https://github.com/grpc-ecosystem/go-grpc-prometheus/tree/master/examples/grpc-server-with-prometheus

var (
	// 创建一个Registry
	cliReg = prometheus.NewRegistry()

	// 初始化客户端默认的metrics
	grpcClientMetrics = grpc_prometheus.NewClientMetrics()

	// 执行一次
	once sync.Once
)

func registerMetrics() {
	once.Do(func() {
		// 注册metrics才能进行采集，自定义的metrics也需要注册
		cliReg.MustRegister(grpcClientMetrics)
	})
}

// Serve 初始化客户端的prometheus的exporter服务，使用 http://ip:port/metrics 获取数据
func Serve(addr string) {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: promhttp.HandlerFor(cliReg, promhttp.HandlerOpts{}),
	}

	// 启动http服务
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			panic("Unable to start a http server.")
		}
	}()
}

// ---------------------------------- client interceptor ----------------------------------

// UnaryClientMetrics 一元rpc metrics
func UnaryClientMetrics() grpc.UnaryClientInterceptor {
	registerMetrics() // 在拦截器之前完成注册metrics，只执行一次
	return grpcClientMetrics.UnaryClientInterceptor()
}

// StreamClientMetrics 流rpc metrics
func StreamClientMetrics() grpc.StreamClientInterceptor {
	registerMetrics() // 在拦截器之前完成注册metrics，只执行一次
	return grpcClientMetrics.StreamClientInterceptor()
}
