package runtime

import (
	"math"
	"testing"
)

const (
	floatCompareBuffer = 1e-6
)

func TestNumberAsInt(t *testing.T) {
	cases := []struct {
		x   float64
		exp int64
	}{
		{x: 0.0, exp: 0},
		{x: 1.0, exp: 1},
		{x: -1.0, exp: -1},
		{x: 0.00001, exp: 0},
		{x: 0.99999, exp: 0},
		{x: 1.9, exp: 1},
		{x: 123.456789, exp: 123},
		{x: -1.987654321, exp: -1},
		{x: -999.999999, exp: -999},
	}

	for _, c := range cases {
		vx := Number(c.x)
		res := vx.Int()
		if c.exp != res {
			t.Errorf("%f as int : expected %d, got %d", c.x, c.exp, res)
		}
	}
}

func TestNumberAsFloat(t *testing.T) {
	cases := []struct {
		x   float64
		exp float64
	}{
		{x: 0.0, exp: 0.0},
		{x: 1.0, exp: 1.0},
		{x: -1.0, exp: -1.0},
		{x: 0.00001, exp: 0.00001},
		{x: 0.99999, exp: 0.99999},
		{x: 1.9, exp: 1.9},
		{x: 123.456789, exp: 123.456789},
		{x: -1.987654321, exp: -1.987654321},
		{x: -999.999999, exp: -999.999999},
	}

	for _, c := range cases {
		vx := Number(c.x)
		res := vx.Float()
		if c.exp != res {
			t.Errorf("%f as float : expected %f, got %f", c.x, c.exp, res)
		}
	}
}

func TestNumberAsString(t *testing.T) {
	cases := []struct {
		x   float64
		exp string
	}{
		{x: 0.0, exp: "0"},
		{x: 1.0, exp: "1"},
		{x: -1.0, exp: "-1"},
		{x: 0.00001, exp: "0.00001"},
		{x: 0.99999, exp: "0.99999"},
		{x: 1.9, exp: "1.9"},
		{x: 123.456789, exp: "123.456789"},
		{x: -1.987654321, exp: "-1.987654321"},
		{x: -999.999999, exp: "-999.999999"},
	}

	for _, c := range cases {
		vx := Number(c.x)
		res := vx.String()
		if c.exp != res {
			t.Errorf("%f as string : expected %s, got %s", c.x, c.exp, res)
		}
	}
}

func TestNumberAsBool(t *testing.T) {
	cases := []struct {
		x   float64
		exp bool
	}{
		{x: 0.0, exp: false},
		{x: 1.0, exp: true},
		{x: -1.0, exp: true},
		{x: 0.00001, exp: true},
		{x: 0.99999, exp: true},
		{x: 1.9, exp: true},
		{x: 123.456789, exp: true},
		{x: -1.987654321, exp: true},
		{x: -999.999999, exp: true},
	}

	for _, c := range cases {
		vx := Number(c.x)
		res := vx.Bool()
		if c.exp != res {
			t.Errorf("%f as bool : expected %v, got %v", c.x, c.exp, res)
		}
	}
}

func TestAddNumber(t *testing.T) {
	cases := []struct {
		x   float64
		y   float64
		exp float64
	}{
		{x: 0.0, y: 0.0, exp: 0.0},
		{x: 1.0, y: 0.0, exp: 1.0},
		{x: 0.0, y: 1.0, exp: 1.0},
		{x: 1.0, y: 1.0, exp: 2.0},
		{x: 1.1, y: 0.9, exp: 2.0},
		{x: -10.90, y: 1.1, exp: -9.8},
		{x: 10.123, y: 9.456, exp: 19.579},
	}

	for _, c := range cases {
		vx, vy := Number(c.x), Number(c.y)
		res := vx.Add(vy)
		if fres := res.Float(); math.Abs(c.exp-fres) > floatCompareBuffer {
			t.Errorf("%f + %f : expected %f, got %f", c.x, c.y, c.exp, res.Float())
		}
	}
}

