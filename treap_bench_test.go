package treap_test

import (
	"testing"

	"github.com/lthibault/treap"
)

var (
	discard, discardRight *treap.Node
	discardNode           *treap.Node
)

func BenchmarkInsertSync(b *testing.B) {
	// To make this a fair benchmark, let's measure single-inserts to a non-empty,
	// balanced tree that is consistent across runs.
	root := newPrefilledTreap(handle, 1000)

	b.Run("NoPool", func(b *testing.B) {
		benchmarkInsertSync(b, handle, root)
	})

	b.Run("MemPool", func(b *testing.B) {
		var handle = treap.Handle{
			CompareWeights: treap.IntComparator,
			CompareKeys:    treap.IntComparator,
			NodeFactory:    treap.NewMemPool(),
		}

		b.Run("MemPool/Cold", func(b *testing.B) {
			benchmarkInsertSync(b, handle, root)
		})

		b.Run("MemPool/Warm", func(b *testing.B) {
			// warm up the mempool
			for i := 0; i < b.N; i++ {
				handle.NewNode().Free()
			}

			benchmarkInsertSync(b, handle, root)
		})
	})
}

func benchmarkInsertSync(b *testing.B, handle treap.Handle, root *treap.Node) {
	toInsert := mkTestCases(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for _, tc := range toInsert {
		discard, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
		discard.Free() // TODO:  consider testing without call to Free when benchmarking no-pool.
	}
}

func BenchmarkSplitSync(b *testing.B) {
	var root *treap.Node

	for _, tc := range mkTestCases(b.N) {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		discard, discardRight = handle.Split(root, i)
	}
}

func BenchmarkMergeSync(b *testing.B) {
	var root *treap.Node
	splits := make([]struct{ left, right *treap.Node }, b.N)

	for _, tc := range mkTestCases(b.N) {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	for i := 0; i < b.N; i++ {
		splits[i].left, splits[i].right = handle.Split(root, i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for _, s := range splits {
		discard = handle.Merge(s.left, s.right)
	}
}

func BenchmarkDeleteSync(b *testing.B) {
	root := newPrefilledTreap(handle, b.N)

	b.Run("NoPool", func(b *testing.B) {
		benchmarkDeleteSync(b, handle, root)
	})

	b.Run("MemPool", func(b *testing.B) {
		var handle = treap.Handle{
			CompareWeights: treap.IntComparator,
			CompareKeys:    treap.IntComparator,
			NodeFactory:    treap.NewMemPool(),
		}

		b.Run("Cold", func(b *testing.B) {
			benchmarkDeleteSync(b, handle, root)
		})

		b.Run("Warm", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				handle.NewNode().Free() // warm up the mempool
			}

			benchmarkDeleteSync(b, handle, root)
		})
	})
}

func benchmarkDeleteSync(b *testing.B, handle treap.Handle, root *treap.Node) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		discard = handle.Delete(root, i)
	}
}

func BenchmarkPopSync(b *testing.B) {
	root := newPrefilledTreap(handle, 1000)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, discard = handle.Pop(root)
		// don't call free; the popped node is still in use!
	}
}

func BenchmarkSetWeightSync(b *testing.B) {
	var root *treap.Node
	cs := mkTestCases(b.N)

	for _, tc := range cs {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i, tc := range cs {
		discard, _ = handle.SetWeight(root, tc.key, i)
	}
}

func BenchmarkIterSync(b *testing.B) {
	var root *treap.Node
	cs := mkTestCases(10)

	for _, tc := range cs {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for it := handle.Iter(root); it.Node != nil; it.Next() {
			discardNode = it.Node
		}

	}
}

func newPrefilledTreap(handle treap.Handle, n int) *treap.Node {
	var root *treap.Node
	cs := mkTestCases(n)

	for _, tc := range cs {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	return root
}
