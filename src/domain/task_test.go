package domain

import (
	"testing"
)

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
