package config

import (
	"net/http"
	configservice "reverse-study-server/internal/service/ai/config"

	"github.com/gin-gonic/gin"
)

type updateProgramSavePathRequest struct {
	ProgramSavePath string `json:"programSavePath"`
}

// GetProgramSavePath 获取题目保存路径配置。
func GetProgramSavePath(c *gin.Context) {
	item, err := configservice.GetProgramSavePathConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "status": http.StatusOK, "data": item})
}

// UpdateProgramSavePath 更新题目保存路径配置。
func UpdateProgramSavePath(c *gin.Context) {
	var input updateProgramSavePathRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	item, err := configservice.UpdateProgramSavePathConfig(c.Request.Context(), input.ProgramSavePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "status": http.StatusOK, "data": item})
}
