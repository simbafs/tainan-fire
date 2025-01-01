package bucket

import (
	"time"
)

type BucketItem[T any] struct {
	updateTime time.Time
	data       T
}

type PairItem[T any] struct {
	HasOld bool
	Old    T
	New    T
}

type Bucket[T any] struct {
	data      map[string]BucketItem[T]
	compare   func(T, T) bool
	aliveTime time.Duration
}

func New[T any](aliveTime time.Duration, compare func(T, T) bool) *Bucket[T] {
	return &Bucket[T]{
		data:      map[string]BucketItem[T]{},
		compare:   compare,
		aliveTime: aliveTime,
	}
}

func (b *Bucket[T]) Set(id string, item T) {
	b.data[id] = BucketItem[T]{updateTime: time.Now(), data: item}
}

func (b *Bucket[T]) Get(id string) (T, bool) {
	item, ok := b.data[id]
	return item.data, ok
}

// GC delete the data older than
func (b *Bucket[T]) GC() {
	for k, v := range b.data {
		if time.Since(v.updateTime) > b.aliveTime {
			delete(b.data, k)
		}
	}
}

func (b *Bucket[T]) Len() int {
	return len(b.data)
}
