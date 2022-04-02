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

func NewTaskPb(userId string) *taskpb.Task {

	task := &taskpb.Task{
		Id:          RandomID(),
		Summary:     RandomSummary(),
		PerformedAt: timestamppb.New(time.Now()),
		UserId:      userId,
	}

	return task
}
