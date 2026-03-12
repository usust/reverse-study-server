package prompt

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/internal/model"

	"gorm.io/gorm"
)

// DeleteByID 按 ID 删除提示词。
func DeleteByID(ctx context.Context, db *gorm.DB, id uint64) error {
	if db == nil {
		return fmt.Errorf("gorm db is not initialized")
	}
	if id == 0 {
		return fmt.Errorf("id is required")
	}

	result := db.WithContext(ctx).Where("id = ?", id).Delete(&dbmodel.Prompt{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

