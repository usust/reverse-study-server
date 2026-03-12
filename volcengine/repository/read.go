package repository

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/volcengine/model"
	"strings"

	"gorm.io/gorm"
)

// GetAPIByID 按 ID 查询模型 API 配置。
func GetAPIByID(ctx context.Context, db *gorm.DB, id string) (dbmodel.ModelAPIModel, error) {
	if db == nil {
		return dbmodel.ModelAPIModel{}, fmt.Errorf("gorm db is not initialized")
	}

	trimmedID := strings.TrimSpace(id)
	if trimmedID == "" {
		return dbmodel.ModelAPIModel{}, fmt.Errorf("id is required")
	}

	var item dbmodel.ModelAPIModel
	if err := db.WithContext(ctx).Where("id = ?", trimmedID).First(&item).Error; err != nil {
		return dbmodel.ModelAPIModel{}, err
	}

	return item, nil
}

// ListAPIs 查询全部模型 API 配置，按更新时间倒序返回。
func ListAPIs(ctx context.Context, db *gorm.DB) ([]dbmodel.ModelAPIModel, error) {
	if db == nil {
		return nil, fmt.Errorf("gorm db is not initialized")
	}

	var items []dbmodel.ModelAPIModel
	if err := db.WithContext(ctx).Order("updated_at desc").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
