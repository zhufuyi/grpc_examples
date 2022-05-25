package middleware

import (
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// ---------------------------------- server interceptor ----------------------------------

var (
	// 自定义打印kv
	loggingFields = map[string]interface{}{}

	// 跳过打印日志的调用方法
	skipLoggingMethods = map[string]struct{}{}
)

// AddSkipMethods 添加跳过认证方法，在服务初始化时设置
func AddLoggingFields(kvs map[string]interface{}) {
	loggingFields = kvs
}

// AddSkipLoggingMethods 添加跳过打印日志方法，在服务初始化时设置
func AddSkipLoggingMethods(methodNames ...string) {
	for _, name := range methodNames {
		skipLoggingMethods[name] = struct{}{}
	}
}

// UnaryServerZapLogging 日志打印拦截器
func UnaryServerZapLogging(logger *zap.Logger) grpc.UnaryServerInterceptor {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// 日志设置，默认打印客户端断开连接信息，示例 https://pkg.go.dev/github.com/grpc-ecosystem/go-grpc-middleware/logging/zap
	zapOptions := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds()) // 默认打印耗时字段
		}),
	}

	// 自定义打印字段
	for key, val := range loggingFields {
		zapOptions = append(zapOptions, grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Any(key, val)
		}))
	}

	// 自定义跳过打印日志的调用方法
	for method := range skipLoggingMethods {
		zapOptions = append(zapOptions, grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
			if err == nil && fullMethodName == method {
				return false
			}
			return true
		}))
	}

	return grpc_zap.UnaryServerInterceptor(logger, zapOptions...)
}

// UnaryServerCtxTags field extractor logging
func UnaryServerCtxTags() grpc.UnaryServerInterceptor {
	return grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor))
}
