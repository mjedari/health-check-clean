package contract

import (
	"context"
	"github.com/mjedari/health-checker/domain"
)

type IHealthService interface {
	FetchAllEndpoints(ctx context.Context) ([]domain.Endpoint, error)
	CreateEndpoint(ctx context.Context, endpoint domain.Endpoint) error
	DeleteEndpoint(ctx context.Context, id uint) error
	FetchEndpoint(ctx context.Context, id uint) (domain.Endpoint, error)
	StartWatching(ctx context.Context, endpoint domain.Endpoint) error
	StopWatching(ctx context.Context, endpoint domain.Endpoint) error
}

type ITaskService interface {
	GetOrCreateTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error)
	CreateTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error)
	GetTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error)
	RemoveTask(ctx context.Context, endpoint domain.Endpoint) error
}
