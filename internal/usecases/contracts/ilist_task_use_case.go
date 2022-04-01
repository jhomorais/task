package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/usecases/ports/input"
	"github.com/fgmaia/task/internal/usecases/ports/output"
)

type ListTaskUseCase interface {
	Execute(ctx context.Context, findTask *input.ListTaskInput) (*output.ListTaskOutput, error)
}
