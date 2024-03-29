package main

import (
	"context"
	"net/http"

	"github.com/zhufuyi/grpc_examples/tracing/http2rpc"
	pb "github.com/zhufuyi/grpc_examples/tracing/http2rpc/proto/hellopb"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/tracer"
)

var (
	// web addr
	webAddr = ":8080"

	// rpc addr
	rpcAddr = "127.0.0.1:8282"

	serviceName = "hello-client"
)

func sayHello(c *gin.Context) {
	resp, err := helloClient.SayHello(c.Request.Context(), &pb.HelloRequest{Name: "foo"})
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

func main() {
	tracing.InitTrace(serviceName)
	defer tracer.Close(context.Background()) //nolint

	// 连接rpc服务端
	connectRPCServer(rpcAddr)

	r := gin.Default()
	r.Use(middleware.Tracing(serviceName))
	r.GET("/hello", sayHello)

	err := r.Run(webAddr)
	if err != nil {
		panic(err)
	}
}
