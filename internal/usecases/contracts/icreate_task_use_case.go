package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/usecases/ports/input"
	"github.com/fgmaia/task/internal/usecases/ports/output"
)

type CreateTaskUseCase interface {
	Execute(ctx context.Context, createTask *input.CreateTaskInput) (*output.CreateTaskOutput, error)
}
