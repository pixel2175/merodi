package parser

import (
	"fmt"
	"merodi/src/config"
)

type Action string

const (
	Help    Action = "help"
	Init    Action = "init"
	Version Action = "version"
	Build   Action = "build"
	Watch   Action = "watch"
	Clean   Action = "clean"
	Hook    Action = "hook"
)

var names = map[string]Action{
	"init":    Init,
	"new":     Init,
	"version": Version,
	"build":   Build,
	"watch":   Watch,
	"clean":   Clean,
	"hook":    Hook,
}

type Option struct {
	Action     Action
	Mode       config.Mode
	ProjectDir string
	Hook       string
	Args       []string
}

func Parse(raw []string) Option {
	opt := Option{Action: Help, Mode: config.Draft}

	args := make([]string, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		switch arg := raw[i]; arg {
		case "--release":
			opt.Mode = config.Release
		case "-h", "--help":
			return opt
		case "-v", "--version":
			opt.Action = Version
			return opt
		case "-C":
			if i+1 < len(raw) {
				i++
				opt.ProjectDir = raw[i]
			}
		default:
			args = append(args, arg)
		}
	}

	if len(args) == 0 {
		return opt
	}

	action, ok := names[args[0]]
	if !ok {
		return opt
	}
	opt.Action, args = action, args[1:]

	if opt.Action == Hook {
		if len(args) == 0 {
			opt.Action = Help
			return opt
		}
		opt.Hook, args = args[0], args[1:]
	}
	opt.Args = args

	return opt
}

var VERSION string

func PrintHelp() {
	fmt.Printf(`Merodi %s - static site generator

USAGE:
  merodi <command> [args] [options]

COMMANDS:
  init, new           Create a new project
  build               Build the site
  watch               Rebuild on file changes
  clean               Remove build output
  hook <name> [args]  Run a hook
  version             Show version

OPTIONS:
  -C <dir>            Project directory (default: current)
  --release           Build in release mode (default: draft)
  -h, --help          Show this help
  -v, --version       Show version

EXAMPLES:
  merodi new -C blog
  merodi build -C blog --release
  merodi hook deploy arg1 arg2 -C blog
  merodi watch  

`, VERSION)
}
