package config

import (
	"net/http"
	configservice "reverse-study-server/internal/service/ai/config"

	"github.com/gin-gonic/gin"
)

type updateCreateCCodePromptRequest struct {
	Prompt string `json:"prompt"`
}

// GetCreateCCodePrompt 获取“创建C代码提示词”配置。
func GetCreateCCodePrompt(c *gin.Context) {
	item, err := configservice.GetCreateCCodePromptConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "status": http.StatusOK, "data": item})
}

// UpdateCreateCCodePrompt 更新“创建C代码提示词”配置。
func UpdateCreateCCodePrompt(c *gin.Context) {
	var input updateCreateCCodePromptRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	item, err := configservice.UpdateCreateCCodePromptConfig(c.Request.Context(), input.Prompt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "status": http.StatusOK, "data": item})
}
