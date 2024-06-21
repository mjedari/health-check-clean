package healthsrv

import (
	"context"
	"errors"
	"github.com/mjedari/health-checker/app/contract"
	"github.com/mjedari/health-checker/domain"
	"github.com/sirupsen/logrus"
	"time"
)

type HealthService struct {
	client      contract.IClient
	taskService contract.ITaskService
	repo        contract.IRepository
}

func NewHealthService(client contract.IClient, repo contract.IRepository, taskService contract.ITaskService) *HealthService {
	return &HealthService{client: client, repo: repo, taskService: taskService}
}

func (s *HealthService) FetchAllEndpoints(ctx context.Context) ([]domain.Endpoint, error) {
	var endpoints []domain.Endpoint
	s.repo.ReadAll(ctx, &endpoints)

	return endpoints, nil
}

func (s *HealthService) CreateEndpoint(ctx context.Context, endpoint *domain.Endpoint) error {
	err := s.repo.Create(ctx, endpoint)
	return err
}

func (s *HealthService) DeleteEndpoint(ctx context.Context, id uint) error {

	endpoint := domain.Endpoint{ID: id}
	err := s.taskService.RemoveTask(ctx, endpoint)
	if err != nil {
		return err
	}

	// delete from database
	err = s.repo.Delete(ctx, endpoint)
	if err != nil {
		return err
	}

	return err
}

func (s *HealthService) FetchEndpoint(ctx context.Context, id uint) (domain.Endpoint, error) {
	var endpoint domain.Endpoint
	err := s.repo.Read(ctx, id, &endpoint)
	return endpoint, err
}

func (s *HealthService) StartWatching(ctx context.Context, endpoint domain.Endpoint) error {

	task, err := s.taskService.GetOrCreateTask(ctx, endpoint)
	if err != nil {
		return err
	}

	if task.IsRunning() {
		return errors.New("task has been started already")
	}

	task.TaskStart()

	go func() {
		ticker := time.NewTicker(endpoint.Interval * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-task.Stop:
				return
			case <-ctx.Done():
				return
			case <-ticker.C:
				resStatus := s.client.HttpCall(endpoint)
				if task.IsStatusChanged(resStatus) {
					logrus.Infof("endpoint %s status changed to code: %d", endpoint.URL, resStatus)
					go s.client.HttpWebhookCall(endpoint, resStatus)
				}

				task.UpdateStatus(resStatus)
			}
		}
	}()

	logrus.Infof("watching started for %s", endpoint.URL)
	return nil

}

func (s *HealthService) StopWatching(ctx context.Context, endpoint domain.Endpoint) error {
	err := s.taskService.RemoveTask(ctx, endpoint)
	if err != nil {
		return err
	}

	logrus.Infof("watching stopped for %s", endpoint.URL)
	return nil
}
