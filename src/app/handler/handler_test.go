package handler

import (
	"context"
	"github.com/mjedari/health-checker/app/services/mock_files"
	"github.com/mjedari/health-checker/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HandlerTestSetup struct {
	HealthService *mock_files.MockHealthService
	Ctx           context.Context
}

func setupHandlerTest() HandlerTestSetup {
	mockHealthService := new(mock_files.MockHealthService)
	return HandlerTestSetup{
		HealthService: mockHealthService,
		Ctx:           context.Background(),
	}
}

func TestHealthHandler_Index(t *testing.T) {
	// arrange
	setup := setupHandlerTest()

	endpoints := []domain.Endpoint{
		{
			URL:      "http:example.com",
			Method:   "POST",
			Headers:  nil,
			Body:     nil,
			Interval: 0,
		},
	}

	cases := []struct {
		name        string
		mockService func()
		want        int
	}{
		{
			"successful fetch",
			func() {
				setup.HealthService.On("FetchAllEndpoints", setup.Ctx).Return(endpoints, nil)
			},
			http.StatusOK,
		},
		{
			"no endpoint return",
			func() {
				setup.HealthService.On("FetchAllEndpoints", setup.Ctx).Return(nil, nil)
			},
			http.StatusOK,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// here we should test http
			handler := NewHealthHandler(setup.HealthService)
			c.mockService()

			req, err := http.NewRequest("GET", "/endpoint", nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.Index(rr, req)

			assert.Equal(t, c.want, rr.Code)

			setup.HealthService.AssertExpectations(t)

		})
	}
}
