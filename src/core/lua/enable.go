package lua

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	glua "github.com/yuin/gopher-lua"
	"merodi/src/core/build/compile/markdown/attrparser"
)

var extensionsMap = map[string]goldmark.Extender{
	"Table":          extension.Table,
	"Strikethrough":  extension.Strikethrough,
	"Linkify":        extension.Linkify,
	"TaskList":       extension.TaskList,
	"GFM":            extension.GFM,
	"DefinitionList": extension.DefinitionList,
	"Footnote":       extension.Footnote,
	"Typographer":    extension.Typographer,
	"CJK":            extension.CJK,
	"HTMLAttr":       attrparser.Extension,
}

var parsersMap = map[string]parser.Option{
	"AutoHeadingID": parser.WithAutoHeadingID(),
	"EscapedSpace":  parser.WithEscapedSpace(),
}

var blockParsersMap = map[string]parser.BlockParser{
	"SetextHeading":   parser.NewSetextHeadingParser(),
	"ThematicBreak":   parser.NewThematicBreakParser(),
	"List":            parser.NewListParser(),
	"ListItem":        parser.NewListItemParser(),
	"CodeBlock":       parser.NewCodeBlockParser(),
	"ATXHeading":      parser.NewATXHeadingParser(),
	"FencedCodeBlock": parser.NewFencedCodeBlockParser(),
	"Blockquote":      parser.NewBlockquoteParser(),
	"HTMLBlock":       parser.NewHTMLBlockParser(),
	"Paragraph":       parser.NewParagraphParser(),
}

var inlineParsersMap = map[string]parser.InlineParser{
	"CodeSpan": parser.NewCodeSpanParser(),
	"Link":     parser.NewLinkParser(),
	"AutoLink": parser.NewAutoLinkParser(),
	"RawHTML":  parser.NewRawHTMLParser(),
	"Emphasis": parser.NewEmphasisParser(),
}

func (self *Lua) registerEnable() {
	enableTable := self.Context.NewTable()
	meta := self.Context.NewTable()

	self.Context.SetField(meta, "__index", self.Context.NewFunction(
		func(L *glua.LState) int {
			name := L.CheckString(2)

			fn := self.Context.NewFunction(func(L *glua.LState) int {
				if ext, ok := extensionsMap[name]; ok {
					self.State.Markdown.Extensions = append(self.State.Markdown.Extensions, ext)
					return 0
				}

				if option, ok := parsersMap[name]; ok {
					self.State.Markdown.ParsersOption = append(self.State.Markdown.ParsersOption, option)
					return 0
				}

				if p, ok := inlineParsersMap[name]; ok {
					self.State.Markdown.ParsersOption = append(
						self.State.Markdown.ParsersOption,
						parser.WithInlineParsers(
							util.Prioritized(p, 100),
						),
					)
					return 0
				}

				if p, ok := blockParsersMap[name]; ok {
					self.State.Markdown.ParsersOption = append(
						self.State.Markdown.ParsersOption,
						parser.WithBlockParsers(
							util.Prioritized(p, 100),
						),
					)
					return 0
				}

				L.RaiseError("unknown parser or extension: %s", name)
				return 0
			})

			L.Push(fn)
			return 1
		},
	))

	enableTable.Metatable = meta
	self.Merodi.RawSetString("enable", enableTable)
}
