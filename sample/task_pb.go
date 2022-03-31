package sample

import (
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fgmaia/task/pb/taskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	gofakeit.Seed(time.Now().UnixNano())
}

func NewTaskPb() *taskpb.Task {

	task := &taskpb.Task{
		Id:          randomID(),
		Summary:     randomSummary(),
		PerformedAt: timestamppb.New(gofakeit.Date()),
		UserId:      randomID(),
	}

	return task
}
