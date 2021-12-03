package treap_test

import (
	"testing"

	"github.com/lthibault/treap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIter_Empty(t *testing.T) {
	t.Parallel()

	var empty *treap.Node
	assert.Nil(t, handle.Iter(empty).Node,
		"iterator for empty root should have nil node")

	var i int
	for it := handle.Iter(empty); it.Node != nil; it.Next() {
		i++
	}
	assert.Zero(t, i, "empty iterator shoudl not loop")
}

func TestIter_SingleEntry(t *testing.T) {
	t.Parallel()

	var root = &treap.Node{
		Key:    1,
		Value:  1,
		Weight: 1,
	}

	assert.NotNil(t, handle.Iter(root).Node,
		"iterator for non-empty root should have non-nil node")

	var i int
	for it := handle.Iter(root); it.Node != nil; it.Next() {
		i++
	}
	assert.Equal(t, 1, i, "iterator should loop one time")
}

func TestIter_MultiEntry(t *testing.T) {
	t.Parallel()

	var (
		root *treap.Node
		ok   bool
		tt   = []treap.Node{
			{
				Key:    0,
				Value:  0,
				Weight: 0,
			},
			{
				Key:    1,
				Value:  1,
				Weight: 1,
			},
			{
				Key:    2,
				Value:  2,
				Weight: 2,
			},
			{
				Key:    3,
				Value:  3,
				Weight: 2,
			},
			{
				Key:    4,
				Value:  4,
				Weight: 2,
			},
			{
				Key:    5,
				Value:  5,
				Weight: 1,
			},
			{
				Key:    6,
				Value:  6,
				Weight: 2,
			},
		}
	)

	for _, n := range tt {
		root, ok = handle.Insert(root, n.Key, n.Value, n.Weight)
		require.True(t, ok, "precondition failed: insert must succeed")
	}

	var ns []treap.Node
	for it := handle.Iter(root); it.Node != nil; it.Next() {
		t.Log(it.Node.Key)
		ns = append(ns, *it.Node)
	}
	assert.Len(t, ns, len(tt), "iterator should loop %d times", len(tt))

	for i, n := range tt {
		assert.Equal(t, n.Key, ns[i].Key, "iterator should traverse in key order")
	}
}
