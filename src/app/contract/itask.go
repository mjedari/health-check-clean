package contract

import "context"

type ITaskPool interface {
	Add(ctx context.Context, key uint, item any) error
	Exist(ctx context.Context, key uint) bool
	Get(ctx context.Context, key uint, out any) error
	Remove(ctx context.Context, key uint) error
}
