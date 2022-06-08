package tracer

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/zhufuyi/grpc_examples/pkg/tracer/otgrpc"
)

// GinCtx 把gin的context转换为标准的context
func GinCtx(c *gin.Context) context.Context {
	tracerVal, _ := c.Get(otgrpc.GinTracerKey)
	ctx := context.WithValue(context.Background(), otgrpc.GinTracerKey, tracerVal)

	parentSpanVal, _ := c.Get(otgrpc.GinParentSpanKey)
	return context.WithValue(ctx, otgrpc.GinParentSpanKey, parentSpanVal)
}

// GinMiddleware gin的链路跟踪中间件
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//closer, err := InitJaeger("tracing_demo", "192.168.3.36:6831")
		//if err != nil {
		//	panic(err)
		//}
		//defer closer.Close()

		startSpan := opentracing.StartSpan(c.Request.URL.Path)
		defer startSpan.Finish()

		// 通过gin的context传递的tracer和startSpan给下一个使用者
		c.Set(otgrpc.GinTracerKey, opentracing.GlobalTracer())
		c.Set(otgrpc.GinParentSpanKey, startSpan)

		c.Next()
	}
}

// Get 在http请求中添加链路追踪
func Get(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	span, newCtx := opentracing.StartSpanFromContext(ctx, "HTTP GET: "+url, opentracing.Tag{Key: string(ext.Component), Value: "HTTP"})
	defer span.Finish()

	req = req.WithContext(newCtx)
	client := http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
