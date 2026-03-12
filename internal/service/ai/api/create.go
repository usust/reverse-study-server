package api

import (
	"context"
	"fmt"
	"reverse-study-server/internal/bootstrap"
	m "reverse-study-server/volcengine/model"
	r "reverse-study-server/volcengine/repository"
	"strings"

	"github.com/google/uuid"
)

// CreateAPI 创建模型 API 配置。
func CreateAPI(ctx context.Context, input *m.ModelAPIModel) (m.ModelAPIModel, error) {
	if input == nil {
		return m.ModelAPIModel{}, fmt.Errorf("input is required")
	}

	api := *input
	// 新增配置时统一由后端生成主键，忽略前端传入的 id。
	api.ID = uuid.NewString()
	if strings.TrimSpace(api.APIModel) == "" {
		api.APIModel = strings.TrimSpace(api.Model)
	}

	if err := r.CreateAPI(ctx, bootstrap.GormDB, &api); err != nil {
		return m.ModelAPIModel{}, err
	}

	return api, nil
}
