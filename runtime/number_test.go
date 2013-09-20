package runtime

import (
	"testing"
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
