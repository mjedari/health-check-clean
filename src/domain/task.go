package domain

import (
	"net/http"
	"sync"
	"time"
)

/*
Task
*/
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

/*
Task Cache
*/
type TaskCache struct {
	list map[uint]*Task
	sync.Mutex
}

func NewTaskCache() *TaskCache {
	return &TaskCache{list: make(map[uint]*Task)}
}

func (p *TaskCache) Get(key uint) *Task {
	p.Lock()
	v, ok := p.list[key]
	p.Unlock()
	if !ok {
		return nil
	}
	return v
}

func (p *TaskCache) Set(key uint, task *Task) {
	p.Lock()
	p.list[key] = task
	p.Unlock()
}

func (p *TaskCache) Delete(key uint) {
	delete(p.list, key)
}
