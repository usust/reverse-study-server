package reverse_program

import (
	"net/http"
	reverseprogramsvc "reverse-study-server/internal/service/reverse_program"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type updateProgramRequest struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Published       bool     `json:"published"`
	Score           int      `json:"score"`
	ProgramType     string   `json:"programType"`
	Difficulty      string   `json:"difficulty"`
	Tags            []string `json:"tags"`
	BaseDir         string   `json:"baseDir"`
	SourceFileName  string   `json:"sourceFileName"`
	ProgramFileName string   `json:"programFileName"`
	ProgramFileMD5  string   `json:"programFileMd5"`
	CompletedCount  int      `json:"completedCount"`
}

// UpdateProgram 更新指定题目。
func UpdateProgram(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id", "status": http.StatusBadRequest})
		return
	}

	var input updateProgramRequest
	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	item, err := reverseprogramsvc.UpdateProgramByID(c.Request.Context(), id, reverseprogramsvc.UpdateProgramInput{
		Title:           strings.TrimSpace(input.Title),
		Description:     strings.TrimSpace(input.Description),
		Published:       input.Published,
		Score:           input.Score,
		ProgramType:     strings.TrimSpace(input.ProgramType),
		Difficulty:      strings.TrimSpace(input.Difficulty),
		Tags:            input.Tags,
		BaseDir:         strings.TrimSpace(input.BaseDir),
		SourceFileName:  strings.TrimSpace(input.SourceFileName),
		ProgramFileName: strings.TrimSpace(input.ProgramFileName),
		ProgramFileMD5:  strings.TrimSpace(input.ProgramFileMD5),
		CompletedCount:  input.CompletedCount,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"status":  http.StatusOK,
		"data":    item,
	})
}

// PublishProgram 发布指定题目。
func PublishProgram(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid id", "status": http.StatusBadRequest})
		return
	}

	item, err := reverseprogramsvc.PublishProgramByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "发布成功",
		"status":  http.StatusOK,
		"data":    item,
	})
}
