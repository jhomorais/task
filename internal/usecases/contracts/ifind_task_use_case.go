package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/usecases/ports/input"
	"github.com/fgmaia/task/internal/usecases/ports/output"
)

type FindTaskUseCase interface {
	Execute(ctx context.Context, createTask *input.FindTaskInput) (*output.FindTaskOutput, error)
}
