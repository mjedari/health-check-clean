package domain

import (
	"errors"
	"reflect"
	"testing"
)

func TestEndpoint_Validate(t *testing.T) {
	// arrange
	cases := []struct {
		name     string
		endpoint Endpoint
		wantErrs ValidationErrors
	}{
		{
			name:     "valid input",
			endpoint: Endpoint{URL: "http://example.com", Method: "GET", Interval: 10},
			wantErrs: nil,
		},
		{
			name:     "missing URL",
			endpoint: Endpoint{URL: "", Method: "GET", Interval: 10},
			wantErrs: ValidationErrors{errors.New("url is required")},
		},
		{
			name:     "invalid URL",
			endpoint: Endpoint{URL: "://invalid-url", Method: "GET", Interval: 10},
			wantErrs: ValidationErrors{errors.New("url is not valid")},
		},
		{
			name:     "missing method",
			endpoint: Endpoint{URL: "http://example.com", Method: "", Interval: 10},
			wantErrs: ValidationErrors{errors.New("method is required"), errors.New("method is invalid")},
		},
		{
			name:     "invalid method",
			endpoint: Endpoint{URL: "http://example.com", Method: "PUT", Interval: 10},
			wantErrs: ValidationErrors{errors.New("method is invalid")},
		},
		{
			name:     "missing interval",
			endpoint: Endpoint{URL: "http://example.com", Method: "GET", Interval: 0},
			wantErrs: ValidationErrors{errors.New("interval is required")},
		},
		{
			name:     "invalid interval",
			endpoint: Endpoint{URL: "http://example.com", Method: "GET", Interval: -1},
			wantErrs: ValidationErrors{errors.New("interval is invalid")},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			errs := c.endpoint.Validate()

			if !reflect.DeepEqual(errs, c.wantErrs) {
				t.Errorf("want %v got %v", c.wantErrs, errs)
			}
		})
	}
}
