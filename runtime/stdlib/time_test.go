package stdlib

import (
	"testing"
	"time"

	"github.com/PuerkitoBio/agora/runtime"
)

func TestTimeSleep(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	tm := new(TimeMod)
	tm.SetCtx(ctx)
	n := time.Now()
	tm.time_Sleep(runtime.Number(100))
	if diff := time.Now().Sub(n); diff < 100*time.Millisecond {
		t.Errorf("expected at least 100ms, got %f", diff.Seconds()*1000)
	}
}

func TestTimeNow(t *testing.T) {
	ctx := runtime.NewCtx(nil, nil)
	tm := new(TimeMod)
	tm.SetCtx(ctx)
	exp := time.Now()
	ret := tm.time_Now()
	ob := ret.(runtime.Object)
	if yr := ob.Get(runtime.String("Year")); yr.Int() != int64(exp.Year()) {
		t.Errorf("expected year %d, got %d", exp.Year(), yr.Int())
	}
	if mt := ob.Get(runtime.String("Month")); mt.Int() != int64(exp.Month()) {
		t.Errorf("expected month %d, got %d", exp.Month(), mt.Int())
	}
	if dy := ob.Get(runtime.String("Day")); dy.Int() != int64(exp.Day()) {
		t.Errorf("expected day %d, got %d", exp.Day(), dy.Int())
	}
	if hr := ob.Get(runtime.String("Hour")); hr.Int() != int64(exp.Hour()) {
		t.Errorf("expected hour %d, got %d", exp.Hour(), hr.Int())
	}
	if mn := ob.Get(runtime.String("Minute")); mn.Int() != int64(exp.Minute()) {
		t.Errorf("expected minute %d, got %d", exp.Minute(), mn.Int())
	}
	if sc := ob.Get(runtime.String("Second")); sc.Int() != int64(exp.Second()) {
		t.Errorf("expected second %d, got %d", exp.Second(), sc.Int())
	}
	if ns := ob.Get(runtime.String("Nanosecond")); ns.Int() < int64(exp.Nanosecond()) {
		t.Errorf("expected nanosecond %d, got %d", exp.Nanosecond(), ns.Int())
	}
}

func TestTimeDate(t *testing.T) {
	cases := []struct {
		args []runtime.Val
		exp  time.Time
	}{
		0: {
			args: []runtime.Val{
				runtime.Number(1975),
			},
			exp: time.Date(1975, 1, 1, 0, 0, 0, 0, time.Local),
		},
		1: {
			args: []runtime.Val{
				runtime.Number(1975),
				runtime.Number(2),
			},
			exp: time.Date(1975, 2, 1, 0, 0, 0, 0, time.Local),
		},
		2: {
			args: []runtime.Val{
				runtime.Number(1975),
				runtime.Number(2),
				runtime.Number(3),
			},
			exp: time.Date(1975, 2, 3, 0, 0, 0, 0, time.Local),
		},
		3: {
			args: []runtime.Val{
				runtime.Number(1975),
				runtime.Number(2),
				runtime.Number(3),
				runtime.Number(4),
			},
			exp: time.Date(1975, 2, 3, 4, 0, 0, 0, time.Local),
		},
		4: {
			args: []runtime.Val{
				runtime.Number(1975),
				runtime.Number(2),
				runtime.Number(3),
				runtime.Number(4),
				runtime.Number(5),
			},
			exp: time.Date(1975, 2, 3, 4, 5, 0, 0, time.Local),
		},
		5: {
			args: []runtime.Val{
				runtime.Number(1975),
				runtime.Number(2),
				runtime.Number(3),
				runtime.Number(4),
				runtime.Number(5),
				runtime.Number(6),
			},
			exp: time.Date(1975, 2, 3, 4, 5, 6, 0, time.Local),
		},
		6: {
			args: []runtime.Val{
				runtime.Number(1975),
				runtime.Number(2),
				runtime.Number(3),
				runtime.Number(4),
				runtime.Number(5),
				runtime.Number(6),
				runtime.Number(7),
			},
			exp: time.Date(1975, 2, 3, 4, 5, 6, 7, time.Local),
		},
	}
	ctx := runtime.NewCtx(nil, nil)
	tm := new(TimeMod)
	tm.SetCtx(ctx)
	for i, c := range cases {
		ret := tm.time_Date(c.args...)
		ob := ret.(runtime.Object)
		if yr := ob.Get(runtime.String("Year")); yr.Int() != int64(c.exp.Year()) {
			t.Errorf("[%d] - expected year %d, got %d", i, c.exp.Year(), yr.Int())
		}
		if mt := ob.Get(runtime.String("Month")); mt.Int() != int64(c.exp.Month()) {
			t.Errorf("[%d] - expected month %d, got %d", i, c.exp.Month(), mt.Int())
		}
		if dy := ob.Get(runtime.String("Day")); dy.Int() != int64(c.exp.Day()) {
			t.Errorf("[%d] - expected day %d, got %d", i, c.exp.Day(), dy.Int())
		}
		if hr := ob.Get(runtime.String("Hour")); hr.Int() != int64(c.exp.Hour()) {
			t.Errorf("[%d] - expected hour %d, got %d", i, c.exp.Hour(), hr.Int())
		}
		if mn := ob.Get(runtime.String("Minute")); mn.Int() != int64(c.exp.Minute()) {
			t.Errorf("[%d] - expected minute %d, got %d", i, c.exp.Minute(), mn.Int())
		}
		if sc := ob.Get(runtime.String("Second")); sc.Int() != int64(c.exp.Second()) {
			t.Errorf("[%d] - expected second %d, got %d", i, c.exp.Second(), sc.Int())
		}
		if ns := ob.Get(runtime.String("Nanosecond")); ns.Int() < int64(c.exp.Nanosecond()) {
			t.Errorf("[%d] - expected nanosecond %d, got %d", i, c.exp.Nanosecond(), ns.Int())
		}
	}
}
