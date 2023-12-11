package main

import (
	"fmt"
	"net/http"

	"github.com/zhufuyi/grpc_examples/include"
	"github.com/zhufuyi/grpc_examples/swagger"

	"github.com/gin-gonic/gin"
)

var swaggerAddr = ":8080"

func main() {
	router := swagger.Router("/", include.Path("../swagger/example/proto/hellopb/hello.swagger.json"))
	router.POST("/v1/sayHello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello",
		})
	})
	server := &http.Server{
		Addr:    swaggerAddr,
		Handler: router,
	}

	fmt.Println("http service is running", swaggerAddr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
