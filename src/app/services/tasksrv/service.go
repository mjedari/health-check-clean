package tasksrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/mjedari/health-checker/app/contract"
	"github.com/mjedari/health-checker/domain"
)

type TaskService struct {
	cache contract.ICache
	pool  contract.ITaskPool
}

func NewTaskService(storage contract.ICache, pool contract.ITaskPool) *TaskService {
	return &TaskService{cache: storage, pool: pool}
}

func (s *TaskService) GetOrCreateTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {

	existedTask, err := s.GetTask(ctx, endpoint)
	fmt.Println("existedTask", existedTask)
	if existedTask != nil && err == nil {
		return existedTask, nil
	}

	task := domain.NewTask(endpoint)
	// store end point
	err = s.cache.Add(ctx, endpoint.ID, endpoint)
	if err != nil {
		return nil, err
	}

	// store task in memory
	s.pool.Set(endpoint.ID, task)
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {
	var item domain.Endpoint
	err := s.cache.Get(ctx, endpoint.ID, &item)
	if err != nil {
		return nil, err
	}

	v := s.pool.Get(item.ID)
	if v == nil {
		return nil, errors.New("task not found in task service pool")
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
	err = s.cache.Remove(ctx, endpoint.ID)
	if err != nil {
		return err
	}

	// delete from pool
	s.pool.Delete(endpoint.ID)

	return nil
}

func (s *TaskService) CreateTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {
	task := domain.NewTask(endpoint)
	err := s.cache.Add(ctx, endpoint.ID, endpoint)
	if err != nil {
		return nil, err
	}

	// store task in memory
	s.pool.Set(endpoint.ID, task)
	return task, nil
}
