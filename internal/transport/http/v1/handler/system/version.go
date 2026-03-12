package system

import (
	"net/http"
	v1service "reverse-study-server/internal/service/system"

	"github.com/gin-gonic/gin"
)

// GetVersion 返回服务版本信息。
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, v1service.GetVersionInfo())
}
