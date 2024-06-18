package tasksrv

import (
	"context"
	"errors"
	"github.com/mjedari/health-checker/app/contract"
	"github.com/mjedari/health-checker/domain"
)

type TaskService struct {
	storage contract.ITaskPool
	cache   map[uint]*domain.Task
	// todo: mutex for race condition
}

func NewTaskService(storage contract.ITaskPool) *TaskService {
	return &TaskService{storage: storage, cache: make(map[uint]*domain.Task)}
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
	s.cache[endpoint.ID] = task
	return task, nil
}

func (s *TaskService) GetTask(ctx context.Context, endpoint domain.Endpoint) (*domain.Task, error) {
	var item domain.Endpoint
	err := s.storage.Get(ctx, endpoint.ID, &item)
	if err != nil {
		return nil, err
	}

	v, ok := s.cache[item.ID]
	if !ok {
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
	delete(s.cache, endpoint.ID)

	return nil
}
