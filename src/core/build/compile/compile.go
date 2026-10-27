package compile

import (
	. "merodi/src/core/build/compile/jinja"
	. "merodi/src/core/build/compile/markdown"
	"merodi/src/core/state"
	. "github.com/pixel2175/errs"
)

type Compile struct {
	State *state.State
	Jinja JinjaEngine
	MD    Markdown
	Build *state.Build
}

func (self *Compile) Init() (err error) {
	defer Catch(&err)

	self.Jinja = JinjaEngine{Title: self.State.Config.Project.Description}
	TryE(self.Jinja.Init(self.State.Config.Tree.Templates))

	self.MD = Markdown{}
	self.MD.Init(self.State)

	return nil
}

func (self *Compile) Convert() (doccontent string) {
	self.applyJinja()
	self.convertToDoc()
	return self.Build.Content
}

func (self *Compile) Run(mdfile string) (err error) {
	defer Catch(&err)
	defer self.reset()

	self.ResolvePaths(mdfile)
	self.readMarkdown()

	doc := self.Convert()
	self.hook("before_write")
	TryE(self.WriteDoc(doc))
	return
}
