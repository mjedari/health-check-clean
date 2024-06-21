package domain

import (
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
