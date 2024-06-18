package domain

import (
	"net/http"
	"time"
)

// each task should be assigned to a go routine

// to make this scalable we need to store date on redis to access fast to be cloud native
//so we need a service as a TaskService

type Task struct {
	Endpoint   Endpoint
	Stop       chan bool
	LastStatus int
	LastUpdate time.Time
}

func (t *Task) IsStatusChanged(status int) bool {
	if t.LastStatus == 0 && status != http.StatusOK {
		return true
	}

	return t.LastStatus != 0 && t.LastStatus != status
}

func (t *Task) UpdateStatus(status int) {
	t.LastStatus = status
	t.LastUpdate = time.Now()
}

func NewTask(endpoint Endpoint) *Task {
	return &Task{Endpoint: endpoint, Stop: make(chan bool)}
}
