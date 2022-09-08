package routers

import (
	"github.com/zhufuyi/grpc_examples/include"
	"github.com/zhufuyi/grpc_examples/swagger"
	"github.com/zhufuyi/pkg/gin/middleware"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/pkg/logger"
)

// NewRouter 实例化路由
func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.Logging(middleware.WithLog(logger.Get())))

	// 注册swagger路由，通过swag init生成代码
	swagger.Handler(r, include.Path("../gin2grpc/gin/api/user/v1/pb/user.swagger.json"))

	return r
}
