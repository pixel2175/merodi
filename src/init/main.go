package init

import (
	"errors"
	"merodi/src/config"
	"merodi/src/project"
	. "github.com/pixel2175/errs"
	"merodi/src/utils/log"
)

type Init struct {
	ProjectDir string
	cfg        config.Tree
}

var ErrProjectExists = errors.New("project already initialized: `config.toml` already exists")

func (init *Init) GetProjectDir() (err error) {
	defer Catch(&err)
	init.ProjectDir = TryV(project.GetProjectDir(init.ProjectDir))
	if project.IsProjectExists(init.ProjectDir) {
		TryE(ErrProjectExists)
	}
	return
}

func (init *Init) InitConfig() error {
	cfg := config.Config{}
	cfg.ProjectDir = init.ProjectDir
	cfg.Init()
	init.cfg = cfg.Data.Tree
	return cfg.Write()
}

func (init Init) Run() (err error) {
	defer Catch(&err)
	TryE(init.GetProjectDir())
	TryE(init.InitConfig())
	TryE(init.Bootstrap())
	log.Info("Project initialized.")
	return nil
}
