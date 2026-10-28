package lua

import (
	lua "github.com/yuin/gopher-lua"
)

func (self *Lua) registerConfigTree() *lua.LTable {
	tree := self.State.Lua.Context.NewTable()
	tree.RawSetString("markdown", lua.LString(self.State.Config.Tree.Markdown))
	tree.RawSetString("templates", lua.LString(self.State.Config.Tree.Templates))
	tree.RawSetString("plugins", lua.LString(self.State.Config.Tree.Plugins))
	tree.RawSetString("draft_dest", lua.LString(self.State.Config.Tree.DraftDest))
	tree.RawSetString("release_dest", lua.LString(self.State.Config.Tree.ReleaseDest))
	return tree
}

func (self *Lua) registerConfigProject() *lua.LTable {
	projectTable := self.State.Lua.Context.NewTable()
	projectTable.RawSetString("name", lua.LString(self.State.Config.Project.Name))
	projectTable.RawSetString("version", lua.LString(self.State.Config.Project.Version))
	projectTable.RawSetString("description", lua.LString(self.State.Config.Project.Description))
	return projectTable
}

func (self *Lua) registerConfig() {
	configTable := self.State.Lua.Context.NewTable()

	configTable.RawSetString("tree", self.registerConfigTree())
	configTable.RawSetString("project", self.registerConfigProject())

	self.Merodi.RawSetString("config", configTable)
}
