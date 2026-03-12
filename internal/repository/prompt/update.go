package prompt

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/internal/model"
	"strings"

	"gorm.io/gorm"
)

// UpdateByID 按 ID 更新提示词。
func UpdateByID(ctx context.Context, db *gorm.DB, id uint64, input dbmodel.Prompt) (dbmodel.Prompt, error) {
	if db == nil {
		return dbmodel.Prompt{}, fmt.Errorf("gorm db is not initialized")
	}

	existing, err := GetByID(ctx, db, id)
	if err != nil {
		return dbmodel.Prompt{}, err
	}

	existing.Name = strings.TrimSpace(input.Name)
	existing.Content = strings.TrimSpace(input.Content)

	if err = db.WithContext(ctx).Save(&existing).Error; err != nil {
		return dbmodel.Prompt{}, err
	}
	return existing, nil
}
