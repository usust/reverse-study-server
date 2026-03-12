package config

import (
	"context"
	"fmt"
	"reverse-study-server/internal/bootstrap"
	promptrepo "reverse-study-server/internal/repository/prompt"
)

// DeletePromptByID 按 ID 删除提示词。
func DeletePromptByID(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("id is required")
	}
	return promptrepo.DeleteByID(ctx, bootstrap.GormDB, id)
}
