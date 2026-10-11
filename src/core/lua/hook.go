package lua

import (
	. "merodi/src/utils/errs"

	lua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerHook(L *lua.LState) int {
	stage := L.ToString(1)
	fn := L.ToFunction(2)
	self.State.Hooks[stage] = fn
	return 0
}

func (self *Lua) luaAbort(L *lua.LState) int {
	self.State.Lua.Aborted = true
	panic(ErrAborted)
}

func (self *Lua) registerHooks() {
	self.Context.SetFuncs(self.Merodi, map[string]lua.LGFunction{
		"hook":  self.registerHook,
		"abort": self.luaAbort,
	})
}
