package core

import (
	"merodi/src/core/build"
	"merodi/src/core/clean"
	"merodi/src/core/hook"
	"merodi/src/core/lua"
	. "merodi/src/core/state"
	"merodi/src/core/watch"
	"merodi/src/parser"
	. "github.com/pixel2175/errs"
)

type Action interface {
	Run(state *State, opt *parser.Option) error
}

var Actions = map[parser.Action]Action{
	parser.Build: &build.Build{},
	parser.Watch: &watch.Watcher{},
	parser.Clean: &clean.Clean{},
	parser.Hook:  &hook.Hook{},
}

func Run(opt *parser.Option) (err error) {
	defer Catch(&err)
	state := InitGlobals()
	state.Action = opt.Action
	state.Build.Mode = opt.Mode

	TryE(state.GoToProjectDir(opt.ProjectDir))
	TryE(state.LoadConfig())
	if opt.Action != parser.Clean {
		TryE(lua.Init(&state))
		defer state.Lua.Context.Close()
	}
	TryE(Actions[opt.Action].Run(&state, opt))

	return nil
}
