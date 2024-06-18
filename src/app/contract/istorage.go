package contract

import (
	"context"
	"time"
)

type IStorage interface {
	Store(ctx context.Context, key, value string, timeToLive time.Duration) error
	Fetch(ctx context.Context, key string) []byte
	Exists(ctx context.Context, key string) bool
	Delete(ctx context.Context, key string) error
}

type IRepository interface {
	Create(ctx context.Context, value any) error
	Read(ctx context.Context, id uint, out any) error
	ReadAll(ctx context.Context, out any) error
	Update(ctx context.Context, value any) error
	Delete(ctx context.Context, value any) error
}
