package chat

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	chatprompt "reverse-study-server/internal/service/config/prompt"
	reverseprogramsvc "reverse-study-server/internal/service/reverse_program"
	chatapi "reverse-study-server/volcengine/chat"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCCode 接收前端提示词，并调用模型生成 C 源码。
func CreateCCode(c *gin.Context) {
	var input chatapi.PromptChatRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	result, err := chatprompt.GenerateCCode(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	program, err := reverseprogramsvc.CreateProgram(c.Request.Context(), reverseprogramsvc.CreateProgramInput{
		Prompt: input.Prompt,
		CCode:  result.Response.Content,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "生成成功",
		"status":  http.StatusOK,
		"data": gin.H{
			"promptResult": result,
			"program":      program,
		},
	})
}

// DownloadProgramSource 下载已保存的 C 源码。
func DownloadProgramSource(c *gin.Context) {
	idText := strings.TrimSpace(c.Param("id"))
	if idText == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id is required", "status": http.StatusBadRequest})
		return
	}
	id, err := strconv.ParseUint(idText, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "id must be a positive integer", "status": http.StatusBadRequest})
		return
	}

	item, err := reverseprogramsvc.GetProgramByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "program not found", "status": http.StatusNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}

	baseDir := strings.TrimSpace(item.BaseDir)
	if baseDir == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "program baseDir is empty", "status": http.StatusInternalServerError})
		return
	}
	sourceName := strings.TrimSpace(item.SourceFileName)
	if sourceName == "" {
		sourceName = "main.c"
	}
	sourcePath := filepath.Join(baseDir, sourceName)
	if strings.TrimSpace(sourcePath) == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "program source path is empty", "status": http.StatusInternalServerError})
		return
	}
	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	if len(sourceData) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "program source is empty", "status": http.StatusInternalServerError})
		return
	}

	fileName := sourceName
	if fileName == "" {
		fileName = fmt.Sprintf("reverse_program_%d.c", item.ID)
	}
	fileHash := md5.Sum(sourceData)
	sourceMD5 := hex.EncodeToString(fileHash[:])
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.Header("Content-Type", "text/x-c; charset=utf-8")
	c.Header("X-File-MD5", sourceMD5)
	c.Data(http.StatusOK, "text/x-c; charset=utf-8", sourceData)
}
