package main

import (
	"fmt"
	"net/http"

	"github.com/zhufuyi/grpc_examples/include"
	"github.com/zhufuyi/grpc_examples/swagger"
)

var swaggerAddr = ":9090"

func main() {
	router := swagger.Router("/", include.Path("../swagger/example/proto/hellopb/hello.swagger.json"))
	server := &http.Server{
		Addr:    swaggerAddr,
		Handler: router,
	}

	fmt.Println("start up web server ", swaggerAddr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
