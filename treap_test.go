package treap_test

import (
	"math/rand" // don't seed; keep reproducible.
	"testing"

	"github.com/lthibault/treap"
	"github.com/stretchr/testify/assert"
)

type key int

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var handle = treap.Handle{
	CompareWeights: treap.IntComparator,
	CompareKeys: func(a, b interface{}) int {
		return treap.IntComparator(int(a.(key)), int(b.(key)))
	},
}

func TestTreap(t *testing.T) {
	var root *treap.Node

	var ok bool
	t.Run("Insert", func(t *testing.T) {
		t.Run("InsertBatch", func(t *testing.T) {
			/*
				Ensure insertion n+1 doesn't invalidate insertion n.
			*/

			root, ok = handle.Insert(root, key(7), "a", 1)
			assert.True(t, ok)
			root, ok = handle.Insert(root, key(2), "b", 11)
			assert.True(t, ok)

			res, ok := handle.Get(root, key(7))
			assert.Equal(t, "a", res)
			assert.True(t, ok)

			res, ok = handle.Get(root, key(2))
			assert.Equal(t, "b", res)
			assert.True(t, ok)
		})

		t.Run("InsertSingle", func(t *testing.T) {
			/*
				Ensure insertion n is immediately valid
			*/

			root, ok = handle.Insert(root, key(13), "c", -1)
			assert.True(t, ok)

			res, ok := handle.Get(root, key(13))
			assert.True(t, ok)
			assert.Equal(t, "c", res)
		})

		t.Run("Update", func(t *testing.T) {
			root, ok = handle.Upsert(root, key(13), "d", -1)
			assert.False(t, ok) // ensure it was created, not updated.

			res, ok := handle.Get(root, key(13))
			assert.True(t, ok)
			assert.Equal(t, "d", res)

		})

		t.Run("ValidateHeapOrder", func(t *testing.T) {
			if root.Weight != -1 {
				t.Error("min-heap ordering violated")
			}
		})
	})

	t.Run("Delete", func(t *testing.T) {
		t.Run("DeleteMissingValue", func(t *testing.T) {
			root = handle.Delete(root, key(5))
			_, ok = handle.Get(root, key(5))
			assert.False(t, ok)
		})

		t.Run("DeleteExistingValue", func(t *testing.T) {
			root = handle.Delete(root, key(13))
			_, ok = handle.Get(root, key(13))
			assert.False(t, ok)
		})

		t.Run("ValidateRemainingEntries", func(t *testing.T) {
			// Ensure old values are still present
			res, ok := handle.Get(root, key(7))
			assert.Equal(t, "a", res)
			assert.True(t, ok)

			res, ok = handle.Get(root, key(2))
			assert.Equal(t, "b", res)
			assert.True(t, ok)
		})

		t.Run("ValidateHeapOrder", func(t *testing.T) {
			if root.Weight != 1 {
				t.Error("min-heap ordering violated")
			}
		})
	})

	t.Run("InsertExistingFails", func(t *testing.T) {
		_, ok = handle.Get(root, key(2))
		assert.True(t, ok)

		// left branch
		_, ok = handle.Insert(root, key(2), "fail", 9001)
		assert.False(t, ok)

		// right branch
		root, _ = handle.Insert(root, key(9001), "d", 9001)
		new, ok := handle.Insert(root, key(9001), "fail", 0)
		assert.False(t, ok)

		if new != nil && new != root {
			t.Error("failed insert returned modified treap")
		}
	})
}

func TestPop(t *testing.T) {
	var root *treap.Node

	cs := mkTestCases(t, 1000)
	for _, tc := range cs {
		root, _ = handle.Insert(root, tc.key, tc.value, tc.weight)
	}

	for w := root.Weight; root != nil; _, root = handle.Pop(root) {
		assert.LessOrEqual(t, w, root.Weight,
			"heap property violated: %s < %s", root.Weight, w)
	}
}

func TestFuzz(t *testing.T) {
	/*
		For good measure, we perform a deterministic fuzz test.  We generate a large
		number of key-value pairs, insert them, and then perform a mix of updates and
		deletes, while ensuring the other entries are not invalidated by this process.
	*/

	const iter = 100
	var root *treap.Node

	cs := mkTestCases(t, iter)

	// Test insertions
	var ok bool
	var v interface{}
	for i, tc := range cs {
		if root, ok = handle.Insert(root, tc.key, tc.value, tc.weight); !ok {
			t.Error("insertion failed (key collision?)")
			t.FailNow()
		}

		v, ok = handle.Get(root, tc.key)
		assert.True(t, ok)
		assert.Equal(t, tc.value, v)

		testOthers(t, handle, root, cs[0:i])
	}

	// Test single deletions
	for i, tc := range cs {
		temp := handle.Delete(root, tc.key)

		v, ok = handle.Get(temp, tc.key)
		assert.False(t, ok)
		assert.Nil(t, v)

		testOthers(t, handle, temp, cs[0:i])
	}

	// Test cumulative deletions
	for i, tc := range cs {
		root = handle.Delete(root, tc.key)

		v, ok = handle.Get(root, tc.key)
		assert.False(t, ok)
		assert.Nil(t, v)

		if i < len(cs)-1 {
			testOthers(t, handle, root, cs[i+1:len(cs)-1])
		}
	}
}

func testOthers(t *testing.T, handle treap.Handle, root *treap.Node, cs []testCase) {
	for _, tc := range cs {
		v, ok := handle.Get(root, tc.key)
		assert.True(t, ok)
		assert.Equal(t, tc.value, v)
	}
}

type testCase struct {
	key    key
	value  string
	weight int
}

func mkTestCases(t *testing.T, n int) []testCase {
	cs := make([]testCase, n)
	for i := range cs {
		cs[i].key = key(i)
		cs[i].weight = i
		cs[i].value = randStr(5) // duplicates possible
	}

	// shuffle weights
	rand.Shuffle(n, func(i, j int) {
		cs[i].weight = j
		cs[j].weight = i
	})

	// shuffle keys
	rand.Shuffle(n, func(i, j int) {
		cs[i], cs[j] = cs[j], cs[i]
	})

	return cs
}

func randStr(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = rune(chars[rand.Intn(len(chars))])
	}
	return string(b)
}
