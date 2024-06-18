package domain

import "time"

type Webhook struct {
	URL     string
	Method  string
	Timeout time.Duration
}
