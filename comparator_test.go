package treap_test

import (
	"testing"
	"time"

	"github.com/lthibault/treap"
	"github.com/stretchr/testify/assert"
)

func TestMaxTreap(t *testing.T) {
	comp := treap.MaxTreap(treap.IntComparator)

	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil > 1",
		test: []interface{}{nil, 1, 1},
	}, {
		desc: "1 < nil",
		test: []interface{}{1, nil, -1},
	}, {
		desc: "1 > 2",
		test: []interface{}{1, 2, 1},
	}, {
		desc: "2 < 1",
		test: []interface{}{2, 1, -1},
	}, {
		desc: "0 == 0",
		test: []interface{}{0, 0, 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], comp(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestIntComparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, 1, -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{1, nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{1, 1, 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{1, 2, -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{2, 1, 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{0, 0, 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.IntComparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestInt8Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, int8(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{int8(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{int8(1), int8(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{int8(1), int8(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{int8(2), int8(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{int8(0), int8(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.Int8Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestInt16Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, int16(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{int16(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{int16(1), int16(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{int16(1), int16(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{int16(2), int16(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{int16(0), int16(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.Int16Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestInt32Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, int32(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{int32(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{int32(1), int32(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{int32(1), int32(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{int32(2), int32(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{int32(0), int32(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.Int32Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestInt64Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, int64(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{int64(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{int64(1), int64(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{int64(1), int64(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{int64(2), int64(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{int64(0), int64(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.Int64Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestStringComparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: `nil < ""`,
		test: []interface{}{nil, "", -1},
	}, {
		desc: `"" > nil`,
		test: []interface{}{"", nil, 1},
	}, {
		desc: `"alpha" == "alpha"`,
		test: []interface{}{"alpha", "alpha", 0},
	}, {
		desc: `"" < "bravo"`,
		test: []interface{}{"", "bravo", -1},
	}, {
		desc: `"alpha" > ""`,
		test: []interface{}{"alpha", "", 1},
	}, {
		desc: `"" == ""`,
		test: []interface{}{"", "", 0},
	}, {
		desc: `"alpha" < "bravo"`,
		test: []interface{}{"alpha", "bravo", -1},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.StringComparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestByteComparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: `nil < empty`,
		test: []interface{}{nil, []byte{}, -1},
	}, {
		desc: `empty > nil`,
		test: []interface{}{[]byte{}, nil, 1},
	}, {
		desc: `"alpha" == "alpha"`,
		test: []interface{}{[]byte("alpha"), []byte("alpha"), 0},
	}, {
		desc: `empty < "bravo"`,
		test: []interface{}{[]byte{}, []byte("bravo"), -1},
	}, {
		desc: `"alpha" > empty`,
		test: []interface{}{[]byte("alpha"), []byte{}, 1},
	}, {
		desc: `empty == empty`,
		test: []interface{}{[]byte{}, []byte{}, 0},
	}, {
		desc: `"alpha" < "bravo"`,
		test: []interface{}{[]byte("alpha"), []byte("bravo"), -1},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.BytesComparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestTimeComparator(t *testing.T) {

	t0 := time.Now()

	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < long ago",
		test: []interface{}{nil, t0.Add(time.Hour * -999999), -1},
	}, {
		desc: "long ago > nil",
		test: []interface{}{t0.Add(time.Hour * -999999), nil, 1},
	}, {
		desc: "t0 == t0",
		test: []interface{}{t0, t0, 0},
	}, {
		desc: "two weeks from now > t0",
		test: []interface{}{t0.Add(24 * 7 * 2 * time.Hour), t0, 1},
	}, {
		desc: "t0 < two weeks from now",
		test: []interface{}{t0, t0.Add(24 * 7 * 2 * time.Hour), -1},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.TimeComparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestUIntComparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, uint(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{uint(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{uint(1), uint(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{uint(1), uint(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{uint(2), uint(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{uint(0), uint(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.UIntComparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestUInt8Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, uint8(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{uint8(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{uint8(1), uint8(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{uint8(1), uint8(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{uint8(2), uint8(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{uint8(0), uint8(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.UInt8Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestUInt16Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, uint16(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{uint16(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{uint16(1), uint16(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{uint16(1), uint16(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{uint16(2), uint16(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{uint16(0), uint16(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.UInt16Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestUInt32Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, uint32(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{uint32(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{uint32(1), uint32(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{uint32(1), uint32(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{uint32(2), uint32(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{uint32(0), uint32(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.UInt32Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestUInt64Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, uint64(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{uint64(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{uint64(1), uint64(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{uint64(1), uint64(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{uint64(2), uint64(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{uint64(0), uint64(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.UInt64Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestFloat32Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, float32(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{float32(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{float32(1), float32(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{float32(1), float32(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{float32(2), float32(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{float32(0), float32(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.Float32Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}

func TestFloat64Comparator(t *testing.T) {
	for _, tc := range []struct {
		desc string
		test []interface{}
	}{{
		desc: "nil < 1",
		test: []interface{}{nil, float64(1), -1},
	}, {
		desc: "1 > nil",
		test: []interface{}{float64(1), nil, 1},
	}, {
		desc: "1 == 1",
		test: []interface{}{float64(1), float64(1), 0},
	}, {
		desc: "1 < 2",
		test: []interface{}{float64(1), float64(2), -1},
	}, {
		desc: "2 > 1",
		test: []interface{}{float64(2), float64(1), 1},
	}, {
		desc: "0 == 0",
		test: []interface{}{float64(0), float64(0), 0},
	}} {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.test[2], treap.Float64Comparator(tc.test[0], tc.test[1]),
				"constraint %s violated", tc.desc)
		})
	}
}
