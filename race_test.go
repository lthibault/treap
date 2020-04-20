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
	var root = unsafe.Pointer(&treap.Node{
		Weight: 0,
		Key:    key(0),
		Value:  "a",
	})

	var wg sync.WaitGroup
	wg.Add(len(chars))

	ch := make(chan struct{})

	for i := 0; i < len(chars); i++ {
		go func(k int, val rune) {
			defer wg.Done()

			<-ch // try to get as many read/writes happening at the same time

			for i := 0; i < 1000; i++ {
				switch {
				case i&k == 0:
					// Write
					for {
						old := (*treap.Node)(atomic.LoadPointer(&root))
						if new, _ := handle.Upsert(old, key(k), val, k); atomic.CompareAndSwapPointer(
							&root,
							unsafe.Pointer(old),
							unsafe.Pointer(new),
						) {
							break
						}
					}
				case i&k == k-1:
					// Delete
					for {
						old := (*treap.Node)(atomic.LoadPointer(&root))
						if new := handle.Delete(old, key(k)); atomic.CompareAndSwapPointer(
							&root,
							unsafe.Pointer(old),
							unsafe.Pointer(new),
						) {
							break
						}
					}
				default:
					// Read
					v, ok := handle.Get((*treap.Node)(atomic.LoadPointer(&root)), key(k))
					if ok && v.(rune) != val {
						t.Error("violation")
					}
				}
			}
		}(i, getRune(i))
	}

	close(ch)
	wg.Wait()
}
