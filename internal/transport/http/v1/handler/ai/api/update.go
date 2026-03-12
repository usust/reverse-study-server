package api

import (
	"errors"
	"net/http"
	v1service "reverse-study-server/internal/service/ai/api"
	dbmodel "reverse-study-server/volcengine/model"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateModelAPIInfo 更新模型 API 配置。
func UpdateModelAPIInfo(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id is required", "status": http.StatusBadRequest})
		return
	}

	var input dbmodel.ModelAPIModel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	item, err := v1service.UpdateModelAPI(c.Request.Context(), id, input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "ModelAPI not found", "status": http.StatusNotFound})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "status": http.StatusOK, "data": item})
	return
}
