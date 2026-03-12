package prompt

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/internal/model"

	"gorm.io/gorm"
)

// Create 创建一条提示词记录。
func Create(ctx context.Context, db *gorm.DB, item *dbmodel.Prompt) error {
	if db == nil {
		return fmt.Errorf("gorm db is not initialized")
	}
	if item == nil {
		return fmt.Errorf("prompt item is required")
	}

	return db.WithContext(ctx).Create(item).Error
}

