// Code generated by mockery v2.10.1. DO NOT EDIT.

package mocks

import (
	context "context"

	entities "github.com/fgmaia/task/internal/domain/entities"

	mock "github.com/stretchr/testify/mock"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, entity
func (_m *UserRepository) CreateUser(ctx context.Context, entity *entities.User) error {
	ret := _m.Called(ctx, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.User) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
