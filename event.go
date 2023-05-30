//go:build darwin || freebsd || netbsd || openbsd || linux
// +build darwin freebsd netbsd openbsd linux

package psnotify

import (
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-process"
)

type event struct {
	state bool
	eType uint32
	ppid  int
	pid   int
	proc  *process.Process
}

func newEv(ppid, pid int, et uint32) *event {
	return &event{state: false, pid: pid, ppid: ppid, eType: et, proc: &process.Process{Pid: -1}}
}

func (ev *event) eTypeToString() string {
	switch ev.eType {
	case PROC_EVENT_FORK:
		return "fork"
	case PROC_EVENT_EXEC:
		return "exec"
	case PROC_EVENT_EXIT:
		return "exit"
	default:
		return "null"
	}
}

func (ev *event) ps() *process.Process {
	if ev.state {
		return ev.proc
	}

	p, e := process.Pid(ev.pid)
	ev.state = true
	if e != nil {
		xEnv.Infof("pid:%d ppid:%d event:%s got process fail %v", ev.pid, ev.ppid, ev.eTypeToString(), e)
		return ev.proc
	}
	ev.proc = p
	return ev.proc
}

func (ev *event) Byte() []byte {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.KV("type", ev.eTypeToString())
	enc.KV("pid", ev.pid)
	enc.KV("ppid", ev.ppid)
	if ev.proc == nil {
		enc.KV("proc", []byte("{}"))
	} else {
		enc.Raw("proc", ev.proc.Byte())
	}
	enc.End("}")

	return enc.Bytes()
}

func (ev *event) Proc() *process.Process {
	if ev.proc != nil {
		return ev.proc
	}

	var proc *process.Process
	var err error
	if ev.ppid == -1 {
		proc, err = process.Pid(ev.pid)
	} else {
		proc, err = process.Pid(ev.ppid)
	}

	if err != nil {
		return nil
	}

	ev.proc = proc
	return proc
}

func (ev *event) info() []byte {
	enc := kind.NewJsonEncoder()
	p := ev.ps()

	enc.Tab("")
	enc.KV("type", ev.eTypeToString())
	enc.KV("pid", ev.pid)
	enc.KV("ppid", ev.ppid)
	enc.KV("gid", p.Pgid)
	enc.KV("name", p.Name)
	enc.KV("exe", p.Executable)
	enc.KV("cmdline", p.Cmdline)
	enc.KV("cwd", p.Cwd)
	enc.KV("user", p.Username)
	enc.KV("args", p.ArgsToString())
	enc.End("}")
	return enc.Bytes()
}
