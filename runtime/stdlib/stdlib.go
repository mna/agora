package stdlib

import (
	"github.com/PuerkitoBio/goblin/runtime"
)

var (
	lib = make(map[string]map[string]func(*runtime.Ctx, ...Val) Val)
)

func RegisterPackage(ctx *runtime.Ctx, pkg string) {

}
