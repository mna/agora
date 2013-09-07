package runtime

import (
	"testing"
)

func TestAddMixed(t *testing.T) {
	cases := []struct {
		x   Val
		y   Val
		exp Val
		p   bool
	}{
		// Number leads
		{x: Number(5), y: Number(3.24), exp: Number(8.24)},
		{x: Number(5), y: String("2"), exp: Number(7)},
		{x: Number(5), y: String("2.2"), exp: Number(7.2)},
		{x: Number(5), y: String("whatever"), exp: Nil, p: true},
		{x: Number(5), y: Bool(true), exp: Number(6)},
		{x: Number(5), y: Nil, exp: Nil, p: true},
		{x: Number(2.2), y: Number(3), exp: Number(5.2)},
		{x: Number(2.2), y: String("3"), exp: Number(5.2)},
		{x: Number(2.2), y: String("3.4"), exp: Number(5.6)},
		{x: Number(2.2), y: String("test"), exp: Nil, p: true},
		{x: Number(2.2), y: Bool(true), exp: Number(3.2)},
		{x: Number(2.2), y: Nil, exp: Nil, p: true},
		// String leads
		{x: String("some"), y: Number(45), exp: String("some45")},
		{x: String("some"), y: Number(-5), exp: String("some-5")},
		{x: String("some"), y: Number(2.23), exp: String("some2.23")},
		{x: String("some"), y: Number(2.00023), exp: String("some2.00023")},
		{x: String("some"), y: Number(-0.04), exp: String("some-0.04")},
		{x: String("some"), y: Bool(true), exp: String("sometrue")},
		{x: String("some"), y: Bool(false), exp: String("somefalse")},
		{x: String("some"), y: Nil, exp: String("somenil")},
		// Bool and Nil panic for Add
	}

	assert := func(x, y Val) {
		if err := recover(); err == nil {
			t.Errorf("%v and %v : expected an error, got none", x, y)
		}
	}
	for _, c := range cases {
		func() {
			if c.p {
				defer assert(c.x, c.y)
			}
			res := c.x.Add(c.y)
			if res != c.exp {
				t.Errorf("%v + %v : expected %v, got %v", c.x, c.y, c.exp, res)
			}
		}()
	}
}
