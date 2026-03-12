package config

import (
	"reverse-study-server/internal/transport/http/v1/handler/ai/config"

	"github.com/gin-gonic/gin"
)

func RegisterConfigCommonRoutes(r gin.IRouter) {
	commonGroup := r.Group("/common")
	{
		commonGroup.GET("/storage/base", config.GetStorageConfig)
		commonGroup.PUT("/storage/base", config.UpdateStorageConfig)
	}
}
