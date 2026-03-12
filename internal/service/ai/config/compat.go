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

const createCCodePromptName = "创建C代码提示词"

// ProgramSavePathConfig 兼容旧接口，实际读取当前运行时 storage.base_dir。
type ProgramSavePathConfig struct {
	ProgramSavePath string `json:"programSavePath"`
}

// CreateCCodePromptConfig 兼容旧接口，实际读写 prompts 表中的“创建C代码提示词”。
type CreateCCodePromptConfig struct {
	Prompt string `json:"prompt"`
}

// GetProgramSavePathConfig 获取题目保存路径配置。
func GetProgramSavePathConfig(ctx context.Context) (ProgramSavePathConfig, error) {
	item, err := GetStorageConfig(ctx)
	if err != nil {
		return ProgramSavePathConfig{}, err
	}
	return ProgramSavePathConfig{ProgramSavePath: item.BaseDir}, nil
}

// UpdateProgramSavePathConfig 更新题目保存路径配置。
func UpdateProgramSavePathConfig(ctx context.Context, path string) (ProgramSavePathConfig, error) {
	item, err := UpdateStorageConfig(ctx, path)
	if err != nil {
		return ProgramSavePathConfig{}, err
	}
	return ProgramSavePathConfig{ProgramSavePath: item.BaseDir}, nil
}

// GetCreateCCodePromptConfig 获取“创建C代码提示词”。
func GetCreateCCodePromptConfig(ctx context.Context) (CreateCCodePromptConfig, error) {
	item, err := promptrepo.GetByName(ctx, bootstrap.GormDB, createCCodePromptName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return CreateCCodePromptConfig{Prompt: ""}, nil
		}
		return CreateCCodePromptConfig{}, err
	}
	return CreateCCodePromptConfig{Prompt: strings.TrimSpace(item.Content)}, nil
}

// UpdateCreateCCodePromptConfig 更新“创建C代码提示词”。
func UpdateCreateCCodePromptConfig(ctx context.Context, prompt string) (CreateCCodePromptConfig, error) {
	content := strings.TrimSpace(prompt)
	if content == "" {
		return CreateCCodePromptConfig{}, fmt.Errorf("prompt is required")
	}

	existing, err := promptrepo.GetByName(ctx, bootstrap.GormDB, createCCodePromptName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item := dbmodel.Prompt{
				Name:    createCCodePromptName,
				Content: content,
			}
			if createErr := promptrepo.Create(ctx, bootstrap.GormDB, &item); createErr != nil {
				return CreateCCodePromptConfig{}, createErr
			}
			return CreateCCodePromptConfig{Prompt: content}, nil
		}
		return CreateCCodePromptConfig{}, err
	}

	existing.Content = content
	if _, err = promptrepo.UpdateByID(ctx, bootstrap.GormDB, existing.ID, existing); err != nil {
		return CreateCCodePromptConfig{}, err
	}
	return CreateCCodePromptConfig{Prompt: content}, nil
}
