package main

import (
	"fmt"

	"github.com/lthibault/treap"
)

// Treap operations are performed by a lightweight handle.  Usually, you'll create a
// single global handle and share it between goroutines.  Handle's methods are thread-
// safe.
//
// A handle is defined by it's comparison functions (type `treap.Comparator`).
var handle = treap.Handle{
	// CompareKeys is used to store and receive mapped entries.  The comparator must be
	// compatible with the Go type used for keys.  In this example, we'll use strings as
	// keys.
	CompareKeys: treap.StringComparator,

	// CompareWeights is used to maintain priority-ordering of mapped entries, providing
	// us with fast `Pop`, `Insert` and `SetWeight` operations.  You'll usually want
	// to use a `treap.IntComparator` for weights, but you can use any comparison
	// function you require.  Try it with `treap.TimeComparator`!
	//
	// Note that treaps are min-heaps by default, so `Pop` will always return the item
	// with the _smallest_ weight.  You can easily switch to a max-heap by using
	// `treap.MaxTreap`, if required.
	CompareWeights: treap.IntComparator,
}

func main() {
	// We define an empty root node.  Don't worry -- there's no initialization required!
	var root *treap.Node

	// We're going to insert each of these boxers into the treap, and observe how the
	// treap treap provides us with a combination of map and heap semantics.
	for _, boxer := range []struct {
		FirstName, LastName string
		Weight              int
	}{{
		FirstName: "Cassius",
		LastName:  "Clay",
		Weight:    210,
	}, {
		FirstName: "Joe",
		LastName:  "Frazier",
		Weight:    215,
	}, {
		FirstName: "Marcel",
		LastName:  "Cerdan",
		Weight:    154,
	}, {
		FirstName: "Jake",
		LastName:  "LaMotta",
		Weight:    160,
	}} {
		// Again, the treap is a purely-functional, persistent data structure.  `Insert`
		// returns a _new_ heap, which replaces `root` on each iteration.
		//
		// When used in conjunction with `atomic.CompareAndSwapPointer`, it is possible
		// to read from a treap without ever blocking -- even in the presence of
		// concurrent writers!
		root, _ = handle.Insert(root, boxer.FirstName, boxer.LastName, boxer.Weight)
	}

	// Now that we've populated the treap, we can query it like an ordinary map.
	lastn, _ := handle.Get(root, "Cassius")
	fmt.Printf("Cassius => %s\n", lastn) // prints:  "Cassius => Clay"

	// Treaps also behave like binary heaps.  Let's start by peeking at the first value
	// in the resulting priority queue.  Remember:  this is a min-heap by default.
	fmt.Printf("Head node is:\t\t%s %s, %d lbs\n", root.Key, root.Value, root.Weight)

	// Woah, that was easy!  Now let's Pop that first value off of the heap.
	// Remember:  this is an immutable data-structure, so `Pop` doesn't actually mutate
	// any state!
	lastn, _ = handle.Pop(root)
	fmt.Printf("Popped head node:\tMarcel %s\n", lastn) // prints:  "Marcel Cerdan"

	// Jake LaMotta moved up to the heavyweight class late in his career.  Let's made an
	// adjustment to his weight.
	root, _ = handle.SetWeight(root, "Jake", 205)

	// Let's list our boxers in ascending order of weight.  You may have noticed
	// there's no `PopNode` method on `treap.Handler`.  This is not a mistake!  A `Pop`
	// is just a merge on the root node's subtrees.  Check it out:
	fmt.Println("\n[ heap traversal... ]")
	for n := root; n != nil; {
		fmt.Printf("%s %s: %d\n", n.Key, n.Value, n.Weight)
		n = handle.Merge(n.Left, n.Right)
	}

	// Lastly, we can iterate through the treap in key-order (smallest to largest).
	// To do this, we use an iterator.  Contrary to treaps, iterators are stateful and
	// mutable!  As such, they are NOT thread-safe.  However, multiple concurrent
	// iterators can traverse the same treap safely.
	var i int
	fmt.Println("\n[ binary search-tree traversal (notice keys are sorted alphabetically)... ]")
	for iterator := handle.Iter(root); iterator.Next(); i++ {
		fmt.Printf("[%d] %s %s: %d\n", i, iterator.Key, iterator.Value, iterator.Weight)
	}
}
