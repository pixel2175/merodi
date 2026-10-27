package compile

import (
	. "merodi/src/utils/errs"
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
	defer Handle(&err)
	mdRoot := CheckV(filepath.Abs(self.State.Config.Tree.Markdown))
	absMd := CheckV(filepath.Abs(mdPath))
	rel := CheckV(filepath.Rel(mdRoot, absMd))

	name := strings.TrimSuffix(rel, ".md")
	if name != rel {
		name += ".html"
	}
	dest = CheckV(filepath.Abs(filepath.Join(self.destRoot(), name)))
	return
}

func (self *Compile) WriteDoc(docContent string) (err error) {
	defer Handle(&err)

	CheckE(os.MkdirAll(filepath.Dir(self.Build.DocPath), 0o755))
	CheckE(os.WriteFile(self.Build.DocPath, []byte(docContent), 0o644))
	return
}
