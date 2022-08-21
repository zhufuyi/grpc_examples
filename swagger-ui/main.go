package main

import (
	"fmt"
	"net/http"

	jsonfile "github.com/zhufuyi/grpc_examples/swagger-ui/proto"
	"github.com/zhufuyi/grpc_examples/swagger-ui/swagger"
)

var swaggerAddr = "127.0.0.1:8080"

func main() {
	mux := http.NewServeMux()

	// 注册swagger-ui和swagger json文件路由
	swagger.RegisterRoute(mux, jsonfile.Path("/hellopb"))

	fmt.Println("start up swagger server ", swaggerAddr)
	err := http.ListenAndServe(swaggerAddr, mux)
	if err != nil {
		panic(err)
	}
}

/*
(1) 方式一，使用自身swagger服务
启动之后，在浏览器打开 http://127.0.0.1:8080/swagger-ui
然后输入查看本地json文件 http://127.0.0.1:8080/swagger/hello.swagger.json


(2) 方式二，使用官方swag-ui指定json文件查看
docker run -p 8080:8080 -e SWAGGER_JSON=/openapiv2/hello.swagger.json -v $PWD/openapiv2/:/openapiv2 swaggerapi/swagger-ui:v4.14.0

-e SWAGGER_JSON值表示指定swagger json文件，
-v 把json文件到容器/openapiv2目录下，

启动容器，在浏览器打开 http://127.0.0.1:8080 就可以自动进入接口文档
*/
