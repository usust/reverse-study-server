package api

import (
	"context"
	"reverse-study-server/internal/bootstrap"
	m "reverse-study-server/volcengine/model"
	rep "reverse-study-server/volcengine/repository"
)

// GetAPIInfoByID 按 ID 查询模型 API 配置。
func GetAPIInfoByID(ctx context.Context, id string) (m.ModelAPIModel, error) {
	return rep.GetAPIByID(ctx, bootstrap.GormDB, id)
}

// ListAPIs 查询全部模型 API 配置。
func ListAPIs(ctx context.Context) ([]m.ModelAPIModel, error) {
	return rep.ListAPIs(ctx, bootstrap.GormDB)
}
