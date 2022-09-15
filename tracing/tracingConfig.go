package tracing

import (
	"context"
	"time"

	"github.com/zhufuyi/pkg/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	sn = ""
)

// InitTrace 初始化链路跟踪
func InitTrace(serviceName string) {
	sn = serviceName
	//exporter, _, err := tracer.NewFileExporter(serviceName + "-trace.json")
	//exporter, err := tracer.NewJaegerExporter("http://192.168.3.37:14268/api/traces")
	exporter, err := tracer.NewJaegerAgentExporter("192.168.3.37", "6831")
	if err != nil {
		panic(err)
	}

	resource := tracer.NewResource(
		tracer.WithServiceName(serviceName),
		tracer.WithEnvironment("dev"),
		tracer.WithServiceVersion("demo"),
	)

	tracer.Init(exporter, resource) // 默认采集全部
}

// SpanDemo 设置一个span
func SpanDemo(ctx context.Context, spanName string) {
	_, span := otel.Tracer(sn).Start(
		ctx, spanName,
		trace.WithAttributes(attribute.String(spanName, time.Now().String())),
	)
	defer span.End()

	time.Sleep(20 * time.Millisecond)
}
