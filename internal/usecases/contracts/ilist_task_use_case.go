package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/usecases/ports/output"
)

type ListTaskUseCase interface {
	Execute(ctx context.Context, userID string) (*output.ListTaskOutput, error)
}
