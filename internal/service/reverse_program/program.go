package reverse_program

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"reverse-study-server/internal/bootstrap"
	dbmodel "reverse-study-server/internal/model"
	reverseprogramrepo "reverse-study-server/internal/repository/reverse_program"
	configservice "reverse-study-server/internal/service/ai/config"
	"strings"

	"github.com/google/uuid"
)

// CreateProgramInput 是创建逆向程序记录的输入。
type CreateProgramInput struct {
	Prompt     string
	CCode      string
	Score      int
	Difficulty string
}

// CreateProgram 保存一条逆向程序记录。
func CreateProgram(ctx context.Context, input CreateProgramInput) (dbmodel.ReverseProgram, error) {
	code := strings.TrimSpace(input.CCode)
	if code == "" {
		return dbmodel.ReverseProgram{}, fmt.Errorf("c code is required")
	}

	score := input.Score
	if score <= 0 {
		score = 100
	}
	difficulty := strings.TrimSpace(input.Difficulty)
	if difficulty == "" {
		difficulty = "medium"
	}

	sourceName := "main.c"
	hash := md5.Sum([]byte(code))
	sourceMD5 := hex.EncodeToString(hash[:])

	programSavePathConfig, err := configservice.GetProgramSavePathConfig(ctx)
	if err != nil {
		return dbmodel.ReverseProgram{}, err
	}

	saveDir := strings.TrimSpace(programSavePathConfig.ProgramSavePath)
	if saveDir == "" {
		saveDir = "./data/reverse-programs"
	}
	if err := os.MkdirAll(saveDir, 0o755); err != nil {
		return dbmodel.ReverseProgram{}, fmt.Errorf("create program save directory failed: %w", err)
	}

	sourceDir := filepath.Join(saveDir, uuid.NewString())
	if err := os.MkdirAll(sourceDir, 0o755); err != nil {
		return dbmodel.ReverseProgram{}, fmt.Errorf("create program source directory failed: %w", err)
	}

	sourcePath := filepath.Join(sourceDir, sourceName)
	if err := os.WriteFile(sourcePath, []byte(code), 0o644); err != nil {
		return dbmodel.ReverseProgram{}, fmt.Errorf("write program source file failed: %w", err)
	}

	item := dbmodel.ReverseProgram{
		Title:           "未命名题目",
		Description:     "",
		Published:       false,
		SourceFileName:  sourceName,
		ProgramFileName: sourceName,
		Score:           score,
		ProgramType:     "",
		Difficulty:      difficulty,
		Tags:            nil,
		BaseDir:         sourceDir,
		ProgramFileMD5:  sourceMD5,
		CompletedCount:  0,
	}
	if err := reverseprogramrepo.Create(ctx, bootstrap.GormDB, &item); err != nil {
		return dbmodel.ReverseProgram{}, err
	}
	return item, nil
}

// GetProgramByID 按 ID 查询逆向程序记录。
func GetProgramByID(ctx context.Context, id uint64) (dbmodel.ReverseProgram, error) {
	return reverseprogramrepo.GetByID(ctx, bootstrap.GormDB, id)
}

// ListPrograms 查询题目列表。
func ListPrograms(ctx context.Context) ([]dbmodel.ReverseProgram, error) {
	return reverseprogramrepo.List(ctx, bootstrap.GormDB)
}
