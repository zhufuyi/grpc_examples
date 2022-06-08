package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/grpc_examples/pkg/tracer"
	pb "github.com/zhufuyi/grpc_examples/tracing/api2rpc/proto/hellopb"
)

var (
	// web addr
	webAddr = ":6060"

	// jager agent addr
	serviceName = "tracing_demo"
	agentAddr   = "192.168.3.36:6831"

	// rpc addr
	rpcAddr = "127.0.0.1:8080"
)

func sayHello(c *gin.Context) {
	// tracer.GinCtx(c) 从gin的context获取tracing相关字段，生成新的context
	// 调用rpc1
	resp, err := helloClient.SayHello(tracer.GinCtx(c), &pb.HelloRequest{Name: "zhangsan"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func main() {
	// 连接jaeger
	closer, err := tracer.InitJaeger(serviceName, agentAddr)
	if err != nil {
		panic(err)
	}
	defer closer.Close()

	// 连接rpc服务端
	connectRPCServer(rpcAddr)

	r := gin.Default()
	r.Use(tracer.GinMiddleware()) // 添加链路跟踪中间件
	r.POST("/hello", sayHello)

	err = r.Run(webAddr)
	if err != nil {
		panic(err)
	}
}
