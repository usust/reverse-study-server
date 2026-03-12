package config

import (
	"reverse-study-server/internal/transport/http/v1/handler/ai/config"

	"github.com/gin-gonic/gin"
)

func RegisterConfigPromptRoutes(r gin.IRouter) {
	configGroup := r.Group("/prompt")
	{
		configGroup.GET("/prompts", config.ListPrompts)
		configGroup.PUT("/prompts", config.SavePrompt)
		configGroup.DELETE("/prompts/:id", config.DeletePrompt)
	}

	commonGroup := r.Group("/common")
	{
		commonGroup.GET("/storage", config.GetStorageConfig)
		commonGroup.PUT("/storage", config.UpdateStorageConfig)
	}
}
