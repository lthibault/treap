package treap

// Node is the recurisve datastructure that defines a persistent treap.
//
// The zero value is ready to use.
type Node struct {
	Weight, Key, Value interface{}
	Left, Right        *Node
}

// Handle performs purely functional transformations on a treap.
type Handle struct {
	CompareWeights, CompareKeys Comparator
}

// Get an element by key.  Returns nil if the key is not in the treap.
// O(log n) if the treap is balanced (i.e. has uniformly distributed weights).
func (h Handle) Get(n *Node, key interface{}) (v interface{}, found bool) {
	if n, found = h.GetNode(n, key); found {
		v = n.Value
	}
	return
}

// GetNode returns the subtree whose root has the specified key.  This is equivalent to
// Get, but returns a full node.
func (h Handle) GetNode(n *Node, key interface{}) (*Node, bool) {
	if n == nil {
		return nil, false
	}

	switch comp := h.CompareKeys(key, n.Key); {
	case comp < 0:
		return h.GetNode(n.Left, key)
	case comp > 0:
		return h.GetNode(n.Right, key)
	default:
		return n, true
	}
}

// Insert an element into the treap, returning false if the element is already present.
//
// O(log n) if the treap is balanced (see Get).
func (h Handle) Insert(n *Node, key, val, weight interface{}) (new *Node, ok bool) {
	return h.upsert(n, key, val, weight, true, false, nil)
}

// SetWeight adjusts the weight of the specified item.  It is a nop if the key is not in
// the treap, in which case the returned bool is `false`.
//
// O(log n) if the treap is balanced (see Get).
func (h Handle) SetWeight(n *Node, key, weight interface{}) (new *Node, ok bool) {
	new, _ = h.upsert(n, key, nil, weight, false, true, nil)
	ok = new != nil
	return
}

// Upsert updates an element, creating one if it is missing.
//
// O(log n) if the treap is balanced (see Get).
func (h Handle) Upsert(n *Node, key, val, weight interface{}) (new *Node, created bool) {
	return h.upsert(n, key, val, weight, true, true, nil)
}

// UpsertIf the comparison function returns true.  This is functionally equivalent to a
// Get followed by an Upsert, but faster.
//
// Return values and time complexity are identical to Upsert.
func (h Handle) UpsertIf(n *Node, key, val, weight interface{}, comp func(*Node) bool) (*Node, bool) {
	return h.upsert(n, key, val, weight, true, true, comp)
}

func (h Handle) upsert(n *Node, k, v, w interface{}, create, update bool, fn func(*Node) bool) (res *Node, created bool) {
	if n == nil {
		if create {
			created = true
			res = &Node{Weight: w, Key: k, Value: v}
		}

		return
	}

	switch h.CompareKeys(k, n.Key) {
	case -1:
		// use res as temp variable to avoid extra allocation
		if res, created = h.upsert(n.Left, k, v, w, create, update, fn); res == nil {
			return
		}

		res = &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Value:  n.Value,
			Left:   res,
			Right:  n.Right,
		}
	case 1:
		// use res as temp variable to avoid extra allocation
		if res, created = h.upsert(n.Right, k, v, w, create, update, fn); res == nil {
			return
		}

		res = &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Value:  n.Value,
			Left:   n.Left,
			Right:  res,
		}
	default:
		if !update { // insert only (no upsert)
			return
		}

		if fn != nil && !fn(n) { // InsertIf decided to ignore
			res = n
			return
		}

		res = &Node{
			Weight: w,
			Key:    n.Key,
			Value:  n.Value,
			Left:   n.Left,
			Right:  n.Right,
		}

		if create { // not SetWeight
			res.Value = v // upsert; set new value.
		}
	}

	if res.Left != nil && h.CompareWeights(res.Left.Weight, res.Weight) < 0 {
		res = h.leftRotation(res)
	} else if res.Right != nil && h.CompareWeights(res.Right.Weight, res.Weight) < 0 {
		res = h.rightRotation(res)
	}

	return
}

// Split a treap into its left and right branches at point `key`.  Key need not be
// present in the treap.  If it is, it WILL NOT be present in either of the resulting
// subtreaps.
//
// O(log n) if the treap is balanced (see Get).
func (h Handle) Split(n *Node, key interface{}) (*Node, *Node) {
	ins, _ := h.Upsert(n, key, nil, nil)
	return ins.Left, ins.Right
}

// Merge two treaps.  The root will be the root of the input treap with the lowest
// weight.
//
// O(log n) if the treap is balanced (see Get).
func (h Handle) Merge(left, right *Node) *Node {
	switch {
	case left == nil:
		return right
	case right == nil:
		return left
	case h.CompareWeights(left.Weight, right.Weight) < 0:
		return &Node{
			Weight: left.Weight,
			Key:    left.Key,
			Value:  left.Value,
			Left:   left.Left,
			Right:  h.Merge(left.Right, right),
		}
	default:
		return &Node{
			Weight: right.Weight,
			Key:    right.Key,
			Value:  right.Value,
			Left:   h.Merge(left, right.Left),
			Right:  right.Right,
		}
	}
}

// Delete a value.
//
// O(log n) if treap is balanced (see Get).
func (h Handle) Delete(n *Node, key interface{}) *Node {
	return h.Merge(h.Split(n, key))
}

// Pop the next value off the heap.  By default, this is the item with the lowest
// weight.
//
// Pop is equivalent to calling Delete on a root node's key, but avoids an O(n) insert
// operation.
//
// O(log n)
func (h Handle) Pop(n *Node) (interface{}, *Node) {
	if n == nil {
		return nil, nil
	}

	return n.Value, h.Merge(n.Left, n.Right)
}

// Iter walks the tree in key-order.
func (h Handle) Iter(n *Node) Iterator {
	return Iterator{
		n:     n,
		stack: make([]*Node, 0, 16),
	}
}

func (h Handle) leftRotation(n *Node) *Node {
	return &Node{
		Weight: n.Left.Weight,
		Key:    n.Left.Key,
		Value:  n.Left.Value,
		Left:   n.Left.Left,
		Right: &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Value:  n.Value,
			Left:   n.Left.Right,
			Right:  n.Right,
		},
	}
}

func (h Handle) rightRotation(n *Node) *Node {
	return &Node{
		Weight: n.Right.Weight,
		Key:    n.Right.Key,
		Value:  n.Right.Value,
		Left: &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Value:  n.Value,
			Left:   n.Left,
			Right:  n.Right.Left,
		},
		Right: n.Right.Right,
	}
}

// Iterator contains treap iteration state.  Its methods are NOT thread-safe, but
// multiple concurrent iterators are supported.
type Iterator struct {
	*Node
	n     *Node
	stack []*Node
}

// Next item.
func (it *Iterator) Next() (more bool) {
	for more = it.More(); more; {
		if it.n != nil {
			it.stack = append(it.stack, it.n)
			it.n = it.n.Left
			continue
		}

		it.n, it.stack = pop(it.stack)
		it.Node = it.n
		it.n = it.n.Right

		break
	}

	return
}

// More returns true if the iterator has not traversed the entire treap.
func (it *Iterator) More() bool {
	return len(it.stack) > 0 || it.n != nil
}

func pop(stack []*Node) (*Node, []*Node) {
	return stack[len(stack)-1], stack[:len(stack)-1]
}
