package tasksrv

import (
	"github.com/mjedari/health-checker/domain"
	"github.com/stretchr/testify/mock"
)

type MockTaskCache struct {
	mock.Mock
}

func (m *MockTaskCache) Get(key uint) *domain.Task {
	args := m.Called(key)
	return args.Get(0).(*domain.Task)
}

func (m *MockTaskCache) Set(key uint, task *domain.Task) {
	m.Called(key, task)
}

func (m *MockTaskCache) Delete(key uint) {
	//TODO implement me
	panic("implement me")
}
