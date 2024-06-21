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
	sync.Mutex
}

func NewTaskPool() *TaskPool {
	return &TaskPool{list: make(map[uint]*Task)}
}

func (p *TaskPool) Get(key uint) *Task {
	p.Lock()
	v, ok := p.list[key]
	p.Unlock()
	if !ok {
		return nil
	}
	return v
}

func (p *TaskPool) Set(key uint, task *Task) {
	p.Lock()
	p.list[key] = task
	p.Unlock()
}

func (p *TaskPool) Delete(key uint) {
	delete(p.list, key)
}
