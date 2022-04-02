package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/usecases/ports/output"
)

type FindTaskUseCase interface {
	Execute(ctx context.Context, userID string, taskID string) (*output.FindTaskOutput, error)
}
