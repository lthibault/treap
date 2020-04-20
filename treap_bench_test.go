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
	var root *treap.Node
	cs := mkTestCases(b.N * 2)
	is := cs[b.N:]
	cs = cs[0:b.N]

	// To make this a fair benchmark, let's measure single-inserts to a non-empty,
	// balanced tree that is consistent across runs.
	for _, tc := range cs {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for _, tc := range is {
		discard, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
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
	var root *treap.Node

	for _, tc := range mkTestCases(b.N) {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		discard = handle.Delete(root, i)
	}
}

func BenchmarkPopSync(b *testing.B) {
	var root *treap.Node

	for _, tc := range mkTestCases(b.N) {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, root = handle.Pop(root)
	}

	discard = root
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
		for it := handle.Iter(root); it.Next(); {
			discardNode = it.Node
		}

	}
}
