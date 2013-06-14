package runtime

import (
	"testing"
)

// TODO Test conversions too

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

func TestSub(t *testing.T) {
	cases := []struct {
		x   int
		y   int
		exp int
	}{
		{x: 0, y: 0, exp: 0},
		{x: 0, y: 1, exp: -1},
		{x: 1, y: 0, exp: 1},
		{x: 2, y: 5, exp: -3},
		{x: -12, y: 356, exp: -368},
		{x: -1, y: 0, exp: -1},
		{x: -1, y: 1, exp: -2},
		{x: -1, y: -1, exp: 0},
		{x: 4294967296, y: 1, exp: 4294967295},
		{x: 1000, y: -100, exp: 1100},
	}

	for _, c := range cases {
		vx, vy := Int(c.x), Int(c.y)
		res := vx.Sub(vy)
		if ires := int(res.(Int)); c.exp != ires {
			t.Errorf("%d - %d : expected %d, got %d", c.x, c.y, c.exp, ires)
		}
	}
}

func TestMul(t *testing.T) {
	cases := []struct {
		x   int
		y   int
		exp int
	}{
		{x: 0, y: 0, exp: 0},
		{x: 0, y: 1, exp: 0},
		{x: 1, y: 0, exp: 0},
		{x: 2, y: 5, exp: 10},
		{x: -12, y: 356, exp: -4272},
		{x: -1, y: 0, exp: 0},
		{x: -1, y: 1, exp: -1},
		{x: -1, y: -1, exp: 1},
		{x: 4294967296, y: 1, exp: 4294967296},
		{x: 1000, y: -100, exp: -100000},
	}

	for _, c := range cases {
		vx, vy := Int(c.x), Int(c.y)
		res := vx.Mul(vy)
		if ires := int(res.(Int)); c.exp != ires {
			t.Errorf("%d * %d : expected %d, got %d", c.x, c.y, c.exp, ires)
		}
	}
}

func TestDiv(t *testing.T) {
	cases := []struct {
		x   int
		y   int
		exp int
	}{
		{x: 0, y: 1, exp: 0},
		{x: 20, y: 5, exp: 4},
		{x: -12, y: 356, exp: 0},
		{x: -1, y: 1, exp: -1},
		{x: -1, y: -1, exp: 1},
		{x: 4294967296, y: 1, exp: 4294967296},
		{x: 1000, y: -100, exp: -10},
		{x: 10, y: 3, exp: 3},
	}

	for _, c := range cases {
		vx, vy := Int(c.x), Int(c.y)
		res := vx.Div(vy)
		if ires := int(res.(Int)); c.exp != ires {
			t.Errorf("%d / %d : expected %d, got %d", c.x, c.y, c.exp, ires)
		}
	}
}

func TestMod(t *testing.T) {
	cases := []struct {
		x   int
		y   int
		exp int
	}{
		{x: 0, y: 1, exp: 0},
		{x: 20, y: 5, exp: 0},
		{x: -12, y: 356, exp: -12},
		{x: -1, y: 1, exp: 0},
		{x: -1, y: -1, exp: 0},
		{x: 4294967296, y: 1, exp: 0},
		{x: 1000, y: -100, exp: 0},
		{x: 10, y: 3, exp: 1},
	}

	for _, c := range cases {
		vx, vy := Int(c.x), Int(c.y)
		res := vx.Mod(vy)
		if ires := int(res.(Int)); c.exp != ires {
			t.Errorf("%d % %d : expected %d, got %d", c.x, c.y, c.exp, ires)
		}
	}
}

func TestPow(t *testing.T) {
	cases := []struct {
		x   int
		y   int
		exp int
	}{
		{x: 0, y: 1, exp: 0},
		{x: 20, y: 5, exp: 3200000},
		{x: -12, y: 4, exp: 20736},
		{x: -1, y: 1, exp: -1},
		{x: -1, y: -1, exp: -1},
		{x: 4294967296, y: 1, exp: 4294967296},
		{x: 1000, y: -100, exp: 0},
		{x: 10, y: 3, exp: 1000},
	}

	for _, c := range cases {
		vx, vy := Int(c.x), Int(c.y)
		res := vx.Pow(vy)
		if ires := int(res.(Int)); c.exp != ires {
			t.Errorf("%d ^ %d : expected %d, got %d", c.x, c.y, c.exp, ires)
		}
	}
}

func TestNot(t *testing.T) {
	cases := []struct {
		x   int
		exp int
	}{
		{x: 0, exp: -1},
	}

	for _, c := range cases {
		vx := Int(c.x)
		res := vx.Not()
		if ires := int(res.(Int)); c.exp != ires {
			t.Errorf("!%d : expected %d, got %d", c.x, c.exp, ires)
		}
	}
}

func TestUnm(t *testing.T) {
	cases := []struct {
		x   int
		exp int
	}{
		{x: 0, exp: 0},
		{x: 1, exp: -1},
		{x: -1, exp: 1},
		{x: -12, exp: 12},
		{x: 234, exp: -234},
	}

	for _, c := range cases {
		vx := Int(c.x)
		res := vx.Unm()
		if ires := int(res.(Int)); c.exp != ires {
			t.Errorf("-%d : expected %d, got %d", c.x, c.exp, ires)
		}
	}
}
