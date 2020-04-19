package treap_test

import (
	"testing"

	"github.com/lthibault/treap"
)

var handle = treap.Handle{
	CompareWeights: treap.IntComparator,
	CompareKeys:    treap.IntComparator,
}

func TestTreap(t *testing.T) {
	var root *treap.Node

	t.Run("Insert", func(t *testing.T) {
		root = handle.Upsert(root, 1, 7, "a")
		root = handle.Upsert(root, 11, 2, "b")
		assertEq(t, "a", handle.Get(root, 7))
		assertEq(t, "b", handle.Get(root, 2))

		root = handle.Upsert(root, -1, 13, "c")
		assertEq(t, "c", handle.Get(root, 13))
	})

	t.Run("Update", func(t *testing.T) {
		root = handle.Upsert(root, -1, 13, "d")
		assertEq(t, "d", handle.Get(root, 13))
	})

	t.Run("Delete", func(t *testing.T) {
		root = handle.Delete(root, 5)
		assertNil(t, handle.Get(root, 5))

		root = handle.Delete(root, 13)
		assertNil(t, handle.Get(root, 13))
	})
}

func assertEq(t *testing.T, expected string, actual interface{}) {
	if actual.(string) != expected {
		t.Errorf("expected %s, got %s", expected, actual)
		t.FailNow()
	}
}

func assertNil(t *testing.T, v interface{}) {
	if v != nil {
		t.Errorf("expected nil, got %v", v)
		t.FailNow()
	}
}
