package build

import (
	"errors"
	"io/fs"
	"merodi/src/core/build/compile"
	"merodi/src/core/state"
	"merodi/src/parser"
	"merodi/src/utils"
	path "merodi/src/utils/path"
	"strings"

	. "github.com/pixel2175/errs"
)

type Build struct {
	state   *state.State
	compile *compile.Compile
}

func (self *Build) InitCompile() {
	self.compile = &compile.Compile{State: self.state}
	TryE(self.compile.Init())
	self.compile.Build = &self.state.Build
}

func (self *Build) Visit(filepath string, entry fs.DirEntry, err error) (rerr error) {
	defer Catch(&rerr, func(err *error) {
		if errors.Is(*err, utils.ErrAborted) {
			*err = self.state.RunHook("on_abort")
		}
	})

	TryE(err)
	if !self.isSource(entry) {
		return
	}

	if strings.HasSuffix(filepath, ".md") {
		TryE(self.compile.Run(filepath))
	} else {
		self.compile.ResolvePaths(filepath)
		TryE(path.CopyFile(self.compile.Build.MDPath, self.compile.Build.DocPath))
	}

	return
}

func (self *Build) Run(state *state.State, opt *parser.Option) (err error) {
	defer Catch(&err)

	self.state = state
	self.InitCompile()

	TryE(self.walkAndBuild())
	return
}
