package repository

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/volcengine/model"
	"strings"

	"gorm.io/gorm"
)

// DeleteAPIByID 按 ID 删除模型 API 配置。
func DeleteAPIByID(ctx context.Context, db *gorm.DB, id string) error {
	if db == nil {
		return fmt.Errorf("gorm db is not initialized")
	}

	trimmedID := strings.TrimSpace(id)
	if trimmedID == "" {
		return fmt.Errorf("id is required")
	}

	result := db.WithContext(ctx).Where("id = ?", trimmedID).Delete(&dbmodel.ModelAPIModel{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
