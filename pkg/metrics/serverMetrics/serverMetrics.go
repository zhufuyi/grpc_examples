package serverMetrics

import (
	"net/http"
	"sync"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

// https://github.com/grpc-ecosystem/go-grpc-prometheus/tree/master/examples/grpc-server-with-prometheus

var (
	// 创建一个Registry
	srvReg = prometheus.NewRegistry()

	// 初始化服务端默认的metrics
	grpcServerMetrics = grpc_prometheus.NewServerMetrics()

	// go metrics
	goMetrics = collectors.NewGoCollector()

	// 用户自定义指标 https://prometheus.io/docs/concepts/metric_types/#histogram
	customizedCounterMetrics   = []*prometheus.CounterVec{}
	customizedSummaryMetrics   = []*prometheus.SummaryVec{}
	customizedGaugeMetrics     = []*prometheus.GaugeVec{}
	customizedHistogramMetrics = []*prometheus.HistogramVec{}

	// 执行一次
	once sync.Once
)

// AddCounterMetrics 添加Counter类型指标
func AddCounterMetrics(metrics ...*prometheus.CounterVec) {
	customizedCounterMetrics = append(customizedCounterMetrics, metrics...)
}

// AddSummaryMetrics 添加Summary类型指标
func AddSummaryMetrics(metrics ...*prometheus.SummaryVec) {
	customizedSummaryMetrics = append(customizedSummaryMetrics, metrics...)
}

// AddGaugeMetrics 添加Gauge类型指标
func AddGaugeMetrics(metrics ...*prometheus.GaugeVec) {
	customizedGaugeMetrics = append(customizedGaugeMetrics, metrics...)
}

// AddHistogramMetrics 添加Histogram类型指标
func AddHistogramMetrics(metrics ...*prometheus.HistogramVec) {
	customizedHistogramMetrics = append(customizedHistogramMetrics, metrics...)
}

func registerMetrics() {
	once.Do(func() {
		// 开启了对RPCs处理时间的记录
		grpcServerMetrics.EnableHandlingTimeHistogram()

		// 注册go metrics
		srvReg.MustRegister(goMetrics)

		// 注册metrics才能进行采集，自定义的metrics也需要注册
		srvReg.MustRegister(grpcServerMetrics)

		// 注册自定义counter metric
		for _, metric := range customizedCounterMetrics {
			srvReg.MustRegister(metric)
		}
		for _, metric := range customizedSummaryMetrics {
			srvReg.MustRegister(metric)
		}
		for _, metric := range customizedGaugeMetrics {
			srvReg.MustRegister(metric)
		}
		for _, metric := range customizedHistogramMetrics {
			srvReg.MustRegister(metric)
		}
	})
}

// Serve 初始化服务端的prometheus的exporter服务，使用 http://ip:port/metrics 获取数据
func Serve(addr string, grpcServer *grpc.Server) {
	httpServer := &http.Server{
		Addr:    addr,
		Handler: promhttp.HandlerFor(srvReg, promhttp.HandlerOpts{}),
	}

	// 启动http服务
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			panic("Unable to start a http server.")
		}
	}()

	// 所有gRPC方法初始化Metrics
	grpcServerMetrics.InitializeMetrics(grpcServer)
}

// ---------------------------------- server interceptor ----------------------------------

// UnaryServerMetrics 一元rpc metrics
func UnaryServerMetrics() grpc.UnaryServerInterceptor {
	registerMetrics() // 在拦截器之前完成注册metrics，只执行一次
	return grpcServerMetrics.UnaryServerInterceptor()
}

// StreamServerMetrics 流rpc metrics
func StreamServerMetrics() grpc.StreamServerInterceptor {
	registerMetrics() // 在拦截器之前完成注册metrics，只执行一次
	return grpcServerMetrics.StreamServerInterceptor()
}
