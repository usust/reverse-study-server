package api

import (
	"net/http"
	v1service "reverse-study-server/internal/service/ai/api"
	m "reverse-study-server/volcengine/model"

	"github.com/gin-gonic/gin"
)

// CreateModelAPI 创建模型 API 配置。
func CreateModelAPI(c *gin.Context) {
	var input m.ModelAPIModel
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	item, err := v1service.CreateAPI(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "创建成功", "status": http.StatusOK, "data": item})
}
