package mock_files

import (
	"context"
	"github.com/mjedari/health-checker/domain"
	"github.com/stretchr/testify/mock"
)

type MockHealthService struct {
	mock.Mock
}

func (s *MockHealthService) FetchAllEndpoints(ctx context.Context) ([]domain.Endpoint, error) {
	args := s.Called(ctx)
	return args.Get(0).([]domain.Endpoint), args.Error(1)
}

func (s *MockHealthService) CreateEndpoint(ctx context.Context, endpoint *domain.Endpoint) error {
	//TODO implement me
	panic("implement me")
}

func (s *MockHealthService) DeleteEndpoint(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}

func (s *MockHealthService) FetchEndpoint(ctx context.Context, id uint) (domain.Endpoint, error) {
	//TODO implement me
	panic("implement me")
}

func (s *MockHealthService) StartWatching(ctx context.Context, endpoint domain.Endpoint) error {
	//TODO implement me
	panic("implement me")
}

func (s *MockHealthService) StopWatching(ctx context.Context, endpoint domain.Endpoint) error {
	//TODO implement me
	panic("implement me")
}
