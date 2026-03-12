package prompt

import (
	"context"
	"errors"
	"fmt"
	dbmodel "reverse-study-server/internal/model"
	"strings"

	"gorm.io/gorm"
)

// GetByID 按 ID 查询提示词。
func GetByID(ctx context.Context, db *gorm.DB, id uint64) (dbmodel.Prompt, error) {
	if db == nil {
		return dbmodel.Prompt{}, fmt.Errorf("gorm db is not initialized")
	}
	if id == 0 {
		return dbmodel.Prompt{}, fmt.Errorf("id is required")
	}

	var item dbmodel.Prompt
	if err := db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return dbmodel.Prompt{}, err
	}
	return item, nil
}

// GetByName 按名称查询提示词。
func GetByName(ctx context.Context, db *gorm.DB, name string) (dbmodel.Prompt, error) {
	if db == nil {
		return dbmodel.Prompt{}, fmt.Errorf("gorm db is not initialized")
	}

	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return dbmodel.Prompt{}, fmt.Errorf("name is required")
	}

	var item dbmodel.Prompt
	if err := db.WithContext(ctx).Where("name = ?", trimmedName).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dbmodel.Prompt{}, err
		}
		return dbmodel.Prompt{}, err
	}
	return item, nil
}

// List 按更新时间倒序查询提示词列表。
func List(ctx context.Context, db *gorm.DB) ([]dbmodel.Prompt, error) {
	if db == nil {
		return nil, fmt.Errorf("gorm db is not initialized")
	}

	var items []dbmodel.Prompt
	if err := db.WithContext(ctx).Order("updated_at desc").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
