package treap

import (
	"time"
	"unsafe"
)

// Comparator establishes ordering between two elements.
// It returns -1 if a < b, 0 if a == b, and 1 if a > b.
// Nil values are treated as -Inf.
type Comparator func(a, b interface{}) int

// MaxTreap wraps a comparator, resulting in a treap with max-heap ordering.
func MaxTreap(f Comparator) Comparator {
	return func(a, b interface{}) int {
		return -f(a, b)
	}
}

// IntComparator compares integers.  Nil values are considered infinite.
func IntComparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(int)
	bAsserted := b.(int)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// StringComparator provides a fast comparison on strings.
// Nil values are treated as infinite.
func StringComparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	s1 := a.(string)
	s2 := b.(string)
	min := len(s2)
	if len(s1) < len(s2) {
		min = len(s1)
	}
	diff := 0
	for i := 0; i < min && diff == 0; i++ {
		diff = int(s1[i]) - int(s2[i])
	}
	if diff == 0 {
		diff = len(s1) - len(s2)
	}
	if diff < 0 {
		return -1
	}
	if diff > 0 {
		return 1
	}
	return 0
}

// Int8Comparator provides a basic comparison on int8.
// Nil values are treated as infinite.
func Int8Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(int8)
	bAsserted := b.(int8)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// Int16Comparator provides a basic comparison on int16.
// Nil values are treated as infinite.
func Int16Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(int16)
	bAsserted := b.(int16)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// Int32Comparator provides a basic comparison on int32.
// Nil values are treated as infinite.
func Int32Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(int32)
	bAsserted := b.(int32)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// Int64Comparator provides a basic comparison on int64.
// Nil values are treated as infinite.
func Int64Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(int64)
	bAsserted := b.(int64)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// UIntComparator provides a basic comparison on uint.
// Nil values are treated as infinite.
func UIntComparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(uint)
	bAsserted := b.(uint)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// UInt8Comparator provides a basic comparison on uint8.
// Nil values are treated as infinite.
func UInt8Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(uint8)
	bAsserted := b.(uint8)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// UInt16Comparator provides a basic comparison on uint16.
// Nil values are treated as infinite.
func UInt16Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(uint16)
	bAsserted := b.(uint16)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// UInt32Comparator provides a basic comparison on uint32.
// Nil values are treated as infinite.
func UInt32Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(uint32)
	bAsserted := b.(uint32)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// UInt64Comparator provides a basic comparison on uint64.
// Nil values are treated as infinite.
func UInt64Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(uint64)
	bAsserted := b.(uint64)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// Float32Comparator provides a basic comparison on float32.
// Nil values are treated as infinite.
func Float32Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(float32)
	bAsserted := b.(float32)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// Float64Comparator provides a basic comparison on float64.
// Nil values are treated as infinite.
func Float64Comparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(float64)
	bAsserted := b.(float64)
	switch {
	case aAsserted < bAsserted:
		return -1
	case aAsserted > bAsserted:
		return 1
	default:
		return 0
	}
}

// BytesComparator provides a basic comparison on []byte.
// Nil values are treated as infinite.
func BytesComparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.([]byte)
	bAsserted := b.([]byte)

	return StringComparator(
		*(*string)(unsafe.Pointer(&aAsserted)),
		*(*string)(unsafe.Pointer(&bAsserted)),
	)
}

// TimeComparator provides a basic comparison on time.Time.
// Nil values are treated as infinite.
func TimeComparator(a, b interface{}) int {
	switch {
	case a == nil:
		return -1 // N.B.:  treap is a min-heap by default
	case b == nil:
		return 1
	}

	aAsserted := a.(time.Time)
	bAsserted := b.(time.Time)

	switch {
	case aAsserted.After(bAsserted):
		return 1
	case aAsserted.Before(bAsserted):
		return -1
	default:
		return 0
	}
}
