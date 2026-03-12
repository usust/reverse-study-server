package reverse_program

import (
	"net/http"
	reverseprogramsvc "reverse-study-server/internal/service/reverse_program"

	"github.com/gin-gonic/gin"
)

// ListPrograms 返回题目列表。
func ListPrograms(c *gin.Context) {
	items, err := reverseprogramsvc.ListPrograms(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "查询成功",
		"status":  http.StatusOK,
		"data":    items,
	})
}
