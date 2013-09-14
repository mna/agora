package stdlib

import (
	"time"

	"github.com/PuerkitoBio/agora/runtime"
)

// The time module, as documented in
// https://github.com/PuerkitoBio/agora/wiki/Standard-library
type TimeMod struct {
	ctx *runtime.Ctx
	ob  runtime.Object
}

func (t *TimeMod) ID() string {
	return "time"
}

func (t *TimeMod) Run(_ ...runtime.Val) (v runtime.Val, err error) {
	defer runtime.PanicToError(&err)
	if t.ob == nil {
		// Prepare the object
		t.ob = runtime.NewObject()
		t.ob.Set(runtime.String("Date"), runtime.NewNativeFunc(t.ctx, "time.Date", t.time_Date))
		t.ob.Set(runtime.String("Now"), runtime.NewNativeFunc(t.ctx, "time.Now", t.time_Now))
		t.ob.Set(runtime.String("Sleep"), runtime.NewNativeFunc(t.ctx, "time.Sleep", t.time_Sleep))
	}
	return t.ob, nil
}

func (t *TimeMod) SetCtx(c *runtime.Ctx) {
	t.ctx = c
}

func (t *TimeMod) time_Sleep(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	time.Sleep(time.Duration(args[0].Int()) * time.Millisecond)
	return runtime.Nil
}

type _time struct {
	runtime.Object
	t time.Time
}

func (t *TimeMod) newTime(tm time.Time) runtime.Val {
	ob := &_time{
		runtime.NewObject(),
		tm,
	}
	ob.Set(runtime.String("__toInt"), runtime.NewNativeFunc(t.ctx, "time._time.__toInt", func(args ...runtime.Val) runtime.Val {
		return runtime.Number(ob.t.Unix())
	}))
	ob.Set(runtime.String("__toString"), runtime.NewNativeFunc(t.ctx, "time._time.__toString", func(args ...runtime.Val) runtime.Val {
		return runtime.String(ob.t.Format(time.RFC3339))
	}))
	ob.Set(runtime.String("Year"), runtime.Number(tm.Year()))
	ob.Set(runtime.String("Month"), runtime.Number(tm.Month()))
	ob.Set(runtime.String("Day"), runtime.Number(tm.Day()))
	ob.Set(runtime.String("Hour"), runtime.Number(tm.Hour()))
	ob.Set(runtime.String("Minute"), runtime.Number(tm.Minute()))
	ob.Set(runtime.String("Second"), runtime.Number(tm.Second()))
	ob.Set(runtime.String("Nanosecond"), runtime.Number(tm.Nanosecond()))
	return ob
}

func (t *TimeMod) time_Now(args ...runtime.Val) runtime.Val {
	return t.newTime(time.Now())
}

func (t *TimeMod) time_Date(args ...runtime.Val) runtime.Val {
	runtime.ExpectAtLeastNArgs(1, args)
	yr := int(args[0].Int())
	mth := 1
	if len(args) > 1 {
		mth = int(args[1].Int())
	}
	dy := 1
	if len(args) > 2 {
		dy = int(args[2].Int())
	}
	hr := 0
	if len(args) > 3 {
		hr = int(args[3].Int())
	}
	min := 0
	if len(args) > 4 {
		min = int(args[4].Int())
	}
	sec := 0
	if len(args) > 5 {
		sec = int(args[5].Int())
	}
	nsec := 0
	if len(args) > 6 {
		nsec = int(args[6].Int())
	}
	return t.newTime(time.Date(yr, time.Month(mth), dy, hr, min, sec, nsec, time.UTC))
}
