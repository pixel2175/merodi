package init

import (
	"merodi/src/config"
	. "github.com/pixel2175/errs"
	"os"
	"path/filepath"
)

func (init *Init) CreateDirectories(tree config.Tree) (err error) {
	defer Catch(&err)
	for _, dir := range []string{
		tree.Markdown,
		tree.Templates,
		tree.Plugins,
	} {
		path := filepath.Join(init.ProjectDir, dir)
		TryE(os.MkdirAll(path, 0755))
	}
	return
}

func (init *Init) CreateFiles(tree config.Tree) (err error) {
	defer Catch(&err)
	files := map[string]string{
		filepath.Join(init.ProjectDir, tree.Markdown, "index.md"):     MarkdownContent,
		filepath.Join(init.ProjectDir, tree.Templates, "layout.html"): HTMLContent,
		filepath.Join(init.ProjectDir, tree.Plugins, "main.lua"):      PluginsContent,
	}

	for path, content := range files {
		TryE(os.WriteFile(path, []byte(content), 0644))
	}
	return
}

func (self *Init) Bootstrap() (err error) {
	defer Catch(&err)
	TryE(self.CreateDirectories(self.cfg))
	TryE(self.CreateFiles(self.cfg))
	return
}
