package reverse_program

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/internal/model"

	"gorm.io/gorm"
)

// DeleteByID 按 ID 删除逆向程序记录。
func DeleteByID(ctx context.Context, db *gorm.DB, id uint64) error {
	if db == nil {
		return fmt.Errorf("gorm db is not initialized")
	}

	if id == 0 {
		return fmt.Errorf("id is required")
	}

	result := db.WithContext(ctx).Where("id = ?", id).Delete(&dbmodel.ReverseProgram{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
