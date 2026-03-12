package api

import (
	"errors"
	"net/http"
	v1service "reverse-study-server/internal/service/ai/api"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAPIByID 获取单个模型 API 配置。
func GetAPIByID(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id is required", "status": http.StatusBadRequest})
		return
	}

	item, err := v1service.GetAPIInfoByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "config not found", "status": http.StatusNotFound})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "status": http.StatusOK, "data": item})
}

// ListAPIs 获取模型 API 配置列表。
func ListAPIs(c *gin.Context) {
	items, err := v1service.ListAPIs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "status": http.StatusOK, "data": items})
	return
}
