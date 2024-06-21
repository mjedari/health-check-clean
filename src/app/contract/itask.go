package contract

import (
	"github.com/mjedari/health-checker/domain"
)

type ITaskPool interface {
	Get(key uint) *domain.Task
	Set(key uint, task *domain.Task)
	Delete(key uint)
}
