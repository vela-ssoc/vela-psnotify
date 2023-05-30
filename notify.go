//go:build darwin || freebsd || netbsd || openbsd || linux
// +build darwin freebsd netbsd openbsd linux

package psnotify

import (
	"fmt"
	"github.com/vela-ssoc/vela-kit/lua"
	"gopkg.in/tomb.v2"
	"reflect"
	"time"
)

var typeof = reflect.TypeOf((*notify)(nil)).String()

type notify struct {
	lua.SuperVelaData
	tom   *tomb.Tomb
	cfg   *config
	watch *Watcher
}

func newNotify(cfg *config) *notify {
	nty := &notify{cfg: cfg}
	nty.V(lua.VTInit, typeof, time.Now())
	return nty
}

func (nt *notify) Code() string {
	return nt.cfg.co.CodeVM()
}

func (nt *notify) E(er error) {
	xEnv.Errorf("%s pipe vela-event %v", nt.Name(), er)
}

func (nt *notify) handle(ev *event) {
	if nt.cfg.Ignore.Match(ev) {
		return
	}

	if !nt.cfg.filter.Match(ev) {
		return
	}

	nt.cfg.vsh.Do(ev)
	nt.cfg.pipe.Do(ev, nt.cfg.co, nt.E)
}

func (nt *notify) accept() {
	wer := nt.cfg.watch
	defer func() {
		_ = wer.Close()
	}()

	for {
		select {

		case ev := <-wer.Fork:
			nt.handle(newEv(ev.ParentPid, ev.ChildPid, PROC_EVENT_FORK))

		case ev := <-wer.Exec:
			nt.handle(newEv(-1, ev.Pid, PROC_EVENT_EXEC))

		case ev := <-wer.Exit:
			nt.handle(newEv(-1, ev.Pid, PROC_EVENT_EXIT))

		case ev := <-wer.Error:
			xEnv.Infof("%s linux netlink notify got error %v", nt.Name(), ev.Error())

		case <-nt.tom.Dying():
			return
		}
	}
}

func (nt *notify) Name() string {
	return nt.cfg.name
}

func (nt *notify) Type() string {
	return typeof
}
func (nt *notify) Start() error {
	nt.tom = new(tomb.Tomb)

	nt.tom.Go(func() error {
		nt.accept()
		return nil
	})

	return nil
}

func (nt *notify) Close() error {
	defer func() {
		nt.V(lua.VTClose)
	}()

	nt.tom.Kill(fmt.Errorf("close"))
	return nil
}
