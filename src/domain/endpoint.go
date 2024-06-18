package domain

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/url"
	"slices"
	"strings"
	"time"
)

type Endpoint struct {
	ID        uint          `gorm:"primaryKey" json:"id"`
	URL       string        `json:"url"`
	Method    string        `json:"method"`
	Headers   Headers       `json:"headers"`
	Body      Body          `json:"body"`
	Interval  time.Duration `json:"interval"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type ValidationErrors []error

func (v ValidationErrors) Error() string {
	return ""
}

func (v ValidationErrors) Converted() []string {
	var m []string
	for _, err := range v {
		m = append(m, err.Error())
	}
	return m
}

func (e *Endpoint) Validate() ValidationErrors {
	var errs ValidationErrors
	// validate url
	if e.URL == "" {
		errs = append(errs, errors.New("url is required"))
	}

	if _, err := url.Parse(e.URL); err != nil {
		errs = append(errs, errors.New("url is not valid"))
	}

	// validate Method
	if e.Method == "" {
		errs = append(errs, errors.New("method is required"))
	}
	found := slices.Contains([]string{"post", "get"}, strings.ToLower(e.Method))
	if !found {
		errs = append(errs, errors.New("method is invalid"))
	}

	if e.Interval == 0 {
		errs = append(errs, errors.New("internal is required"))
	}

	if e.Interval < 0 {
		errs = append(errs, errors.New("internal is invalid"))
	}

	// todo: validate headers, body

	return errs
}

type Headers map[string]string

func (b *Headers) Scan(value any) error {
	var headers Headers
	err := json.Unmarshal(value.([]byte), &headers)
	*b = headers
	return err
}

func (h Headers) Value() (driver.Value, error) {
	b, err := json.Marshal(h)
	return string(b), err
}

type Body map[string]any

func (b *Body) Scan(value any) error {
	var body Body
	err := json.Unmarshal(value.([]byte), &body)
	*b = body
	return err
}

func (h Body) Value() (driver.Value, error) {
	b, err := json.Marshal(h)
	return string(b), err
}
