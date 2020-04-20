package treap_test

import (
	"math/rand"
	"testing"

	"github.com/lthibault/treap"
)

var root, discardLeft, discardRight *treap.Node

func BenchmarkInsertSync(b *testing.B) {
	ws := make([]int, b.N)
	vals := make([]rune, b.N)
	for i := range ws {
		ws[i] = i
		vals[i] = getRune(i)
	}

	rand.Shuffle(b.N, func(i, j int) {
		ws[i], ws[j] = ws[j], ws[i]
	})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		root, _ = handle.Upsert(root, key(i), vals[i], ws[i])
	}
}

func BenchmarkSplitSync(b *testing.B) {
	ws := make([]int, b.N)
	vals := make([]rune, b.N)
	for i := range ws {
		ws[i] = i
		vals[i] = getRune(i)
	}

	rand.Shuffle(b.N, func(i, j int) {
		ws[i], ws[j] = ws[j], ws[i]
	})

	for i := 0; i < b.N; i++ {
		root, _ = handle.Upsert(root, key(i), vals[i], ws[i])
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		root = handle.Delete(root, key(i))
	}
}

func BenchmarkDeleteSync(b *testing.B) {
	ws := make([]int, b.N)
	vals := make([]rune, b.N)
	for i := range ws {
		ws[i] = i
		vals[i] = getRune(i)
	}

	rand.Shuffle(b.N, func(i, j int) {
		ws[i], ws[j] = ws[j], ws[i]
	})

	for i := 0; i < b.N; i++ {
		root, _ = handle.Upsert(root, key(i), vals[i], ws[i])
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		discardLeft, discardRight = handle.Split(root, key(i))
	}
}

func BenchmarkPopSync(b *testing.B) {
	ws := make([]int, b.N)
	vals := make([]rune, b.N)
	for i := range ws {
		ws[i] = i
		vals[i] = getRune(i)
	}

	rand.Shuffle(b.N, func(i, j int) {
		ws[i], ws[j] = ws[j], ws[i]
	})

	for i := 0; i < b.N; i++ {
		root, _ = handle.Upsert(root, key(i), vals[i], ws[i])
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
