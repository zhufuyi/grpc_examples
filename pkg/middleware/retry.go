package middleware

import (
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// ---------------------------------- client option ----------------------------------

var (
	defaultTimes    uint = 3                                                                       // 重试次数
	defaultInterval      = time.Millisecond * 50                                                   // 重试间隔50毫秒
	defaultCodes         = []codes.Code{codes.Internal, codes.DeadlineExceeded, codes.Unavailable} // 默认触发重试的错误码
)

// SetRetryTimes 设置重试次数，最大10次
func SetRetryTimes(n uint) {
	if n > 10 {
		n = 10
	}
	defaultTimes = n
}

// SetRetryInterval 设置重试时间间隔
func SetRetryInterval(t time.Duration) {
	if t < time.Millisecond {
		t = time.Millisecond
	} else if t > time.Second {
		t = time.Second
	}
	defaultInterval = t
}

// SetRetryErrCodes 设置触发重试错误码
func SetRetryErrCodes(errCodes ...codes.Code) {
	for _, errCode := range errCodes {
		switch errCode {
		case codes.Internal, codes.DeadlineExceeded, codes.Unavailable:
		default:
			defaultCodes = append(defaultCodes, errCode)
		}
	}
}

// Retry 重试
func Retry() grpc.UnaryClientInterceptor {
	return grpc_retry.UnaryClientInterceptor(
		grpc_retry.WithMax(defaultTimes), // 设置重试次数
		grpc_retry.WithBackoff(func(attempt uint) time.Duration { // 设置重试间隔
			return defaultInterval
		}),
		grpc_retry.WithCodes(codes.Internal, codes.DeadlineExceeded, codes.Unavailable), // 设置重试错误码
	)
}
