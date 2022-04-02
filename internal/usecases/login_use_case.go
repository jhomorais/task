package usecases

import (
	"context"
	"errors"

	"github.com/fgmaia/task/internal/repositories"
	"github.com/fgmaia/task/internal/usecases/contracts"
	"github.com/fgmaia/task/internal/usecases/ports/output"
	"github.com/fgmaia/task/utils"
)

type loginUseCase struct {
	userRepository repositories.UserRepository
}

func NewLoginUseCase(userRepository repositories.UserRepository) contracts.LoginUseCase {

	return &loginUseCase{
		userRepository: userRepository,
	}
}

func (l *loginUseCase) Execute(ctx context.Context, email string, password string) (*output.LoginOutput, error) {

	user, err := l.userRepository.FindByEmail(context.Background(), email)

	if err != nil {
		return nil, err
	}

	if user == nil || user.ID == "" {
		return nil, errors.New("error when try to find user")
	}

	output := &output.LoginOutput{}

	if utils.DoPasswordsMatch(user.Password, password, []byte(utils.SALT)) {
		output.User = user
	}

	return output, nil
}
