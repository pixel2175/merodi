package config

import (
	. "github.com/pixel2175/errs"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Mode int

const ConfigFileName = "config.toml"
const (
	None Mode = iota
	Release
	Draft
)

type Project struct {
	Name        string `toml:"name"`
	Version     string `toml:"version"`
	Description string `toml:"description"`
}

type Tree struct {
	Markdown    string `toml:"markdown"`
	Templates   string `toml:"templates"`
	Plugins     string `toml:"plugins"`
	DraftDest   string `toml:"draft_dest"`
	ReleaseDest string `toml:"release_dest"`
}

type Data struct {
	Project Project `toml:"project"`
	Tree    Tree    `toml:"tree"`
}

type Config struct {
	Data       Data
	ProjectDir string
}

func (self *Config) DefaultData() {
	self.Data = Data{
		Project: Project{
			Name:        filepath.Base(self.ProjectDir),
			Version:     "0.1.0",
			Description: "Add your description here",
		},
		Tree: Tree{
			Markdown:    "src/md",
			Templates:   "src/templates",
			Plugins:     "src/plugins",
			DraftDest:   "build/draft",
			ReleaseDest: "build/release",
		},
	}
}

func ModeToString(mode Mode) string {
	switch mode {
	case Release:
		return "release"
	case Draft:
		return "Draft"
	default:
		return "None"
	}
}

func (self *Config) ConfigPath() string {
	return filepath.Join(self.ProjectDir, ConfigFileName)
}

func (self *Config) Init() (err error) {
	defer Catch(&err)
	TryE(os.MkdirAll(self.ProjectDir, 0755))
	f := TryV(os.Create(self.ConfigPath()))
	defer f.Close()
	self.DefaultData()
	TryE(toml.NewEncoder(f).Encode(self.Data))
	return
}

func (self *Config) Read() (err error) {
	defer Catch(&err)
	TryV(toml.DecodeFile(self.ConfigPath(), &self.Data))
	return
}

func (self *Config) Write() (err error) {
	defer Catch(&err)
	f := TryV(os.Create(self.ConfigPath()))
	defer f.Close()
	TryE(toml.NewEncoder(f).Encode(self.Data))
	return
}