func TestSubNumber(t *testing.T) {
	cases := []struct {
		x   float64
		y   float64
		exp float64
	}{
		{x: 0.0, y: 0.0, exp: 0.0},
		{x: 1.0, y: 0.0, exp: 1.0},
		{x: 0.0, y: 1.0, exp: -1.0},
		{x: 1.0, y: 1.0, exp: 0.0},
		{x: 1.1, y: 0.9, exp: 0.2},
		{x: -10.90, y: 1.1, exp: -12.0},
		{x: 10.123, y: 9.456, exp: 0.667},
	}

	for _, c := range cases {
		vx, vy := Number(c.x), Number(c.y)
		res := vx.Sub(vy)
		if fres := res.Float(); math.Abs(c.exp-fres) > floatCompareBuffer {
			t.Errorf("%f - %f : expected %f, got %f", c.x, c.y, c.exp, fres)
		}
	}
}

func TestMulNumber(t *testing.T) {
	cases := []struct {
		x   float64
		y   float64
		exp float64
	}{
		{x: 0.0, y: 0.0, exp: 0.0},
		{x: 1.0, y: 0.0, exp: 0.0},
		{x: 0.0, y: 1.0, exp: 0.0},
		{x: 1.0, y: 1.0, exp: 1.0},
		{x: 1.1, y: 0.9, exp: 0.99},
		{x: -10.90, y: 1.1, exp: -11.99},
		{x: 10.123, y: 9.456, exp: 95.723088},
	}

	for _, c := range cases {
		vx, vy := Number(c.x), Number(c.y)
		res := vx.Mul(vy)
		if fres := res.Float(); math.Abs(c.exp-fres) > floatCompareBuffer {
			t.Errorf("%f * %f : expected %f, got %f", c.x, c.y, c.exp, fres)
		}
	}
}

func TestDivNumber(t *testing.T) {
	cases := []struct {
		x   float64
		y   float64
		exp float64
	}{
		{x: 0.0, y: 1.0, exp: 0.0},
		{x: 1.0, y: 1.0, exp: 1.0},
		{x: 1.1, y: 0.9, exp: 1.222222222},
		{x: -10.90, y: 1.1, exp: -9.909090909},
		{x: 10.123, y: 9.456, exp: 1.070537225},
	}

	for _, c := range cases {
		vx, vy := Number(c.x), Number(c.y)
		res := vx.Div(vy)
		if fres := res.Float(); math.Abs(c.exp-fres) > floatCompareBuffer {
			t.Errorf("%f / %f : expected %f, got %f", c.x, c.y, c.exp, fres)
		}
	}
}

func TestModNumber(t *testing.T) {
	cases := []struct {
		x   float64
		y   float64
		exp float64
	}{
		{x: 0.0, y: 1.0, exp: 0.0},
		{x: 1.0, y: 1.0, exp: 0.0},
		{x: 1.1, y: 0.9, exp: 0.2},
		{x: -10.90, y: 1.1, exp: -1.0},
		{x: 10.123, y: 9.456, exp: 0.667},
	}

	for _, c := range cases {
		vx, vy := Number(c.x), Number(c.y)
		res := vx.Mod(vy)
		if fres := res.Float(); math.Abs(c.exp-fres) > floatCompareBuffer {
			t.Errorf("%f %% %f : expected %f, got %f", c.x, c.y, c.exp, fres)
		}
	}
}

func TestUnmNumber(t *testing.T) {
	cases := []struct {
		x   float64
		exp float64
	}{
		{x: 0.0, exp: 0.0},
		{x: 1.0, exp: -1.0},
		{x: 1.1, exp: -1.1},
		{x: -10.90, exp: 10.90},
		{x: 10.123, exp: -10.123},
	}

	for _, c := range cases {
		vx := Number(c.x)
		res := vx.Unm()
		if fres := res.Float(); math.Abs(c.exp-fres) > floatCompareBuffer {
			t.Errorf("-%f : expected %f, got %f", c.x, c.exp, fres)
		}
	}
}
