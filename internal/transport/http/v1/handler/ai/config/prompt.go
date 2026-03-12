package config

import (
	"net/http"
	dbmodel "reverse-study-server/internal/model"
	v1service "reverse-study-server/internal/service/ai/config"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type savePromptRequest struct {
	ID      uint64 `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// ListPrompts 获取提示词列表。
func ListPrompts(c *gin.Context) {
	items, err := v1service.ListPrompts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "查询成功", "status": http.StatusOK, "data": items})
}

// SavePrompt 创建或更新提示词。
func SavePrompt(c *gin.Context) {
	var input savePromptRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	item, err := v1service.SavePrompt(c.Request.Context(), dbmodel.Prompt{
		ID:      input.ID,
		Name:    strings.TrimSpace(input.Name),
		Content: strings.TrimSpace(input.Content),
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "保存成功", "status": http.StatusOK, "data": item})
}

// DeletePrompt 删除提示词。
func DeletePrompt(c *gin.Context) {
	idText := strings.TrimSpace(c.Param("id"))
	id, err := strconv.ParseUint(idText, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id is required", "status": http.StatusBadRequest})
		return
	}

	if err := v1service.DeletePromptByID(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功", "status": http.StatusOK})
}
