package reverse_program

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	reverseprogramsvc "reverse-study-server/internal/service/reverse_program"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DownloadProgram 下载题目的程序文件，文件名使用程序 MD5。
func DownloadProgram(c *gin.Context) {
	idText := strings.TrimSpace(c.Param("id"))
	id, err := strconv.ParseUint(idText, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id", "status": http.StatusBadRequest})
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
	programName := strings.TrimSpace(item.ProgramFileName)
	if programName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "program file name is empty", "status": http.StatusInternalServerError})
		return
	}
	fileMD5 := strings.TrimSpace(item.ProgramFileMD5)
	if fileMD5 == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "program file md5 is empty", "status": http.StatusInternalServerError})
		return
	}

	programPath := filepath.Join(baseDir, "programs", strconv.FormatUint(item.ID, 10), programName)
	data, err := os.ReadFile(programPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "status": http.StatusInternalServerError})
		return
	}
	if len(data) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "program file is empty", "status": http.StatusInternalServerError})
		return
	}

	downloadName := fileMD5 + filepath.Ext(programName)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", downloadName))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("X-File-MD5", fileMD5)
	c.Data(http.StatusOK, "application/octet-stream", data)
}
