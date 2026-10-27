package clean

import (
	"merodi/src/config"
	"merodi/src/core/state"
	"merodi/src/parser"
	. "github.com/pixel2175/errs"
	"os"
)

type Clean struct {
	State *state.State
}

func (self *Clean) Run(state *state.State, opt *parser.Option) (err error) {
	defer Catch(&err)
	self.State = state

	if self.State.Build.Mode == config.Release {
		TryE(os.RemoveAll(self.State.Config.Tree.ReleaseDest))
	} else {
		TryE(os.RemoveAll(self.State.Config.Tree.DraftDest))
	}
	return
}
