package mock_files

import (
	"github.com/mjedari/health-checker/domain"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) HttpCall(endpoint domain.Endpoint) int {
	//TODO implement me
	panic("implement me")
}

func (m *MockClient) HttpWebhookCall(endpoint domain.Endpoint, status int) {
	//TODO implement me
	panic("implement me")
}
