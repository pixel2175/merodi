package state

import (
	"merodi/src/config"
	"merodi/src/project"
	. "github.com/pixel2175/errs"
	"os"
)

func (self *State) GoToProjectDir(dir string) (err error) {
	defer Catch(&err)
	dir = TryV(project.SetProjectDir(dir))
	self.ProjectDir = dir
	TryE(os.Chdir(dir))
	return
}

func (self *State) LoadConfig() (err error) {
	defer Catch(&err)
	var cfg config.Config
	cfg.ProjectDir = self.ProjectDir
	TryE(cfg.Read())
	self.Config = cfg.Data
	return
}
