package reverse_program

import (
	"net/http"
	reverseprogramsvc "reverse-study-server/internal/service/reverse_program"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteProgram 删除指定题目。
func DeleteProgram(c *gin.Context) {
	idText := c.Param("id")
	id, err := strconv.ParseUint(idText, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id", "status": http.StatusBadRequest})
		return
	}

	if err := reverseprogramsvc.DeleteProgramByID(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
}

