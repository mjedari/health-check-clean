package healthsrv

import (
	"context"
	"github.com/mjedari/health-checker/app/contract"
	"github.com/mjedari/health-checker/domain"
)

type EndpointRepository struct {
	repo contract.IRepository
}

func NewEndpointRepository(repo contract.IRepository) *EndpointRepository {
	return &EndpointRepository{repo: repo}
}

func (r *EndpointRepository) GetALLEndpoints(ctx context.Context) ([]domain.Endpoint, error) {
	var endpoints []domain.Endpoint
	r.repo.ReadAll(ctx, &endpoints)

	return endpoints, nil
}

func (r *EndpointRepository) CreateEndpoint(ctx context.Context, endpoint *domain.Endpoint) error {
	return r.repo.Create(ctx, endpoint)
}

func (r *EndpointRepository) DeleteEndpoint(ctx context.Context, endpoint *domain.Endpoint) error {
	return r.repo.Delete(ctx, endpoint)
}

func (r *EndpointRepository) FetchEndpoint(ctx context.Context, id uint) (domain.Endpoint, error) {
	var endpoint domain.Endpoint
	err := r.repo.Read(ctx, id, &endpoint)
	return endpoint, err
}
