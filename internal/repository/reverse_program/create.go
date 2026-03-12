package reverse_program

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/internal/model"

	"gorm.io/gorm"
)

// Create 创建一条逆向程序记录。
func Create(ctx context.Context, db *gorm.DB, item *dbmodel.ReverseProgram) error {
	if db == nil {
		return fmt.Errorf("gorm db is not initialized")
	}

	return db.WithContext(ctx).Create(item).Error
}
