package mock_files

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (r *MockRepository) Create(ctx context.Context, value any) error {
	//TODO implement me
	panic("implement me")
}

func (r *MockRepository) Read(ctx context.Context, id uint, out any) error {
	//TODO implement me
	panic("implement me")
}

func (r *MockRepository) ReadAll(ctx context.Context, out any) error {
	return r.Called(ctx, out).Error(0)
}

func (r *MockRepository) Update(ctx context.Context, value any) error {
	//TODO implement me
	panic("implement me")
}

func (r *MockRepository) Delete(ctx context.Context, value any) error {
	return r.Called(ctx, value).Error(0)
}
