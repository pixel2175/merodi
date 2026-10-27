package state

import (
	"merodi/src/config"
	. "merodi/src/config"
	"merodi/src/core/build/compile/jinja"
	"merodi/src/parser"
	"merodi/src/utils"

	. "github.com/pixel2175/errs"

	"github.com/yuin/goldmark"
	gparser "github.com/yuin/goldmark/parser"
	glua "github.com/yuin/gopher-lua"
)

type State struct {
	ProjectDir string
	Config     Data
	Build      Build
	Jinja      map[string]any
	Hooks      map[string]*glua.LFunction
	Lua        Lua
	Markdown   Markdown
	Watch      Watcher
	Action     parser.Action
}

type Lua struct {
	Context *glua.LState
	Aborted bool
}

type Build struct {
	MDPath  string
	DocPath string
	Content string
	Mode    config.Mode
}

type Markdown struct {
	Extensions    []goldmark.Extender
	ParsersOption []gparser.Option
	Page          jinja.Page
}

type Watcher struct {
	Add []string
}

func InitGlobals() State {
	return State{
		ProjectDir: "",
		Config:     Data{},
		Jinja:      make(map[string]any),
		Build:      Build{},
		Watch:      Watcher{},
		Hooks:      make(map[string]*glua.LFunction),
		Lua:        Lua{},
		Action:     parser.Help,
	}
}

func (s *State) RunHook(stage string, args ...string) (err error) {
	defer Catch(&err)
	fn, ok := s.Hooks[stage]
	if !ok {
		return
	}

	s.Lua.Aborted = false

	luaArgs := make([]glua.LValue, len(args))
	for i, arg := range args {
		luaArgs[i] = glua.LString(arg)
	}

	cerr := s.Lua.Context.CallByParam(glua.P{
		Fn:      fn,
		NRet:    0,
		Protect: true,
	}, luaArgs...)

	if s.Lua.Aborted {
		TryE(utils.ErrAborted)
	}
	TryE(cerr)
	return
}
