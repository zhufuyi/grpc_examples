package tracer

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

// InitJaeger 连接jaeger服务
func InitJaeger(serviceName string, agentAddr string) (io.Closer, error) {
	cfg := &jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: time.Second,
			LocalAgentHostPort:  agentAddr,
		},
	}

	jLogger := jaegerlog.StdLogger
	option := []jaegercfg.Option{
		jaegercfg.Logger(jLogger),
	}

	tracer, closer, err := cfg.NewTracer(option...)
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer) // 免去设置全局tracer变量，在其他地方使用opentracing.GlobalTracer()获取tracer，opentracing.StartSpan()创建一个span

	return closer, nil
}

//func InitZipkinTracer(serviceName string,rpcServerAddr, zipkinAddr string) (opentracing.Tracer, error) {
//	// docker run -d -p 9411:9411 -p 9410:9410 openzipkin/zipkin
//
//	// create collector.
//	//collector, err := zipkin.NewHTTPCollector("http://localhost:9411/api/v1/spans")
//	collector, err := zipkin.NewHTTPCollector(zipkinAddr+"/api/v1/spans")
//	if err != nil {
//		return nil, err
//	}
//
//	// create recorder.
//	//recorder := zipkin.NewRecorder(collector, false, "127.0.0.1:50051", serviceName)
//	recorder := zipkin.NewRecorder(collector, false, rpcServerAddr, serviceName)
//
//	// create tracer.
//	tracer, err := zipkin.NewTracer(
//		recorder,
//		zipkin.ClientServerSameSpan(true),
//		zipkin.TraceID128Bit(true),
//	)
//	if err != nil {
//		nil, err
//	}
//
//	// explicitly set our tracer to be the default tracer.
//	//opentracing.InitGlobalTracer(tracer)
//
//	return tracer,nil
//}
