package reverse_program

import (
	"net/http"
	compiler "reverse-study-server/compiler"
	dbmodel "reverse-study-server/internal/model"
	reverseprogramsvc "reverse-study-server/internal/service/reverse_program"
	chatapi "reverse-study-server/volcengine/chat"
	"strings"

	"github.com/gin-gonic/gin"
)

type createNewProgramRequest struct {
	ModelAPIID     string                         `json:"model_api_id"`
	Prompt         string                         `json:"prompt"`
	CompileOptions compiler.CompileRequestOptions `json:"compile_options"`
	Program        createNewProgramInfo           `json:"program"`
}

type createNewProgramInfo struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Published       int      `json:"published"`
	SourceFileName  string   `json:"sourceFileName"`
	ProgramFileName string   `json:"programFileName"`
	Score           int      `json:"score"`
	ProgramType     int      `json:"programType"`
	Difficulty      int      `json:"difficulty"`
	Tags            []string `json:"tags"`
	BaseDir         string   `json:"baseDir"`
	CompletedCount  int      `json:"completedCount"`
}

// NewProgram 调用 service 创建新的逆向题目。
func NewProgram(c *gin.Context) {
	var input createNewProgramRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid json body", "status": http.StatusBadRequest})
		return
	}

	err := reverseprogramsvc.GenerateNewReverseProgram(
		c.Request.Context(),
		chatapi.PromptChatRequest{
			ModelAPIID: input.ModelAPIID,
			Prompt:     input.Prompt,
		},
		input.CompileOptions,
		dbmodel.ReverseProgram{
			Title:           strings.TrimSpace(input.Program.Title),
			Description:     strings.TrimSpace(input.Program.Description),
			Published:       input.Program.Published,
			SourceFileName:  strings.TrimSpace(input.Program.SourceFileName),
			ProgramFileName: strings.TrimSpace(input.Program.ProgramFileName),
			Score:           input.Program.Score,
			ProgramType:     input.Program.ProgramType,
			Difficulty:      input.Program.Difficulty,
			Tags:            normalizeTags(input.Program.Tags),
			BaseDir:         strings.TrimSpace(input.Program.BaseDir),
			CompletedCount:  input.Program.CompletedCount,
			MetaFileName:    "meta.json",
		},
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "status": http.StatusBadRequest})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "题目创建成功",
		"status":  http.StatusOK,
		"data": gin.H{
			"baseDir": input.Program.BaseDir,
		},
	})
}

func normalizeTags(tags []string) []string {
	var result []string
	for _, item := range tags {
		tag := strings.TrimSpace(item)
		if tag == "" {
			continue
		}
		result = append(result, tag)
	}
	return result
}
