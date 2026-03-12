package reverse_program

import (
	"context"
	"fmt"
	"reverse-study-server/internal/bootstrap"
	dbmodel "reverse-study-server/internal/model"
	reverseprogramrepo "reverse-study-server/internal/repository/reverse_program"
	"strings"
)

type UpdateProgramInput struct {
	Title           string
	Description     string
	Published       bool
	Score           int
	ProgramType     string
	Difficulty      string
	Tags            []string
	BaseDir         string
	SourceFileName  string
	ProgramFileName string
	ProgramFileMD5  string
	CompletedCount  int
}

// UpdateProgramByID 更新指定题目。
func UpdateProgramByID(ctx context.Context, id uint64, input UpdateProgramInput) (dbmodel.ReverseProgram, error) {
	if id == 0 {
		return dbmodel.ReverseProgram{}, fmt.Errorf("id is required")
	}

	return reverseprogramrepo.UpdateByID(ctx, bootstrap.GormDB, id, dbmodel.ReverseProgram{
		Title:           strings.TrimSpace(input.Title),
		Description:     strings.TrimSpace(input.Description),
		Published:       input.Published,
		Score:           input.Score,
		ProgramType:     strings.TrimSpace(input.ProgramType),
		Difficulty:      strings.TrimSpace(input.Difficulty),
		Tags:            normalizeTags(input.Tags),
		BaseDir:         strings.TrimSpace(input.BaseDir),
		SourceFileName:  strings.TrimSpace(input.SourceFileName),
		ProgramFileName: strings.TrimSpace(input.ProgramFileName),
		ProgramFileMD5:  strings.TrimSpace(input.ProgramFileMD5),
		CompletedCount:  input.CompletedCount,
	})
}

// PublishProgramByID 发布指定题目。
func PublishProgramByID(ctx context.Context, id uint64) (dbmodel.ReverseProgram, error) {
	if id == 0 {
		return dbmodel.ReverseProgram{}, fmt.Errorf("id is required")
	}

	item, err := GetProgramByID(ctx, id)
	if err != nil {
		return dbmodel.ReverseProgram{}, err
	}

	if item.Published {
		return item, nil
	}

	item.Published = true
	return reverseprogramrepo.UpdateByID(ctx, bootstrap.GormDB, id, item)
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
