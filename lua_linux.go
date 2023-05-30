package psnotify

import (
	"github.com/vela-ssoc/vela-kit/vela"
	"github.com/vela-ssoc/vela-kit/lua"
)

var xEnv vela.Environment

func constructor(L *lua.LState) int {
	cfg := newConfig(L)

	proc := L.NewVelaData(cfg.name, typeof)
	if proc.IsNil() {
		proc.Set(newNotify(cfg))
	} else {
		nty := proc.Data.(*notify)
		xEnv.Free(nty.cfg.co)
		nty.cfg = cfg
	}

	L.Push(proc)
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env
	xEnv.Set("psnotify", lua.NewFunction(constructor))
}
