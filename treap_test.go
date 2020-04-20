package treap_test

import (
	"math/rand" // don't seed; keep reproducible.
	"testing"

	"github.com/lthibault/treap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var handle = treap.Handle{
	CompareWeights: treap.IntComparator,
	CompareKeys:    treap.IntComparator,
}

func TestTreap(t *testing.T) {
	/*
		This test preforms a somewhat ecological test that combines operations.
		It is done in the spirit of an integration test.
	*/
	var root *treap.Node

	var ok bool
	t.Run("Insert", func(t *testing.T) {
		t.Run("InsertBatch", func(t *testing.T) {
			/*
				Ensure insertion n+1 doesn't invalidate insertion n.
			*/

			root, ok = handle.Insert(root, 7, "a", 1)
			assert.True(t, ok)
			root, ok = handle.Insert(root, 2, "b", 11)
			assert.True(t, ok)

			res, ok := handle.Get(root, 7)
			assert.Equal(t, "a", res)
			assert.True(t, ok)

			res, ok = handle.Get(root, 2)
			assert.Equal(t, "b", res)
			assert.True(t, ok)
		})

		t.Run("InsertSingle", func(t *testing.T) {
			/*
				Ensure insertion n is immediately valid
			*/

			root, ok = handle.Insert(root, 13, "c", -1)
			assert.True(t, ok)

			res, ok := handle.Get(root, 13)
			assert.True(t, ok)
			assert.Equal(t, "c", res)
		})

		t.Run("Update", func(t *testing.T) {
			root, ok = handle.Upsert(root, 13, "d", -1)
			assert.False(t, ok) // ensure it was created, not updated.

			res, ok := handle.Get(root, 13)
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
			root = handle.Delete(root, 5)
			_, ok = handle.Get(root, 5)
			assert.False(t, ok)
		})

		t.Run("DeleteExistingValue", func(t *testing.T) {
			root = handle.Delete(root, 13)
			_, ok = handle.Get(root, 13)
			assert.False(t, ok)
		})

		t.Run("ValidateRemainingEntries", func(t *testing.T) {
			// Ensure old values are still present
			res, ok := handle.Get(root, 7)
			assert.Equal(t, "a", res)
			assert.True(t, ok)

			res, ok = handle.Get(root, 2)
			assert.Equal(t, "b", res)
			assert.True(t, ok)
		})

		t.Run("ValidateHeapOrder", func(t *testing.T) {
			if root.Weight != 1 {
				t.Error("min-heap ordering violated")
			}
		})
	})
}

func TestInsert(t *testing.T) {
	var root *treap.Node
	cs := mkTestCases(100)

	t.Run("NewValue", func(t *testing.T) {
		for i, tc := range cs {
			new, ok := handle.Insert(root, tc.key, tc.value, tc.weight)
			assert.True(t, ok)
			assert.NotNil(t, new)

			t.Run("IsImmutable", func(t *testing.T) {
				require.NotEqual(t, root, new, "insertion %d", i)

				_, ok := handle.Get(root, tc.key)
				require.False(t, ok, "inserted value found in root (insertion %d)", i)
			})

			t.Run("IsRetrievable", func(t *testing.T) {
				v, ok := handle.Get(new, tc.key)
				assert.True(t, ok)
				assert.Equal(t, tc.value, v)
			})

			root = new
		}
	})

	t.Run("ExistingValue", func(t *testing.T) {
		t.Run("Fails", func(t *testing.T) {
			for i, tc := range cs {
				_, ok := handle.Insert(root, tc.key, tc.value, tc.weight)
				require.False(t, ok, "insertion %d overwrote value", i)
			}
		})

		t.Run("DoesNotChange", func(t *testing.T) {
			for i, tc := range cs {
				v, ok := handle.Get(root, tc.key)
				require.True(t, ok, "retrieval %d failed", i)
				require.Equal(t, tc.value, v, "value %d was modified", i)
			}
		})
	})
}

func TestSetWeight(t *testing.T) {
	var root *treap.Node
	cs := mkTestCases(100)

	var ok bool
	for _, tc := range cs {
		root, ok = handle.Insert(root, tc.key, tc.value, tc.weight)
		require.True(t, ok)
		require.NotNil(t, root)
	}

	t.Run("MissingKeyIsNop", func(t *testing.T) {
		_, ok = handle.SetWeight(root, 9999999999, nil)
		require.False(t, ok)
	})

	for i := 5; i > 0; i-- {
		new, ok := handle.SetWeight(root, cs[i].key, -i)
		assert.True(t, ok)
		assert.NotNil(t, new)
		require.NotEqual(t, root, new) // immutability
		root = new
	}

	var v interface{}
	for i := 5; i > 0; i-- {
		v, root = handle.Pop(root)
		assert.NotNil(t, root)
		assert.Equal(t, cs[i].value, v)
	}
}

func TestPop(t *testing.T) {
	var root *treap.Node

	v, tail := handle.Pop(nil)
	assert.Nil(t, v)
	assert.Nil(t, tail)

	cs := mkTestCases(1000)
	for _, tc := range cs {
		var ok bool
		root, ok = handle.Insert(root, tc.key, tc.value, tc.weight)
		require.True(t, ok)
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

	var root *treap.Node
	cs := mkTestCases(100)

	// Test insertions
	var ok bool
	var v interface{}
	for i, tc := range cs {
		root, ok = handle.Insert(root, tc.key, tc.value, tc.weight)
		require.True(t, ok)

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
	key    int
	value  string
	weight int
}

func mkTestCases(n int) []testCase {
	cs := make([]testCase, n)
	for i := range cs {
		cs[i].key = i
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
