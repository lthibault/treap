package treap

import (
	"sync"
)

// Node is the recurisve datastructure that defines a persistent treap.
//
// The zero value is ready to use.
type Node struct {
	Weight, Key, Value interface{}
	Left, Right        *Node
	free               func(*Node)
}

func (n *Node) Free() {
	if n != nil && n.free != nil {
		n.free(n)
	}
}

type NodeFactory interface {
	NewNode() *Node
}

// MemPool is a NodeFactory that uses sync.Pool to reuse nodes.
type MemPool sync.Pool

func NewMemPool() NodeFactory {
	p := new(MemPool)
	p.New = func() interface{} {
		return &Node{free: p.free}
	}

	return p
}

func (p *MemPool) NewNode() *Node {
	// TODO(performance):  this seems to be allocating a lot
	return (*sync.Pool)(p).Get().(*Node)
}

func (p *MemPool) free(n *Node) { (*sync.Pool)(p).Put(n) }
