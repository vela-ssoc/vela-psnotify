//go:build darwin || freebsd || netbsd || openbsd || linux
// +build darwin freebsd netbsd openbsd linux

package psnotify

import (
	cond "github.com/vela-ssoc/vela-cond"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	vswitch "github.com/vela-ssoc/vela-switch"
)

type config struct {
	name   string
	Ignore *cond.Ignore
	filter *cond.Combine
	pipe   *pipe.Px
	vsh    *vswitch.Switch
	co     *lua.LState
	watch  *Watcher
}

func newConfig(L *lua.LState) *config {
	val := L.Get(1)
	cfg := &config{
		co:     xEnv.Clone(L),
		pipe:   pipe.New(pipe.Env(xEnv)),
		Ignore: cond.NewIgnore(),
		filter: cond.NewCombine(),
		vsh:    vswitch.NewL(L),
	}

	switch val.Type() {
	case lua.LTString:
		cfg.name = val.String()

	case lua.LTTable:
		tab := val.(*lua.LTable)
		tab.Range(func(key string, val lua.LValue) {
			cfg.NewIndex(L, key, val)
		})
	}

	if e := cfg.valid(); e != nil {
		L.RaiseError("%v", e)
		return nil
	}

	return cfg
}

func (cfg *config) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "name":
		cfg.name = val.String()
	default:
		//todo
	}
}

func (cfg *config) valid() error {
	if e := auxlib.Name(cfg.name); e != nil {
		return e
	}

	wer, err := NewWatcher()
	if err != nil {
		return err
	}

	cfg.watch = wer
	return nil
}
