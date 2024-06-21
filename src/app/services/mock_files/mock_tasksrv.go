package mock_files

import (
	"context"
	"github.com/mjedari/health-checker/domain"
	"github.com/stretchr/testify/mock"
)

type MockTaskService struct {
	mock.Mock
}

func (s *MockTaskService) GetOrCreateTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MockTaskService) CreateTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MockTaskService) GetTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MockTaskService) RemoveTask(ctx context.Context, endpoint domain.Endpoint) error {
	return s.Called(ctx, endpoint).Error(0)
}
