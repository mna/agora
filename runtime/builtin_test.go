package runtime

import (
	"testing"
)

func TestLen(t *testing.T) {
	cases := []struct {
		src Val
		exp int64
	}{
		0: {
			src: Nil,
			exp: 0,
		},
		1: {
			src: Number(3.14),
			exp: 4,
		},
		2: {
			src: String("hi, there"),
			exp: 9,
		},
		3: {
			src: Bool(true),
			exp: 4,
		},
		4: {
			src: String(`this
has
new
lines`),
			exp: 18,
		},
		5: {
			src: &object{
				map[Val]Val{
					Number(1):      String("val1"),
					String("name"): Bool(false),
					String("subobj"): &object{
						map[Val]Val{
							String("key"): Number(10),
						},
					},
				},
			},
			exp: 3,
		},
		6: {
			src: &object{},
			exp: 0,
		},
		7: {
			src: String(""),
			exp: 0,
		},
	}

	bi := new(builtinMod)
	ctx := NewCtx(nil, nil)
	bi.SetCtx(ctx)
	for i, c := range cases {
		ret := bi._len(c.src)
		if c.exp != ret.Int() {
			t.Errorf("[%d] - expected %d, got %d", i, c.exp, ret.Int())
		}
	}
}
