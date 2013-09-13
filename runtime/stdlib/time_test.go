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
