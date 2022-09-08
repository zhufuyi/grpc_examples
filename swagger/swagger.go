package swagger

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
)

// Router swagger 路由
func Router(path string, jsonFile string) *gin.Engine {
	path = strings.TrimRight(path, "/") + "/"

	r := gin.Default()

	registerSwagger(jsonFile)
	// 访问路径 http://ip:port/<path>/swagger/index.html
	r.GET(path+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// Handler swagger handler
func Handler(r *gin.Engine, jsonFile string) {
	registerSwagger(jsonFile)
	// 访问路径 http://ip:port/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// 注册swagger
func registerSwagger(jsonFile string) {
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	swaggerInfo := &swag.Spec{
		Schemes:          []string{},
		InfoInstanceName: "swagger",
		SwaggerTemplate:  string(data),
	}
	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)
}
