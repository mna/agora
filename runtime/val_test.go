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
		// Int leads
		{x: Int(5), y: Float(3.24), exp: Int(8)},
		{x: Int(5), y: String("2"), exp: Int(7)},
		{x: Int(5), y: String("2.2"), exp: Nil, p: true}, // TODO : is this the correct behaviour?
		{x: Int(5), y: String("whatever"), exp: Nil, p: true},
		{x: Int(5), y: Bool(true), exp: Int(6)},
		{x: Int(5), y: Nil, exp: Nil, p: true},
		// Float leads
		{x: Float(2.2), y: Int(3), exp: Float(5.2)},
		{x: Float(2.2), y: String("3"), exp: Float(5.2)},
		{x: Float(2.2), y: String("3.4"), exp: Float(5.6)},
		{x: Float(2.2), y: String("test"), exp: Nil, p: true},
		{x: Float(2.2), y: Bool(true), exp: Float(3.2)},
		{x: Float(2.2), y: Nil, exp: Nil, p: true},
		// String leads
		{x: String("some"), y: Int(45), exp: String("some45")},
		{x: String("some"), y: Int(-5), exp: String("some-5")},
		{x: String("some"), y: Float(2.23), exp: String("some2.23")},
		{x: String("some"), y: Float(2.00023), exp: String("some2.00023")},
		{x: String("some"), y: Float(-0.04), exp: String("some-0.04")},
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
