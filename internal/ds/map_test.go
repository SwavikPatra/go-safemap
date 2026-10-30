package ds

import (
	"sort"
	"sync"
	"testing"
)

// Many goroutines write different keys, then we verify every write landed.
func TestSafeMap_ConcurrentPut(t *testing.T) {
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

// Delete returns true for an existing key and false for a missing one.
func TestSafeMap_Delete(t *testing.T) {
	m := NewSafeMap[string, int]()
	m.Put("a", 1)

	if !m.Delete("a") {
		t.Fatal("Delete(existing key) = false, want true")
	}
	if _, ok := m.Get("a"); ok {
		t.Fatal("key 'a' still present after Delete")
	}
	if m.Delete("a") {
		t.Fatal("Delete(already deleted key) = true, want false")
	}
	if m.Delete("never-existed") {
		t.Fatal("Delete(missing key) = true, want false")
	}
}

// Each key is deleted by exactly one goroutine, so every Delete must succeed
// and the map must end up empty.
func TestSafeMap_ConcurrentDelete(t *testing.T) {
	m := NewSafeMap[int, int]()
	const n = 500
	for i := 0; i < n; i++ {
		m.Put(i, i)
	}

	var wg sync.WaitGroup
	failures := make(chan int, n)
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

	for k := range failures {
		t.Errorf("Delete(%d) returned false, want true", k)
	}
	if got := len(m.Keys()); got != 0 {
		t.Fatalf("len(Keys()) = %d after deleting everything, want 0", got)
	}
}

// Many goroutines race to delete the SAME key: exactly one should win.
func TestSafeMap_ConcurrentDeleteSameKey(t *testing.T) {
	m := NewSafeMap[string, int]()
	m.Put("shared", 1)

	const goroutines = 100
	var wg sync.WaitGroup
	var mu sync.Mutex
	wins := 0

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if m.Delete("shared") {
				mu.Lock()
				wins++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if wins != 1 {
		t.Fatalf("Delete succeeded %d times, want exactly 1", wins)
	}
}

// Clear removes everything, and the map is still usable afterwards.
func TestSafeMap_Clear(t *testing.T) {
	m := NewSafeMap[int, string]()
	for i := 0; i < 10; i++ {
		m.Put(i, "v")
	}

	m.Clear()

	if got := len(m.Keys()); got != 0 {
		t.Fatalf("len(Keys()) = %d after Clear, want 0", got)
	}
	if _, ok := m.Get(3); ok {
		t.Fatal("key 3 still present after Clear")
	}

	// Must not panic (the internal map was re-created, not left nil).
	m.Put(1, "again")
	if v, ok := m.Get(1); !ok || v != "again" {
		t.Fatalf("Get(1) after Clear+Put = (%q, %v), want (\"again\", true)", v, ok)
	}
}

// Clear running while other goroutines write and read: mainly a race check.
func TestSafeMap_ConcurrentClear(t *testing.T) {
	m := NewSafeMap[int, int]()

	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(3)
		go func(i int) {
			defer wg.Done()
			m.Put(i, i)
		}(i)
		go func(i int) {
			defer wg.Done()
			m.Get(i)
		}(i)
		go func() {
			defer wg.Done()
			m.Clear()
		}()
	}
	wg.Wait()
	// No assertion on contents: the result depends on goroutine ordering.
	// The test passes if there is no race or panic under -race.
}

// Keys and Vals return everything (order is random, so sort before comparing).
func TestSafeMap_KeysAndVals(t *testing.T) {
	m := NewSafeMap[int, int]()
	for i := 0; i < 100; i++ {
		m.Put(i, i*10)
	}

	keys := m.Keys()
	vals := m.Vals()

	if len(keys) != 100 || len(vals) != 100 {
		t.Fatalf("len(keys)=%d len(vals)=%d, want 100 each", len(keys), len(vals))
	}

	sort.Ints(keys)
	sort.Ints(vals)
	for i := 0; i < 100; i++ {
		if keys[i] != i {
			t.Fatalf("keys[%d] = %d, want %d", i, keys[i], i)
		}
		if vals[i] != i*10 {
			t.Fatalf("vals[%d] = %d, want %d", i, vals[i], i*10)
		}
	}
}

// On an empty map, Keys and Vals return nothing (and don't panic).
func TestSafeMap_KeysAndValsEmpty(t *testing.T) {
	m := NewSafeMap[string, int]()

	if got := len(m.Keys()); got != 0 {
		t.Fatalf("len(Keys()) = %d on empty map, want 0", got)
	}
	if got := len(m.Vals()); got != 0 {
		t.Fatalf("len(Vals()) = %d on empty map, want 0", got)
	}
}

// The returned slices are snapshots: changing the map afterwards must not
// change a slice we already have.
func TestSafeMap_KeysAreSnapshot(t *testing.T) {
	m := NewSafeMap[int, int]()
	m.Put(1, 1)
	m.Put(2, 2)

	keys := m.Keys()
	m.Put(3, 3)
	m.Delete(1)

	if len(keys) != 2 {
		t.Fatalf("snapshot length changed to %d, want 2", len(keys))
	}
}

// Writers keep changing the map while readers take snapshots and iterate
// over them with no lock held. Mainly a race check.
func TestSafeMap_ConcurrentKeysVals(t *testing.T) {
	m := NewSafeMap[int, int]()
	for i := 0; i < 100; i++ {
		m.Put(i, i)
	}

	var wg sync.WaitGroup
	for i := 100; i < 300; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Put(i, i)
		}(i)
	}
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, k := range m.Keys() {
				_ = k
			}
			for _, v := range m.Vals() {
				_ = v
			}
		}()
	}
	wg.Wait()

	if got := len(m.Keys()); got != 300 {
		t.Fatalf("len(Keys()) = %d, want 300", got)
	}
}
