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
