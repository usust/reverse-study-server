package config

import (
	"net/http"
	configservice "reverse-study-server/internal/service/config/storage"

	"github.com/gin-gonic/gin"
)

type updateStorageConfigRequest struct {
	BaseDir string `json:"base_dir"`
}

// GetStorageConfig 获取通用存储配置。
func GetStorageConfig(c *gin.Context) {
	item, err := configservice.GetStorageConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "status": http.StatusOK, "data": item})
}

// UpdateStorageConfig 更新通用存储配置。
func UpdateStorageConfig(c *gin.Context) {
	var input updateStorageConfigRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	item, err := configservice.UpdateStorageConfig(c.Request.Context(), input.BaseDir)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "status": http.StatusOK, "data": item})
}
