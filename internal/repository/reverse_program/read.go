package reverse_program

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/internal/model"

	"gorm.io/gorm"
)

// GetByID 按 ID 查询逆向程序记录。
func GetByID(ctx context.Context, db *gorm.DB, id uint64) (dbmodel.ReverseProgram, error) {
	if db == nil {
		return dbmodel.ReverseProgram{}, fmt.Errorf("gorm db is not initialized")
	}

	if id == 0 {
		return dbmodel.ReverseProgram{}, fmt.Errorf("id is required")
	}

	var item dbmodel.ReverseProgram
	if err := db.WithContext(ctx).Where("id = ?", id).First(&item).Error; err != nil {
		return dbmodel.ReverseProgram{}, err
	}

	return item, nil
}

// List 按更新时间倒序查询题目列表。
func List(ctx context.Context, db *gorm.DB) ([]dbmodel.ReverseProgram, error) {
	if db == nil {
		return nil, fmt.Errorf("gorm db is not initialized")
	}

	var items []dbmodel.ReverseProgram
	if err := db.WithContext(ctx).Order("updated_at desc").Find(&items).Error; err != nil {
		return nil, err
	}

	return items, nil
}
