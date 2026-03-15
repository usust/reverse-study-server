package system

import (
	v1handler "reverse-study-server/internal/transport/http/v1/handler/system"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 汇总当前 router 子包下的全部路由。
func RegisterRoutes(r gin.IRouter) {
	systemGroup := r.Group("/system")
	{
		systemGroup.GET("/version", v1handler.GetVersion)
	}
}
