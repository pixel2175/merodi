package hook

import (
	"merodi/src/core/state"
	"merodi/src/parser"
	. "github.com/pixel2175/errs"
)

type Hook struct {
	State *state.State
}

func (self *Hook) Run(state *state.State, opt *parser.Option) (err error) {
	defer Catch(&err)
	self.State = state

	TryE(self.State.RunHook(opt.Hook, opt.Args...))
	return
}
