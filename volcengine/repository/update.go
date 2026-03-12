package repository

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/volcengine/model"
	"strings"

	"gorm.io/gorm"
)

// UpdateAPIByID 按 ID 更新模型 API 配置。
func UpdateAPIByID(ctx context.Context, db *gorm.DB, id string, input dbmodel.ModelAPIModel) (dbmodel.ModelAPIModel, error) {
	if db == nil {
		return dbmodel.ModelAPIModel{}, fmt.Errorf("gorm db is not initialized")
	}

	existing, err := GetAPIByID(ctx, db, id)
	if err != nil {
		return dbmodel.ModelAPIModel{}, err
	}

	existing.Name = strings.TrimSpace(input.Name)
	existing.Provider = strings.TrimSpace(input.Provider)
	existing.BaseURL = strings.TrimSpace(input.BaseURL)
	existing.APIKey = strings.TrimSpace(input.APIKey)
	existing.Model = strings.TrimSpace(input.Model)
	existing.APIModel = strings.TrimSpace(input.APIModel)
	existing.Enabled = input.Enabled

	if err = db.WithContext(ctx).Save(&existing).Error; err != nil {
		return dbmodel.ModelAPIModel{}, err
	}

	return existing, nil
}
