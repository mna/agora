package runtime

import (
	"testing"
)

func TestIntAsInt(t *testing.T) {
	cases := []struct {
		x   int
		exp int
	}{
		{x: 0, exp: 0},
		{x: 1, exp: 1},
		{x: -1, exp: -1},
		{x: 123, exp: 123},
		{x: -999, exp: -999},
	}

	for _, c := range cases {
		vx := Int(c.x)
		res := vx.Int()
		if c.exp != res {
			t.Errorf("%d as int : expected %d, got %d", c.x, c.exp, res)
		}
	}
}

func TestIntAsFloat(t *testing.T) {
	cases := []struct {
		x   int
		exp float64
	}{
		{x: 0, exp: 0.0},
		{x: 1, exp: 1.0},
		{x: -1, exp: -1.0},
		{x: 123, exp: 123.0},
		{x: -999, exp: -999.0},
	}

	for _, c := range cases {
		vx := Int(c.x)
		res := vx.Float()
		if c.exp != res {
			t.Errorf("%d as float : expected %f, got %f", c.x, c.exp, res)
		}
	}
}

func TestIntAsString(t *testing.T) {
	cases := []struct {
		x   int
		exp string
	}{
		{x: 0, exp: "0"},
		{x: 1, exp: "1"},
		{x: -1, exp: "-1"},
		{x: 123, exp: "123"},
		{x: -999, exp: "-999"},
	}

	for _, c := range cases {
		vx := Int(c.x)
		res := vx.String()
		if c.exp != res {
			t.Errorf("%d as string : expected %s, got %s", c.x, c.exp, res)
		}
	}
}

func TestIntAsBool(t *testing.T) {
	cases := []struct {
		x   int
		exp bool
	}{
		{x: 0, exp: false},
		{x: 1, exp: true},
		{x: -1, exp: true},
		{x: 123, exp: true},
		{x: -999, exp: true},
	}

	for _, c := range cases {
		vx := Int(c.x)
		res := vx.Bool()
		if c.exp != res {
			t.Errorf("%d as bool : expected %v, got %v", c.x, c.exp, res)
		}
	}
}

func TestAddInt(t *testing.T) {
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

func TestSubInt(t *testing.T) {
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

func TestMulInt(t *testing.T) {
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

func TestDivInt(t *testing.T) {
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

func TestModInt(t *testing.T) {
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
			t.Errorf("%d %% %d : expected %d, got %d", c.x, c.y, c.exp, ires)
		}
	}
}

func TestUnmInt(t *testing.T) {
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
