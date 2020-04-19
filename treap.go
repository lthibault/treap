package treap

// Handle performs purely functional transformations on a treap.
type Handle struct {
	CompareWeights, CompareKeys Comparator
}

// Node is the recurisve datastructure that defines a persistent treap.
//
// The zero value is ready to use.
type Node struct {
	Weight, Key, Value interface{}
	Left, Right        *Node
}

// Comparator establishes ordering between two elements.
// It returns -1 if a < b, 0 if a == b, and 1 if a > b.
// Nil values are treated as -Inf.
type Comparator func(a, b interface{}) int

// Get an element by key.  Returns nil if the key is not in the treap.
// O(log n) if the tree is balanced (i.e. has uniformly distributed weights).
// Thread-safe.
func (h Handle) Get(n *Node, key interface{}) interface{} {
	if n == nil {
		return nil
	}

	switch comp := h.CompareKeys(key, n.Key); {
	case comp < 0:
		return h.Get(n.Left, key)
	case comp > 0:
		return h.Get(n.Right, key)
	default:
		return n.Value
	}
}

// Upsert updates an element, creating one if it is missing.
//
// O(log n) if the tree is balanced (see Get).
func (h Handle) Upsert(n *Node, weight, key, val interface{}) (res *Node) {
	if n == nil {
		return &Node{Weight: weight, Key: key, Value: val}
	}

	switch comp := h.CompareKeys(key, n.Key); {
	case comp < 0:
		res = &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Value:  n.Value,
			Left:   h.Upsert(n.Left, weight, key, val),
			Right:  n.Right,
		}
	case comp > 0:
		res = &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Value:  n.Value,
			Left:   n.Left,
			Right:  h.Upsert(n.Right, weight, key, val),
		}
	default:
		res = &Node{
			Weight: weight,
			Key:    n.Key,
			Value:  val,
			Left:   n.Left,
			Right:  n.Right,
		}
	}

	if res.Left != nil && h.CompareWeights(res.Left.Weight, res.Weight) < 0 {
		return h.leftRotation(res)
	}

	if res.Right != nil && h.CompareWeights(res.Right.Weight, res.Weight) < 0 {
		return h.rightRotation(res)
	}

	return res
}

// Split a treap into its left and right branches at point `key`.  Key need not be
// present in the treap.  If it is, it WILL NOT be present in either of the resulting
// subtreaps.
//
// O(log n) if the treap is balanced (see Get).
func (h Handle) Split(n *Node, key interface{}) (*Node, *Node) {
	ins := h.Upsert(n, nil, key, nil)
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
			Left:   h.Merge(right.Left, left),
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
