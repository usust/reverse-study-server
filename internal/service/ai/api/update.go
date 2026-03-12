package api

import (
	"context"
	"reverse-study-server/internal/bootstrap"
	dbmodel "reverse-study-server/volcengine/model"
	rep "reverse-study-server/volcengine/repository"
)

// UpdateModelAPI 更新指定 ID 的模型 API 配置。
func UpdateModelAPI(ctx context.Context, id string, input dbmodel.ModelAPIModel) (dbmodel.ModelAPIModel, error) {
	return rep.UpdateAPIByID(ctx, bootstrap.GormDB, id, input)
}
