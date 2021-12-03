package treap

import "sync"

// Iterator contains treap iteration state.  Its methods are NOT thread-safe, but
// multiple concurrent iterators are supported.
type Iterator struct {
	*Node
	stack *stack
}

// Iter walks the tree in key-order.
func (h Handle) Iter(n *Node) *Iterator {
	it := iterPool.Get().(*Iterator)
	it.stack = push(nil, n)
	it.Next()
	return it
}

// Next item.
func (it *Iterator) Next() {
	// are we resuming?
	if it.Node != nil {
		it.Node = it.Node.Right
	} else {
		it.Node, it.stack = pop(it.stack)
	}

	for {
		if it.Node == nil {
			it.Node, it.stack = pop(it.stack)
			break
		}

		it.stack = push(it.stack, it.Node)
		it.Node = it.Node.Left
	}

	if it.Node == nil {
		it.Finish()
	}
}

// Finish SHOULD be called if the caller has not exhausted the iterator.
// An iterator is exhausted when 'it.Node' is nil.
func (it *Iterator) Finish() {
	// return stack frames to the pool
	for it.stack != nil {
		_, it.stack = pop(it.stack)
	}

	iterPool.Put(it)
}

// Stack is a singly-linked list of nodes.
type stack struct {
	*Node
	next *stack
}

func push(s *stack, n *Node) *stack {
	if n == nil {
		return s
	}

	ss := stackPool.Get().(*stack)
	ss.Node = n
	ss.next = s
	return ss
}

func pop(s *stack) (n *Node, tail *stack) {
	if s != nil {
		defer stackPool.Put(s)
		n, tail = s.Node, s.next
	}

	return

}

var stackPool = sync.Pool{
	New: func() interface{} { return &stack{} },
}

var iterPool = sync.Pool{
	New: func() interface{} { return &Iterator{} },
}
