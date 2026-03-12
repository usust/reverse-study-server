package config

import (
	"context"
	"errors"
	"fmt"
	"reverse-study-server/internal/bootstrap"
	dbmodel "reverse-study-server/internal/model"
	promptrepo "reverse-study-server/internal/repository/prompt"
	"strings"

	"gorm.io/gorm"
)

// ListPrompts 返回提示词列表。
func ListPrompts(ctx context.Context) ([]dbmodel.Prompt, error) {
	return promptrepo.List(ctx, bootstrap.GormDB)
}

// SavePrompt 保存提示词（有 ID 则更新，无 ID 则按 name 创建或更新）。
func SavePrompt(ctx context.Context, input dbmodel.Prompt) (dbmodel.Prompt, error) {
	name := strings.TrimSpace(input.Name)
	content := strings.TrimSpace(input.Content)
	if name == "" {
		return dbmodel.Prompt{}, fmt.Errorf("name is required")
	}
	if content == "" {
		return dbmodel.Prompt{}, fmt.Errorf("content is required")
	}

	if input.ID > 0 {
		input.Name = name
		input.Content = content
		return promptrepo.UpdateByID(ctx, bootstrap.GormDB, input.ID, input)
	}

	existing, err := promptrepo.GetByName(ctx, bootstrap.GormDB, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item := dbmodel.Prompt{
				Name:    name,
				Content: content,
			}
			if createErr := promptrepo.Create(ctx, bootstrap.GormDB, &item); createErr != nil {
				return dbmodel.Prompt{}, createErr
			}
			return item, nil
		}
		return dbmodel.Prompt{}, err
	}

	existing.Content = content
	return promptrepo.UpdateByID(ctx, bootstrap.GormDB, existing.ID, existing)
}
