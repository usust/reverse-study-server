package api

import (
	"errors"
	"net/http"
	v1service "reverse-study-server/internal/service/ai/api"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DeleteModelAPI 删除模型 API 配置。
func DeleteModelAPI(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id is required", "status": http.StatusBadRequest})
		return
	}

	if err := v1service.DeleteAPIService(c.Request.Context(), id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "api not found", "status": http.StatusNotFound})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
	return
}
