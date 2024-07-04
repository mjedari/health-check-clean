package domain

import (
	"net/http"
	"sync"
	"time"
)

/*
Task
*/

type TaskState int

const (
	INIT TaskState = iota
	RUNNING
	STOPPED
)

type Task struct {
	Endpoint   Endpoint
	State      TaskState
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
	return &Task{State: INIT, Endpoint: endpoint, Stop: make(chan bool)}
}

func (t *Task) TaskStart() {
	t.State = RUNNING
}

func (t *Task) TaskStop() {
	t.State = STOPPED
}

func (t *Task) IsRunning() bool {
	return t.State == RUNNING
}

/*
Task Pool
*/

type TaskPool struct {
	list map[uint]*Task
	mu   sync.RWMutex
}

func NewTaskPool() *TaskPool {
	return &TaskPool{list: make(map[uint]*Task)}
}

func (p *TaskPool) Get(key uint) *Task {
	p.mu.Lock()
	defer p.mu.Unlock()
	v, ok := p.list[key]
	if !ok {
		return nil
	}
	return v
}

func (p *TaskPool) Set(key uint, task *Task) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.list[key] = task
}

func (p *TaskPool) Delete(key uint) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.list, key)
}

func (p *TaskPool) Length() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.list)
}
