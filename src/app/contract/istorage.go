package contract

import (
	"context"
)

type ICache interface {
	Add(ctx context.Context, key uint, item any) error
	Exist(ctx context.Context, key uint) bool
	Get(ctx context.Context, key uint, out any) error
	Remove(ctx context.Context, key uint) error
}

type IRepository interface {
	Create(ctx context.Context, value any) error
	Read(ctx context.Context, id uint, out any) error
	ReadAll(ctx context.Context, out any) error
	Update(ctx context.Context, value any) error
	Delete(ctx context.Context, value any) error
}
