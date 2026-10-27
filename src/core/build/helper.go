package build

import (
	"io/fs"
	. "github.com/pixel2175/errs"
	"path/filepath"
)

func (build *Build) isSource(entry fs.DirEntry) bool {
	return entry.Type().IsRegular() || entry.Type()&fs.ModeSymlink != 0
}

func (build *Build) walkAndBuild() (err error) {
	defer Catch(&err)
	TryE(filepath.WalkDir(build.state.Config.Tree.Markdown, build.Visit))
	return
}
