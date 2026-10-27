package lua

import (
	"merodi/src/core/build/compile"
	"merodi/src/core/state"
	. "github.com/pixel2175/errs"
	"strings"

	glua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerCompile() {
	c := &compile.Compile{State: self.State}
	build := state.Build{}
	TryE(c.Init())
	c.Build = &build

	compileTable := self.Context.NewTable()

	compileTable.RawSetString("convert", self.Context.NewFunction(func(L *glua.LState) int {
		c.Build.Content = L.ToString(1)
		L.Push(glua.LString(strings.TrimRight(c.Convert(), "\n")))
		return 1
	}))

	compileTable.RawSetString("mdtohtml", self.Context.NewFunction(func(L *glua.LState) int {
		md := L.ToString(1)
		html := TryV(c.MD.MdToDoc(md))
		L.Push(glua.LString(strings.TrimRight(html, "\n")))
		return 1
	}))

	self.Merodi.RawSetString("compile", compileTable)
}
