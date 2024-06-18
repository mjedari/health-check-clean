package storage

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

type InMemory struct {
	sync.Mutex
	list map[uint]any
}

func NewInMemory() *InMemory {
	return &InMemory{list: make(map[uint]any)}
}

func (m *InMemory) Add(ctx context.Context, key uint, item any) error {
	m.list[key] = item
	return nil
}

func (m *InMemory) Get(ctx context.Context, key uint, out any) error {
	v, ok := m.list[key]
	if !ok {
		return errors.New("not found")
	}

	rv := reflect.ValueOf(out)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("pointer needed")
	}

	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(v))

	return nil
}

func (m *InMemory) Remove(ctx context.Context, key uint) error {
	_, ok := m.list[key]
	if !ok {
		return errors.New("not found")
	}

	delete(m.list, key)
	return nil
}

func (m *InMemory) Exist(ctx context.Context, key uint) bool {
	_, ok := m.list[key]
	return ok
}
