package reverse_program

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/internal/model"
	"strings"

	"gorm.io/gorm"
)

// UpdateByID 按 ID 更新逆向程序记录。
func UpdateByID(ctx context.Context, db *gorm.DB, id uint64, input dbmodel.ReverseProgram) (dbmodel.ReverseProgram, error) {
	if db == nil {
		return dbmodel.ReverseProgram{}, fmt.Errorf("gorm db is not initialized")
	}

	existing, err := GetByID(ctx, db, id)
	if err != nil {
		return dbmodel.ReverseProgram{}, err
	}

	existing.Title = strings.TrimSpace(input.Title)
	existing.Description = strings.TrimSpace(input.Description)
	existing.Published = input.Published
	existing.SourceFileName = strings.TrimSpace(input.SourceFileName)
	existing.ProgramFileName = strings.TrimSpace(input.ProgramFileName)
	existing.Score = input.Score
	existing.ProgramType = strings.TrimSpace(input.ProgramType)
	existing.Difficulty = strings.TrimSpace(input.Difficulty)
	existing.Tags = append([]string(nil), input.Tags...)
	existing.BaseDir = strings.TrimSpace(input.BaseDir)
	existing.ProgramFileMD5 = strings.TrimSpace(input.ProgramFileMD5)
	existing.CompletedCount = input.CompletedCount

	if err = db.WithContext(ctx).Save(&existing).Error; err != nil {
		return dbmodel.ReverseProgram{}, err
	}

	return existing, nil
}
