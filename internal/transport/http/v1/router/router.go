package router

import (
	v1reverseprogramrouter "reverse-study-server/internal/transport/http/v1/router/reverse_program"

	configrouter "reverse-study-server/internal/transport/http/v1/router/config"

	"github.com/gin-gonic/gin"
)

// RegisterV1 注册 v1 版本下的一级路由。
func RegisterV1(r gin.IRouter) {
	v1Group := r.Group("/v1")

	// 题目生成
	registerReverseProgramRouters(v1Group)
	// 配置菜单
	registerConfigRouters(v1Group)
}

// 注册配置
func registerConfigRouters(r gin.IRouter) {
	configRouter := r.Group("/config")
	configrouter.RegisterConfigAPIRoutes(configRouter)
	configrouter.RegisterConfigPromptRoutes(configRouter)
	configrouter.RegisterConfigCommonRoutes(configRouter)
}

func registerReverseProgramRouters(r gin.IRouter) {
	v1reverseprogramrouter.RegisterReverseProgramRouter(r)
}
