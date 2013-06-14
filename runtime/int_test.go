package runtime

import (
	"testing"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		x   int
		y   int
		exp int
	}{
		{x: 0, y: 0, exp: 0},
		{x: 0, y: 1, exp: 1},
		{x: 1, y: 0, exp: 1},
		{x: 2, y: 5, exp: 7},
		{x: -12, y: 356, exp: 344},
		{x: -1, y: 0, exp: -1},
		{x: -1, y: 1, exp: 0},
		{x: -1, y: -1, exp: -2},
		{x: 4294967296, y: 1, exp: 4294967297}, // Would fail on 32-bit systems
		{x: 1000, y: -100, exp: 900},
	}

	for _, c := range cases {
		vx, vy := Int(c.x), Int(c.y)
		res := vx.Add(vy)
		if ires := int(res.(Int)); c.exp != ires {
			t.Errorf("%d + %d : expected %d, got %d", c.x, c.y, c.exp, ires)
		}
	}
}
