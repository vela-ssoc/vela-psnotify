package psnotify

import (
	"github.com/vela-ssoc/vela-kit/vela"
)

func WithEnv(env vela.Environment) {
	env.Error("not support psnotify with linux")
}
