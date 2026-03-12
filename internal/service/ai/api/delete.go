package api

import (
	"context"
	"reverse-study-server/internal/bootstrap"
	"reverse-study-server/volcengine/repository"
)

// DeleteAPIService 删除指定 ID 的模型 API 配置。
func DeleteAPIService(ctx context.Context, id string) error {
	return repository.DeleteAPIByID(ctx, bootstrap.GormDB, id)
}
