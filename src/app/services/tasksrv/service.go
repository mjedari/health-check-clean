package tasksrv

import (
	"context"
	"errors"
	"github.com/mjedari/health-checker/app/contract"
	"github.com/mjedari/health-checker/domain"
)

type TaskService struct {
	storage contract.IStorage
	cache   contract.ITaskCache
}

func NewTaskService(storage contract.IStorage, cache contract.ITaskCache) *TaskService {
	return &TaskService{storage: storage, cache: cache}
}

func (s *TaskService) GetOrCreateTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {

	existedTask, err := s.GetTask(ctx, endpoint)
	if existedTask != nil && err != nil {
		return existedTask, nil
	}

	task := domain.NewTask(endpoint)
	// store end point
	err = s.storage.Add(ctx, endpoint.ID, endpoint)
	if err != nil {
		return nil, err
	}

	// store task in memory
	s.cache.Set(endpoint.ID, task)
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {
	var item domain.Endpoint
	err := s.storage.Get(ctx, endpoint.ID, &item)
	if err != nil {
		return nil, err
	}

	v := s.cache.Get(item.ID)
	if v == nil {
		return nil, errors.New("task not found in task service cache")
	}

	return v, nil

}

func (s *TaskService) RemoveTask(ctx context.Context, endpoint domain.Endpoint) error {
	task, err := s.GetTask(ctx, endpoint)
	if err != nil {
		return errors.New("task not found")
	}
	// stop task
	task.Stop <- true

	// delete from pool
	err = s.storage.Remove(ctx, endpoint.ID)
	if err != nil {
		return err
	}

	// delete from cache
	s.cache.Delete(endpoint.ID)

	return nil
}
