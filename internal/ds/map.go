package ds

import (
	"fmt"
	"sync"
)

type safeMap[K comparable, V any] struct {
	mu   sync.RWMutex
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

	m.data[key] = val
	fmt.Println("Inside put function: ", m.data)
}

func (m *safeMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	foundKey, ok := m.data[key]
	return foundKey, ok
}

func (m *safeMap[K, V]) Delete(key K) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.data[key]; !ok {
		return false
	} else {
		delete(m.data, key)
		return true
	}
}

func (m *safeMap[K, V]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[K]V)
}

func (m *safeMap[K, V]) Keys() []K {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []K
	for k := range m.data {
		result = append(result, k)
	}
	return result
}

func (m *safeMap[K, V]) Vals() []V {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var vals []V
	for _, val := range m.data {
		vals = append(vals, val)
	}
	return vals
}
