// +build race

package treap_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/lthibault/treap"
)

// TestRace ensures there are no data races.  Only run when the -race flags is passed to
// `go test`.
func TestRace(t *testing.T) {
	const iter = 100

	var root = unsafe.Pointer(handle.Upsert(nil, 0, 0, "a"))

	var wg sync.WaitGroup
	wg.Add(iter)

	ch := make(chan struct{})

	for i := 0; i < iter; i++ {
		go func(key int, val rune) {
			defer wg.Done()

			<-ch // try to get as many read/writes happening at the same time

			for i := 0; i < 1000; i++ {
				switch {
				case i&key == 0:
					// Write
					for {
						old := (*treap.Node)(atomic.LoadPointer(&root))
						if new := handle.Upsert(old, key, key, val); atomic.CompareAndSwapPointer(
							&root,
							unsafe.Pointer(old),
							unsafe.Pointer(new),
						) {
							break
						}
					}
				case i&key == key-1:
					// Delete
					for {
						old := (*treap.Node)(atomic.LoadPointer(&root))
						if new := handle.Delete(old, key); atomic.CompareAndSwapPointer(
							&root,
							unsafe.Pointer(old),
							unsafe.Pointer(new),
						) {
							break
						}
					}
				default:
					// Read
					v := handle.Get((*treap.Node)(atomic.LoadPointer(&root)), key)
					if v != nil && v.(rune) != val {
						t.Error("violation")
					}
				}
			}
		}(i, getRune(i))
	}

	close(ch)
	wg.Wait()
}
