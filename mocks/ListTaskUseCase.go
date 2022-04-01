// Code generated by mockery v2.10.1. DO NOT EDIT.

package mocks

import (
	context "context"

	input "github.com/fgmaia/task/internal/usecases/ports/input"

	mock "github.com/stretchr/testify/mock"

	output "github.com/fgmaia/task/internal/usecases/ports/output"
)

// ListTaskUseCase is an autogenerated mock type for the ListTaskUseCase type
type ListTaskUseCase struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, findTask
func (_m *ListTaskUseCase) Execute(ctx context.Context, findTask *input.ListTaskInput) (*output.ListTaskOutput, error) {
	ret := _m.Called(ctx, findTask)

	var r0 *output.ListTaskOutput
	if rf, ok := ret.Get(0).(func(context.Context, *input.ListTaskInput) *output.ListTaskOutput); ok {
		r0 = rf(ctx, findTask)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*output.ListTaskOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *input.ListTaskInput) error); ok {
		r1 = rf(ctx, findTask)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
