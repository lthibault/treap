package treap_test

import (
	"testing"

	"github.com/lthibault/treap"
)

var handle = treap.Comparator(func(a, b interface{}) int {
	ai := a.(int)
	bi := b.(int)

	switch {
	case ai < bi:
		return -1
	case ai > bi:
		return 1
	default:
		return 0
	}
})

func TestTreap(t *testing.T) {
	var root *treap.Node

	root = handle.Upsert(root, 0, 5, "a")
	root = handle.Upsert(root, 1, 7, "b")
	assertEq(t, "a", handle.Get(root, 5))
	assertEq(t, "b", handle.Get(root, 7))

	root = handle.Upsert(root, 2, 2, "c")
	assertEq(t, "c", handle.Get(root, 2))

	root = handle.Upsert(root, 2, 2, "d")
	assertEq(t, "d", handle.Get(root, 2))

	root = handle.Delete(root, 5)
	assertNil(t, handle.Get(root, 5))
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
