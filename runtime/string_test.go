package runtime

import (
	"testing"
)

func TestStringAsInt(t *testing.T) {
	cases := []struct {
		x   string
		exp int64
		p   bool
	}{
		{x: "0", exp: 0},
		{x: "1", exp: 1},
		{x: "-1", exp: -1},
		{x: "123", exp: 123},
		{x: "-999", exp: -999},
		{x: "-999.23", exp: 0, p: true},
		{x: "a9", exp: 0, p: true},
		{x: "", exp: 0, p: true},
	}

	assert := func(s string, p bool) {
		if err := recover(); (err == nil && p) || (err != nil && !p) {
			t.Errorf("%s : expected %v, got '%s'", s, p, err)
		}
	}
	for _, c := range cases {
		func() {
			defer assert(c.x, c.p)
			vx := String(c.x)
			res := vx.Int()
			if c.p {
				t.Errorf("%s : expected a panic", c.x)
			}
			if c.exp != res {
				t.Errorf("%s as int : expected %d, got %d", c.x, c.exp, res)
			}
		}()
	}
}

func TestStringAsFloat(t *testing.T) {
	cases := []struct {
		x   string
		exp float64
		p   bool
	}{
		{x: "0.000", exp: 0.0},
		{x: "1", exp: 1.0},
		{x: "-1", exp: -1.0},
		{x: "123", exp: 123.0},
		{x: "123.0000", exp: 123.0},
		{x: "-999.00000", exp: -999.0},
		{x: "-999.23", exp: -999.23},
		{x: "1999.023", exp: 1999.023},
		{x: "1e2", exp: 100},
		{x: "a9", exp: 0, p: true},
		{x: "", exp: 0, p: true},
	}

	assert := func(s string, p bool) {
		if err := recover(); (err == nil && p) || (err != nil && !p) {
			t.Errorf("%s : expected %v, got '%s'", s, p, err)
		}
	}
	for _, c := range cases {
		func() {
			defer assert(c.x, c.p)
			vx := String(c.x)
			res := vx.Float()
			if c.p {
				t.Errorf("%s : expected a panic", c.x)
			}
			if c.exp != res {
				t.Errorf("%s as float : expected %f, got %f", c.x, c.exp, res)
			}
		}()
	}
}

func TestStringAsString(t *testing.T) {
	cases := []struct {
		x   string
		exp string
	}{
		{x: "", exp: ""},
		{x: " ", exp: " "},
		{x: "\n", exp: "\n"},
		{x: "testpatatepoil", exp: "testpatatepoil"},
		{x: "123.0000", exp: "123.0000"},
	}

	for _, c := range cases {
		vx := String(c.x)
		res := vx.String()
		if c.exp != res {
			t.Errorf("%s as string : expected %s, got %s", c.x, c.exp, res)
		}
	}
}

func TestStringAsBool(t *testing.T) {
	cases := []struct {
		x   string
		exp bool
	}{
		{x: "", exp: false},
		{x: " ", exp: true},
		{x: "\n", exp: true},
		{x: "testpatatepoil", exp: true},
		{x: "123.0000", exp: true},
	}

	for _, c := range cases {
		vx := String(c.x)
		res := vx.Bool()
		if c.exp != res {
			t.Errorf("%s as bool : expected %v, got %v", c.x, c.exp, res)
		}
	}
}

func TestAddString(t *testing.T) {
	ari := defaultArithmetic{}
	cases := []struct {
		x   string
		y   string
		exp string
	}{
		{x: "", y: "", exp: ""},
		{x: " ", y: "", exp: " "},
		{x: "", y: " ", exp: " "},
		{x: " ", y: " ", exp: "  "},
		{x: "abc", y: "def", exp: "abcdef"},
		{x: "abc\ndef\t", y: "ghi\njkl\x00m", exp: "abc\ndef\tghi\njkl\x00m"},
	}

	for _, c := range cases {
		vx, vy := String(c.x), String(c.y)
		res := ari.Add(vx, vy)
		if sres := res.String(); c.exp != sres {
			t.Errorf("%s + %s : expected %s, got %s", c.x, c.y, c.exp, sres)
		}
	}
}
