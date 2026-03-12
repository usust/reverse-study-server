package repository

import (
	"context"
	"fmt"
	dbmodel "reverse-study-server/volcengine/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateAPI 创建一条模型 API 配置。
func CreateAPI(ctx context.Context, db *gorm.DB, api *dbmodel.ModelAPIModel) error {
	if db == nil {
		return fmt.Errorf("gorm db is not initialized")
	}

	return db.WithContext(ctx).Create(api).Error
}

// Upsert 按 ID 插入或更新模型 API 配置。
func Upsert(ctx context.Context, db *gorm.DB, api *dbmodel.ModelAPIModel) error {
	if db == nil {
		return fmt.Errorf("gorm db is not initialized")
	}

	return db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(api).Error
}
