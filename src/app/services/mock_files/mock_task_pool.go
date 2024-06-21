package mock_files

import (
	"github.com/mjedari/health-checker/domain"
	"github.com/stretchr/testify/mock"
)

type MockTaskPool struct {
	mock.Mock
}

func (m *MockTaskPool) Get(key uint) *domain.Task {
	args := m.Called(key)
	return args.Get(0).(*domain.Task)
}

func (m *MockTaskPool) Set(key uint, task *domain.Task) {
	m.Called(key, task)
}

func (m *MockTaskPool) Delete(key uint) {
	//TODO implement me
	panic("implement me")
}
