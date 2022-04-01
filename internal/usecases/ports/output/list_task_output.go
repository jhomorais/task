package output

import "github.com/fgmaia/task/internal/domain/entities"

type ListTaskOutput struct {
	Tasks []*entities.Task
}
