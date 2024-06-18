package healthsrv

import (
	"context"
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"github.com/mjedari/health-checker/app/contract"
	"github.com/mjedari/health-checker/domain"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"time"
)

type HealthService struct {
	repo        contract.IRepository
	taskService contract.ITaskService
	config      config.Webhook
}

func NewHealthService(repo contract.IRepository, taskService contract.ITaskService, webhook config.Webhook) *HealthService {
	return &HealthService{repo: repo, taskService: taskService, config: webhook}
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
				// call http
				// retry if needed
				// get the response
				// decide to call webhook or not
				// update database if needed
				res := mockCall(endpoint)
				fmt.Println(task.LastStatus, res.status)
				if task.IsStatusChanged(res.status) {
					logrus.Infof("endpoint %s status changed to code: %d", endpoint.URL, res.status)
					mockWebhookCall(endpoint)
				}

				task.UpdateStatus(res.status)
			}
		}
	}()

	return nil

}

func mockCall(endpoint domain.Endpoint) MockedResponse {
	fmt.Println("called endpoint number: ", endpoint.ID)
	successResponse := MockedResponse{status: http.StatusOK}
	failedResponse := MockedResponse{status: http.StatusInternalServerError}
	r := rand.Intn(100)
	if r%2 == 0 {
		return successResponse
	}
	return failedResponse
}

func mockWebhookCall(endpoint domain.Endpoint) MockedResponse {
	fmt.Println("webhook called for endpoint number: ", endpoint.ID)
	successResponse := MockedResponse{status: http.StatusOK}
	return successResponse
}

type MockedResponse struct {
	status int
}

func (s *HealthService) StopWatching(ctx context.Context, endpoint domain.Endpoint) error {
	err := s.taskService.RemoveTask(ctx, endpoint)
	if err != nil {
		return err
	}

	return nil
}
