package contracts

import (
	"context"

	"github.com/fgmaia/task/internal/usecases/ports/output"
)

type LoginUseCase interface {
	Execute(ctx context.Context, email string, password string) (*output.LoginOutput, error)
}
