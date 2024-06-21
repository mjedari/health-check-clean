package healthsrv

import (
	"context"
	"github.com/mjedari/health-checker/app/services/mock_files"
	"github.com/mjedari/health-checker/domain"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type HealthServiceTestSetup struct {
	MockClient      *mock_files.MockClient
	MockRepo        *mock_files.MockRepository
	MockTaskService *mock_files.MockTaskService
	HealthService   *HealthService
	Ctx             context.Context
}

func setupHealthServiceTest() HealthServiceTestSetup {
	mockClient := new(mock_files.MockClient)
	mockRepo := new(mock_files.MockRepository)
	mockTaskService := new(mock_files.MockTaskService)
	healthService := NewHealthService(mockClient, mockRepo, mockTaskService)
	ctx := context.Background()
	return HealthServiceTestSetup{
		MockClient:      mockClient,
		MockRepo:        mockRepo,
		MockTaskService: mockTaskService,
		HealthService:   healthService,
		Ctx:             ctx,
	}
}

func TestHealthService_FetchAllEndpoints(t *testing.T) {
	// arrange
	setup := setupHealthServiceTest()

	endpoints := []domain.Endpoint{
		{
			URL:      "http://example.com",
			Method:   "get",
			Headers:  map[string]string{"test": "test"},
			Body:     nil,
			Interval: 3,
		},
	}

	setup.MockRepo.On("ReadAll", setup.Ctx, mock.AnythingOfType("*[]domain.Endpoint")).Run(func(args mock.Arguments) {
		out := args.Get(1).(*[]domain.Endpoint)
		*out = endpoints
	}).Return(nil)

	// act
	fetchedData, _ := setup.HealthService.FetchAllEndpoints(setup.Ctx)

	// assert
	setup.MockRepo.AssertExpectations(t)
	if !reflect.DeepEqual(endpoints, fetchedData) {
		t.Error("fetched data do not match the endpoints")
	}
}

func TestHealthService_DeleteEndpoint(t *testing.T) {
	// arrange
	setup := setupHealthServiceTest()

	endpoint := domain.Endpoint{
		ID:       12,
		URL:      "http://example.com",
		Method:   "get",
		Headers:  map[string]string{"test": "test"},
		Body:     nil,
		Interval: 3,
	}

	setup.MockTaskService.On("RemoveTask", setup.Ctx, mock.AnythingOfType("domain.Endpoint")).Return(nil)
	setup.MockRepo.On("Delete", setup.Ctx, mock.AnythingOfType("domain.Endpoint")).Return(nil)

	// act
	err := setup.HealthService.DeleteEndpoint(setup.Ctx, endpoint.ID)

	// assert
	setup.MockRepo.AssertExpectations(t)
	if err != nil {
		t.Error("want no error but got", err)
	}
}
