package routers

import (
	"github.com/zhufuyi/grpc_examples/include"
	"github.com/zhufuyi/grpc_examples/swagger"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/logger"
)

// NewRouter 实例化路由
func NewRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.Logging(middleware.WithLog(logger.Get())))

	// 注册swagger路由，通过swag init生成代码
	swagger.Handler(r, include.Path("../httpToGrpc/ginToGrpc/web-server/api/user/v1/pb/user.swagger.json"))

	return r
}
