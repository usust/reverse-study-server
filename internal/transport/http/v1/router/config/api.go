package config

import (
	"reverse-study-server/internal/transport/http/v1/handler/ai/api"

	"github.com/gin-gonic/gin"
)

func RegisterConfigAPIRoutes(r gin.IRouter) {
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/list", api.ListAPIs)
		apiGroup.GET("/:id", api.GetAPIByID)
		apiGroup.POST("/create", api.CreateModelAPI)
		apiGroup.PUT("/:id", api.UpdateModelAPIInfo)
		apiGroup.DELETE("/:id", api.DeleteModelAPI)
	}
}
