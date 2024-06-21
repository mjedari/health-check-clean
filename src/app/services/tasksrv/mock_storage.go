package tasksrv

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) Add(ctx context.Context, key uint, item any) error {
	args := m.Called(ctx, key, item)
	return args.Error(0)
}

func (m *MockCache) Exist(ctx context.Context, key uint) bool {
	//TODO implement me
	panic("implement me")
}

func (m *MockCache) Get(ctx context.Context, key uint, out any) error {
	args := m.Called(ctx, key, out)
	return args.Error(0)
}

func (m *MockCache) Remove(ctx context.Context, key uint) error {
	//TODO implement me
	panic("implement me")
}
