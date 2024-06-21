package httpsrv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"github.com/mjedari/health-checker/domain"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
)

type HttpService struct {
	client http.Client
	config config.Webhook
}

func NewHttpService(config config.Webhook) *HttpService {
	return &HttpService{config: config, client: http.Client{Timeout: config.Timeout}}
}

func (h *HttpService) HttpCall(endpoint domain.Endpoint) int {
	// todo: use retry pattern

	var endpointBody []byte
	if endpoint.Body != nil {
		byteBody, _ := json.Marshal(endpoint.Body)
		endpointBody = byteBody
	}
	req, err := http.NewRequest(endpoint.Method, endpoint.URL, bytes.NewBuffer(endpointBody))
	for key, value := range endpoint.Headers {
		req.Header.Set(key, value)
	}

	resp, err := h.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	return resp.StatusCode
}

func (h *HttpService) HttpWebhookCall(endpoint domain.Endpoint, status int) {
	logrus.Warnf("webhook called for endpoint number: %d \n", endpoint.ID)
	message := fmt.Sprintf("%s change status to %d", endpoint.URL, status)
	req, _ := http.NewRequest(h.config.Method, h.config.URL, bytes.NewBufferString(message))

	_, err := h.client.Do(req)
	if err != nil {
		logrus.Warn("error while calling webhook: ", err)
	}
}
