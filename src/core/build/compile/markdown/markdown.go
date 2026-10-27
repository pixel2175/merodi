package markdown

import (
	"bytes"
	"merodi/src/core/build/compile/jinja"
	"merodi/src/core/state"
	. "github.com/pixel2175/errs"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type Markdown struct {
	*state.Markdown
	state.State
}

func (self *Markdown) Init(state *state.State) {
	self.Markdown = &state.Markdown
	self.Extensions = state.Markdown.Extensions
	self.ParsersOption = state.Markdown.ParsersOption
}

func (self *Markdown) ReadMD(mdfile string) (content string, err error) {
	defer Catch(&err)
	content = string(TryV(os.ReadFile(mdfile)))
	return
}

func (self *Markdown) MdToDoc(mdcontent string) (doc string, err error) {
	defer Catch(&err)
	var raw_doc bytes.Buffer
	mdBytes := []byte(mdcontent)

	gm := goldmark.New(
		goldmark.WithParserOptions(self.ParsersOption...),
		goldmark.WithExtensions(self.Extensions...),
		goldmark.WithRendererOptions(html.WithUnsafe()),
		goldmark.WithParserOptions(
			parser.WithBlockParsers(util.Prioritized(&jinja.JinjaEscaper{}, 50)),
		),
	)

	node := gm.Parser().Parse(text.NewReader(mdBytes))
	TryE(gm.Renderer().Render(&raw_doc, mdBytes, node))
	body := raw_doc.String()

	if !self.Page.Set {
		return body, nil
	}

	return self.wrapHTML(body), nil
}

func (self *Markdown) wrapHTML(body string) string {
	var b bytes.Buffer

	b.WriteString("<!DOCTYPE html>\n")
	b.WriteString(`<html lang="` + self.Page.Lang + `">` + "\n<head>\n")
	b.WriteString("<title>" + self.Page.Title + "</title>\n")

	for _, m := range self.Page.Metas {
		b.WriteString(`<meta name="` + m.Name + `" content="` + m.Content + `">` + "\n")
	}
	for _, s := range self.Page.Styles {
		b.WriteString(`<link rel="stylesheet" href="` + s + `">` + "\n")
	}
	for _, s := range self.Page.Scripts {
		b.WriteString(`<script src="` + s + `" defer></script>` + "\n")
	}

	b.WriteString("</head>\n<body>\n")
	b.WriteString(body)
	b.WriteString("</body>\n</html>")
	return b.String()
}
