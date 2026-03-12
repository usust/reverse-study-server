package http

import (
	v1 "reverse-study-server/internal/transport/http/v1/router"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Content-Disposition, X-Generated-File-Name")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// InitRouter 创建 HTTP 路由并挂载各版本路由。
func InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())

	v1.RegisterV1(r)
	return r
}
