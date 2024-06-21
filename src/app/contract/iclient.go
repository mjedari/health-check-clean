package contract

import "github.com/mjedari/health-checker/domain"

type IClient interface {
	HttpCall(endpoint domain.Endpoint) int
	HttpWebhookCall(endpoint domain.Endpoint, status int)
}
