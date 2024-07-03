package poolx

import (
	"sync"
)

type Pool[T any] interface {
	Get() T
	Put(t T)
}

var _ Pool[any] = (*PoolNormal[any])(nil)

type PoolNormal[T any] struct {
	pool sync.Pool
}

func NewPoolNormal[T any](fn func() T) Pool[T] {
	fun := func() any {
		return fn()
	}
	p := &PoolNormal[T]{
		pool: sync.Pool{New: fun},
	}
	return p
}

func (p *PoolNormal[T]) Get() T {
	return p.pool.Get()
}

func (p *PoolNormal[T]) Put(t T) {
	p.pool.Put(t)
}
