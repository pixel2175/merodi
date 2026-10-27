package main

import (
	"merodi/src/core"
	. "merodi/src/init"
	"merodi/src/parser"
	. "github.com/pixel2175/errs"
	"merodi/src/utils/log"
	"os"
)

func main() {
	var err error
	defer Catch(&err, func(err *error) {
		log.Die((*err).Error())
	})

	opt := parser.Parse(os.Args[1:])

	switch opt.Action {
	case parser.Help:
		parser.PrintHelp()
	case parser.Init:
		TryE((&Init{ProjectDir: opt.ProjectDir}).Run())
	case parser.Version:
		log.Info(log.Title("Version"), "Merodi %s", parser.VERSION)
	default:
		TryE(core.Run(&opt))
	}
}
