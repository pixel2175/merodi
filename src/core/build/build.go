package build

import (
	"errors"
	"io/fs"
	"merodi/src/core/build/compile"
	"merodi/src/core/state"
	"merodi/src/parser"
	. "merodi/src/utils/errs"
	path "merodi/src/utils/path"
	"strings"
)

type Build struct {
	state   *state.State
	compile *compile.Compile
}

func (self *Build) InitCompile() {
	self.compile = &compile.Compile{State: self.state}
	CheckE(self.compile.Init())
	self.compile.Build = &self.state.Build
}

func (self *Build) Visit(filepath string, entry fs.DirEntry, err error) (rerr error) {
	defer Handle(&rerr, func(err *error) {
		if errors.Is(*err, ErrAborted) {
			*err = self.state.RunHook("on_abort")
		}
	})

	CheckE(err)
	if !self.isSource(entry) {
		return
	}

	if strings.HasSuffix(filepath, ".md") {
		CheckE(self.compile.Run(filepath))
	} else {
		self.compile.ResolvePaths(filepath)
		CheckE(path.CopyFile(self.compile.Build.MDPath, self.compile.Build.DocPath))
	}

	return
}

func (self *Build) Run(state *state.State, opt *parser.Option) (err error) {
	defer Handle(&err)

	self.state = state
	self.InitCompile()

	CheckE(self.walkAndBuild())
	return
}
