package compile

import (
	. "github.com/pixel2175/errs"
	"path/filepath"
)

func (self *Compile) hook(name string) {
	TryE(self.State.RunHook(name))
}

func (self *Compile) ResolvePaths(mdFile string) {
	mdPath := filepath.Join(self.State.ProjectDir, mdFile)
	docPath := TryV(self.docDestPath(mdPath))
	self.Build.MDPath, self.Build.DocPath = mdPath, docPath
}

func (self *Compile) readMarkdown() {
	self.hook("before_read")
	self.Build.Content = TryV(self.MD.ReadMD(self.Build.MDPath))
	self.hook("read_md")
}

func (self *Compile) convertToDoc() {
	self.MD.Page = *self.Jinja.Page
	self.Build.Content = TryV(self.MD.MdToDoc(self.Build.Content))
	self.hook("doc_content")
}

func (self *Compile) applyJinja() {
	self.Jinja.Data = self.State.Jinja
	TryE(self.Jinja.JinjaHandler(&self.Build.Content))
	self.hook("apply_doc_jinja")
}

func (self *Compile) reset() {
	self.Build.Content = ""
	self.Build.DocPath = ""
	self.Build.MDPath = ""
}
