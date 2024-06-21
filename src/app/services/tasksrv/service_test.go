package tasksrv

import (
	"context"
	"errors"
	"github.com/mjedari/health-checker/domain"
	"github.com/stretchr/testify/mock"
	"testing"
)

//func TestGetOrCreateTask_TaskExists(t *testing.T) {
//
//	// imagine we have endpoint on mock cache and mock pool
//
//	// arrange
//	mockStorage := new(MockCache)
//	mockCache := new(MockTaskPool)
//	service := NewTaskService(mockStorage, mockCache)
//
//	ctx := context.Background()
//	endpoint := domain.Endpoint{ID: 123, URL: "http://example.com"}
//	task := domain.NewTask(endpoint)
//
//	mockStorage.On("Get", ctx, uint(123), mock.AnythingOfType("*domain.Endpoint")).Run(func(args mock.Arguments) {
//		sumPointer := args.Get(2).(*domain.Endpoint)
//		*sumPointer = endpoint
//	}).Return(nil)
//	mockCache.On("Get", uint(123)).Return(task)
//
//	// act
//	_, _ = service.GetOrCreateTask(ctx, endpoint)
//
//}

func TestGetTask_TaskExists(t *testing.T) {

	// arrange
	mockCache := new(MockCache)
	mockPool := new(MockTaskPool)
	service := NewTaskService(mockCache, mockPool)

	ctx := context.Background()
	endpoint := domain.Endpoint{ID: 123, URL: "http://example.com"}
	task := domain.NewTask(endpoint)

	mockCache.On("Get", ctx, uint(123), mock.AnythingOfType("*domain.Endpoint")).Run(func(args mock.Arguments) {
		sumPointer := args.Get(2).(*domain.Endpoint)
		*sumPointer = endpoint
	}).Return(nil)
	mockPool.On("Get", uint(123)).Return(task)

	// act
	_, _ = service.GetTask(ctx, endpoint)

	mockCache.AssertExpectations(t)
	mockPool.AssertExpectations(t)
}

func TestGetTask_TaskDoesNotExist(t *testing.T) {
	// arrange
	mockCache := new(MockCache)
	mockPool := new(MockTaskPool)
	service := NewTaskService(mockCache, mockPool)

	ctx := context.Background()
	endpoint := domain.Endpoint{ID: 123, URL: "http://example.com"}
	task := domain.NewTask(endpoint)

	mockCache.On("Get", ctx, uint(123), mock.AnythingOfType("*domain.Endpoint")).Run(func(args mock.Arguments) {
		sumPointer := args.Get(2).(*domain.Endpoint)
		*sumPointer = endpoint
	}).Return(nil, errors.New("not found"))
	mockPool.On("Get", uint(123)).Return(task)

	// act
	_, _ = service.GetTask(ctx, endpoint)

	mockCache.AssertExpectations(t)
	mockPool.AssertExpectations(t)
}
