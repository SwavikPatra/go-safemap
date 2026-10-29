package ds

import (
	"sync"
	"testing"
)

// Many goroutines write different keys, then we verify every write landed.
func TestSafeMap_ConcurrentPut(t *testing.T) {
	m := NewSafeMap[int, int]()
	const n = 1000

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Put(i, i*2)
		}(i)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		v, ok := m.Get(i)
		if !ok || v != i*2 {
			t.Fatalf("key %d: got (%d, %v), want (%d, true)", i, v, ok, i*2)
		}
	}
}

// Readers and writers hit the same key at the same time.
// This test mainly proves there is no data race or crash.
func TestSafeMap_ConcurrentReadWrite(t *testing.T) {
	m := NewSafeMap[string, int]()
	m.Put("counter", 0)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)

		go func(i int) { // writer
			defer wg.Done()
			m.Put("counter", i)
		}(i)

		go func() { // reader
			defer wg.Done()
			m.Get("counter")
		}()
	}
	wg.Wait()

	if _, ok := m.Get("counter"); !ok {
		t.Fatal("expected key 'counter' to exist")
	}
}
