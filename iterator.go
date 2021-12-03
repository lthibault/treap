package treap

import "sync"

// Iterator contains treap iteration state.  Its methods are NOT thread-safe, but
// multiple concurrent iterators are supported.
type Iterator struct {
	*Node
	stack *stack
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
