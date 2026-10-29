package ds

import (
	"fmt"
	"sync"
)

type safeMap[K comparable, V any] struct {
	mu   sync.Mutex
	data map[K]V
}

func NewSafeMap[K comparable, V any]() *safeMap[K, V] {
	return &safeMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *safeMap[K, V]) Put(key K, val V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	var counter int
	counter = 1
	fmt.Println("Inside put function: ", counter)
	counter++
	m.data[key] = val
}

func (m *safeMap[K, V]) Get(key K) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	foundKey, ok := m.data[key]
	return foundKey, ok
}
