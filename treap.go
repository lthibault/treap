package treap

const minInt = -int((^uint(0))>>1) - 1

// Node is the recurisve datastructure that defines a persistent treap.
//
// The zero value is ready to use.
type Node struct {
	Weight      int
	Left, Right *Node
	Key, Value  interface{}
}

// Comparator establishes ordering between two elements.
// It returns -1 if a < b, 0 if a == b, and 1 if a > b.
type Comparator func(a, b interface{}) int

// Get an element by key.  Returns nil if the key is not in the treap.
// O(log n) if the tree is balanced (i.e. has uniformly distributed weights).
// Thread-safe.
func (f Comparator) Get(n *Node, key interface{}) interface{} {
	if n == nil {
		return nil
	}

	switch comp := f(key, n.Key); {
	case comp < 0:
		return f.Get(n.Left, key)
	case comp > 0:
		return f.Get(n.Right, key)
	default:
		return n.Value
	}
}

// Upsert updates an element, creating one if it is missing.
//
// O(log n) if the tree is balanced (see Get).
func (f Comparator) Upsert(n *Node, weight int, key, val interface{}) (res *Node) {
	if n == nil {
		return &Node{Weight: weight, Key: key, Value: val}
	}

	switch comp := f(key, n.Key); {
	case comp < 0:
		res = &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Value:  n.Value,
			Left:   f.Upsert(n.Left, weight, key, val),
			Right:  n.Right,
		}
	case comp > 0:
		res = &Node{
			Weight: n.Weight,
			Key:    n.Key,
			Value:  n.Value,
			Left:   n.Left,
			Right:  f.Upsert(n.Right, weight, key, val),
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

	if res.Left != nil && res.Left.Weight < res.Weight {
		return f.leftRotation(res)
	}

	if res.Right != nil && res.Right.Weight < res.Weight {
		return f.rightRotation(res)
	}

	return res
}

// Split a treap into its left and right branches at point `key`.  Key need not be
// present in the treap.  If it is, it WILL NOT be present in either of the resulting
// subtreaps.
//
// O(log n) if the treap is balanced (see Get).
func (f Comparator) Split(n *Node, key interface{}) (*Node, *Node) {
	ins := f.Upsert(n, minInt, key, nil) // minInt ensures `ins` is root node.
	return ins.Left, ins.Right
}

// Merge two treaps.  The root will be the root of the input treap with the lowest
// weight.
//
// O(log n) if the treap is balanced (see Get).
func (f Comparator) Merge(left, right *Node) *Node {
	switch {
	case left == nil:
		return right
	case right == nil:
		return left
	case left.Weight < right.Weight:
		return &Node{
			Weight: left.Weight,
			Key:    left.Key,
			Value:  left.Value,
			Left:   left.Left,
			Right:  f.Merge(left.Right, right),
		}
	default:
		return &Node{
			Weight: right.Weight,
			Key:    right.Key,
			Value:  right.Value,
			Left:   f.Merge(right.Left, left),
			Right:  right.Right,
		}
	}
}

// Delete a value.
//
// O(log n) if treap is balanced (see Get).
func (f Comparator) Delete(n *Node, key interface{}) *Node {
	return f.Merge(f.Split(n, key))
}

func (f Comparator) leftRotation(n *Node) *Node {
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

func (f Comparator) rightRotation(n *Node) *Node {
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
