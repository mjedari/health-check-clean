package domain

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestTask_IsRunning(t *testing.T) {
	endpoint := Endpoint{}

	cases := []struct {
		name  string
		state TaskState
		want  bool
	}{
		{"task is running", RUNNING, true},
		{"task is in init", INIT, false},
		{"task is in stop state", STOPPED, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			newTask := NewTask(endpoint)
			newTask.State = c.state

			got := newTask.IsRunning()
			if got != c.want {
				t.Errorf("want %v but got %v", c.want, got)
			}
		})
	}
}

func TestTask_TaskChangeState(t *testing.T) {
	endpoint := Endpoint{}
	newTask := NewTask(endpoint)

	cases := []struct {
		name  string
		state TaskState
	}{
		{"change state to running", RUNNING},
		{"change state to stopped", STOPPED},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			switch c.state {
			case RUNNING:
				newTask.TaskStart()
			case STOPPED:
				newTask.TaskStop()
			}

			if newTask.State != c.state {
				t.Errorf("want %v but got %v", c.state, newTask.State)
			}
		})
	}

}

func TestTask_IsStatusChanged(t *testing.T) {
	// arrange
	endpoint := Endpoint{}

	cases := []struct {
		name      string
		status    int
		newStatus int
		want      bool
	}{
		{
			"status changed",
			200,
			500,
			true,
		},
		{
			"status not changed",
			500,
			500,
			false,
		},
		{
			"status zero",
			0,
			200,
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			task := NewTask(endpoint)
			task.UpdateStatus(c.status)
			changed := task.IsStatusChanged(c.newStatus)
			if changed != c.want {
				t.Errorf("want %v got %v", c.want, changed)
			}
		})
	}

}

/*
	Task Pool
*/

func TestTaskPool_Delete(t *testing.T) {
	// arrange
	task := NewTask(Endpoint{})
	taskPool := NewTaskPool()
	key := uint(1)

	taskPool.Set(key, task)

	// act
	taskPool.Delete(key)
	got := taskPool.Get(key)

	// assert
	if !assert.Nil(t, got) {
		t.Log("This is: ", got)
	}
}

func TestTaskPool_Get(t *testing.T) {
	//Get(key uint) *Task
	// arrange
	newTask := NewTask(Endpoint{})
	taskPool := NewTaskPool()
	for i := 0; i < 1000; i++ {
		taskPool.Set(uint(i), newTask)
	}

	// act
	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(key uint) {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				_ = taskPool.Get(key*1000 + uint(j))
			}
		}(uint(1))
	}

	wg.Wait()
}

//func TestTaskPoolConcurrency(t *testing.T) {
//	pool := NewTaskPool()
//	var wg sync.WaitGroup
//
//	numGoroutines := 100
//	numTasks := 1
//
//	var count int64
//	setFunction := func(id int) {
//		defer wg.Done()
//		for i := 1; i <= numTasks; i++ {
//			atomic.AddInt64(&count, 1)
//			task := NewTask(Endpoint{ID: uint(i)})
//			pool.Set(uint(count), task)
//		}
//	}
//
//	testFunc := func(id int) {
//		defer wg.Done()
//		for i := 1; i <= numTasks; i++ {
//			atomic.AddInt64(&count, 1)
//			task := NewTask(Endpoint{ID: uint(i)})
//			pool.Set(uint(count), task)
//		}
//	}
//
//	wg.Add(numGoroutines)
//	for i := 1; i <= numGoroutines; i++ {
//		go setFunction(i)
//	}
//
//	wg.Wait()
//	t.Log("HERE ", count)
//	got := pool.Length()
//	expected := numGoroutines * numTasks
//	if got != expected {
//		t.Errorf("Length expected %d but got %d", expected, got)
//	}
//}
