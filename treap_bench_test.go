package treap_test

import (
	"testing"

	"github.com/lthibault/treap"
)

var root, discardLeft, discardRight *treap.Node

func BenchmarkInsertSync(b *testing.B) {
	root = nil
	cs := mkTestCases(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for _, tc := range cs {
		root, _ = handle.Upsert(root, tc.key, tc.value, tc.weight)
	}
}

func BenchmarkSplitSync(b *testing.B) {
	root = nil
	discardLeft = nil
	discardRight = nil

	cs := mkTestCases(b.N)

	for _, tc := range cs {
		root, _ = handle.Upsert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		discardLeft, discardRight = handle.Split(root, key(i))
	}
}

func BenchmarkMergeSync(b *testing.B) {
	root = nil
	discardLeft = nil
	discardRight = nil

	cs := mkTestCases(b.N)

	for _, tc := range cs {
		root, _ = handle.Upsert(root, tc.key, tc.value, tc.weight)
	}

	splits := make([]struct{ left, right *treap.Node }, b.N)
	for i := 0; i < b.N; i++ {
		splits[i].left, splits[i].right = handle.Split(root, key(i))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for _, s := range splits {
		root = handle.Merge(s.left, s.right)
	}
}

func BenchmarkDeleteSync(b *testing.B) {
	root = nil
	cs := mkTestCases(b.N)

	for _, tc := range cs {
		root, _ = handle.Upsert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		root = handle.Delete(root, key(i))
	}
}

func BenchmarkPopSync(b *testing.B) {
	cs := mkTestCases(b.N)

	for _, tc := range cs {
		root, _ = handle.Upsert(root, tc.key, tc.value, tc.weight)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, root = handle.Pop(root)
	}
}

func getRune(i int) rune {
	return rune(chars[i%(len(chars)-1)])
}
