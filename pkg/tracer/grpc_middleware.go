package tracer

import (
	"grpc_examples/pkg/tracer/otgrpc"

	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

// 使用拦截器之前先调用函数 InitJaeger(name, addr) 初始化，否则opentracing.GlobalTracer()无效

// ---------------------------------- server interceptor ----------------------------------

// UnaryServerTracing 一元rpc的tracing服务端拦截器
func UnaryServerTracing() grpc.UnaryServerInterceptor {
	return otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())
}

// ---------------------------------- client interceptor ----------------------------------

// UnaryClientTracing 一元rpc的tracing客户端拦截器
func UnaryClientTracing() grpc.UnaryClientInterceptor {
	return otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())
}
