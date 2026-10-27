package compile

import (
	. "github.com/pixel2175/errs"
	"os"
	"path/filepath"
	"strings"

	"merodi/src/config"
)

func (self *Compile) destRoot() string {
	if self.State.Build.Mode == config.Release {
		return self.State.Config.Tree.ReleaseDest
	}
	return self.State.Config.Tree.DraftDest
}

func (self *Compile) docDestPath(mdPath string) (dest string, err error) {
	defer Catch(&err)
	mdRoot := TryV(filepath.Abs(self.State.Config.Tree.Markdown))
	absMd := TryV(filepath.Abs(mdPath))
	rel := TryV(filepath.Rel(mdRoot, absMd))

	name := strings.TrimSuffix(rel, ".md")
	if name != rel {
		name += ".html"
	}
	dest = TryV(filepath.Abs(filepath.Join(self.destRoot(), name)))
	return
}

func (self *Compile) WriteDoc(docContent string) (err error) {
	defer Catch(&err)

	TryE(os.MkdirAll(filepath.Dir(self.Build.DocPath), 0o755))
	TryE(os.WriteFile(self.Build.DocPath, []byte(docContent), 0o644))
	return
}
