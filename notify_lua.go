//go:build darwin || freebsd || netbsd || openbsd || linux
// +build darwin freebsd netbsd openbsd linux

package psnotify

import (
	"github.com/vela-ssoc/vela-kit/lua"
)

func (nt *notify) pipeL(L *lua.LState) int {
	n := L.GetTop()
	if n <= 0 {
		return 0
	}
	for i := 1; i <= n; i++ {
		nt.cfg.pipe.Check(L, i)
	}
	return 0
}

func (nt *notify) startL(L *lua.LState) int {
	xEnv.Start(L, nt).From(nt.Code()).Do()
	return 0
}

func (nt *notify) stateL(L *lua.LState) int {
	nt.cfg.watch.All = true
	n := L.GetTop()
	if n <= 0 {
		nt.cfg.watch.Entry = PROC_EVENT_ALL
		return 0
	}

	for i := 1; i <= n; i++ {
		switch L.CheckString(i) {
		case "*":
			nt.cfg.watch.Entry = PROC_EVENT_ALL
			return 0
		case "fork":
			nt.cfg.watch.Entry |= PROC_EVENT_FORK
		case "exec":
			nt.cfg.watch.Entry |= PROC_EVENT_EXEC
		case "exit":
			nt.cfg.watch.Entry |= PROC_EVENT_EXIT
		}
	}
	return 0
}

func (nt *notify) filterL(L *lua.LState) int {
	nt.cfg.filter.CheckMany(L)
	return 0
}

func (nt *notify) ignoreL(L *lua.LState) int {
	nt.cfg.Ignore.CheckMany(L)
	return 0
}

func (nt *notify) Index(L *lua.LState, key string) lua.LValue {
	switch key {

	case "FORK":
		return lua.LInt(PROC_EVENT_FORK)

	case "EXEC":
		return lua.LInt(PROC_EVENT_EXEC)

	case "EXIT":
		return lua.LInt(PROC_EVENT_EXIT)

	case "state":
		return lua.NewFunction(nt.stateL)

	case "pipe":
		return lua.NewFunction(nt.pipeL)

	case "filter":
		return lua.NewFunction(nt.filterL)
	case "ignore":
		return lua.NewFunction(nt.ignoreL)
	case "case":
		return nt.cfg.vsh.Index(L, "case")

	case "start":
		return lua.NewFunction(nt.startL)

	}

	return lua.LNil
}
