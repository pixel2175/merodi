package lua

import (
	"merodi/src/core/state"
	. "github.com/pixel2175/errs"
	"path/filepath"

	glua "github.com/yuin/gopher-lua"
)

type Lua struct {
	State   *state.State
	Context *glua.LState
	Merodi  *glua.LTable
}

func Init(st *state.State) (err error) {
	defer Catch(&err)

	st.Lua.Context = glua.NewState()

	self := Lua{
		State:   st,
		Context: st.Lua.Context,
	}

	self.State.Lua.Aborted = false

	self.registerMerodi()

	path := filepath.Join(self.State.Config.Tree.Plugins, "main.lua")
	TryE(self.loadFile(path))
	return
}

func (self *Lua) loadFile(path string) (err error) {
	defer Catch(&err)
	fn, lerr := self.State.Lua.Context.LoadFile(path)
	if lerr != nil {
		TryE(describeLuaError(path, lerr))
	}

	self.State.Lua.Context.Push(fn)
	if perr := self.State.Lua.Context.PCall(0, glua.MultRet, nil); perr != nil {
		TryE(describeLuaError(path, perr))
	}
	return
}

func (self *Lua) registerMerodi() {
	self.Merodi = self.State.Lua.Context.NewTable()
	self.registerHooks()
	self.registerConfig()
	self.registerJinja()
	self.registerBuild()
	self.registerWatch()
	self.registerEnable()
	self.registerDisable()
	self.registerLog()
	self.registerCompile()
	self.registerAction()
	self.State.Lua.Context.SetGlobal("merodi", self.Merodi)
}
