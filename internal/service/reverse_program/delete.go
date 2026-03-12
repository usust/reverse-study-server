package reverse_program

import (
	"context"
	"reverse-study-server/internal/bootstrap"
	reverseprogramrepo "reverse-study-server/internal/repository/reverse_program"
)

// DeleteProgramByID 按 ID 删除题目。
func DeleteProgramByID(ctx context.Context, id uint64) error {
	return reverseprogramrepo.DeleteByID(ctx, bootstrap.GormDB, id)
}

