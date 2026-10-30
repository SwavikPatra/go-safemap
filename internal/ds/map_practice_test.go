package ds

import (
	"sync"
	"testing"
)

func TestSafeMap_CuncurrentPut(t *testing.T) {
	m := NewSafeMap[int, int]()

	const n = 10
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
		if !ok || v != 2*i {
			t.Fatalf("key %d: got (%d, %v), want (%d, true)", i, v, ok, i*2)
		}
	}

}

func TestSafeMap_CuncurrentGetPut(t *testing.T) {
	m := NewSafeMap[string, int]()
	m.Put("counter", 0)
	var wg sync.WaitGroup
	const n = 100

	for i := 0; i < n; i++ {
		wg.Add(2)

		go func(i int) {
			defer wg.Done()
			m.Put("Counter", 1)
		}(i)

		go func() {
			defer wg.Done()
			m.Get("Counter")
		}()
	}
	wg.Wait()

	v, ok := m.Get("Counter")
	if !ok {
		t.Fatal("Key deleted")
	}
	if v < 0 || v > 100 {
		t.Fatal("Key currupted")
	}
}

func TestSafeMap_CuncurrentDelete(t *testing.T) {
	const n = 100
	var wg sync.WaitGroup
	failures := make(chan int, n)

	m := NewSafeMap[int, int]()

	for i := 0; i < n; i++ {
		m.Put(i, i)
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if !m.Delete(i) {
				failures <- i
			}
		}(i)
	}
	wg.Wait()
	close(failures)
	for i := range failures {
		t.Errorf("Delete %d returned false, we want true", i)
	}

	if got := len(m.Keys()); got != 0 {
		t.Fatalf("%d keys remained after deleting everything, want 0", got)
	}
}

func TestSafeMap_CuncurrentDeleteSameKey(t *testing.T) {
	m := NewSafeMap[string, int]()
	m.Put("Shared", 1)

	var wg sync.WaitGroup
	const n = 100
	count := 0
	var mu sync.Mutex
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if m.Delete("Shared") {
				mu.Lock()
				count++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if count != 1 {
		t.Fatalf("Count succeded %d times, want exactly 1", count)
	}
}
