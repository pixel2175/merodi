package attrparser

import (
	"bytes"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var Extension goldmark.Extender = ext{}

type ext struct{}

func (ext) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(wrapParser{}, 150),
			util.Prioritized(lineParser{}, 151),
		),
		parser.WithInlineParsers(util.Prioritized(inline{}, 199)),
		parser.WithASTTransformers(util.Prioritized(finish{}, 100)),
	)
}

var (
	kindWrap = ast.NewNodeKind("AttrWrap")
	kindLine = ast.NewNodeKind("AttrLine")
	stackKey = parser.NewContextKey()
)

type node struct {
	ast.BaseBlock
	kind  ast.NodeKind
	attrs parser.Attributes
}

func (n *node) Kind() ast.NodeKind       { return n.kind }
func (n *node) Dump(s []byte, level int) { ast.DumpHelper(n, s, level, nil, nil) }

func attrsAt(b []byte) (parser.Attributes, int, bool) {
	if len(b) == 0 || b[0] != '{' {
		return nil, 0, false
	}
	r := text.NewReader(b)
	attrs, ok := parser.ParseAttributes(r)
	_, s := r.Position()
	return attrs, s.Start, ok
}

func add(n ast.Node, attrs parser.Attributes) {
	for _, a := range attrs {
		if old, ok := n.Attribute(a.Name); ok && string(a.Name) == "class" {
			o, _ := old.([]byte)
			v, _ := a.Value.([]byte)
			a.Value = append(append(append([]byte{}, o...), ' '), v...)
		}
		n.SetAttribute(a.Name, a.Value)
	}
}

var wrapRe = regexp.MustCompile(`^\[\s*(?:>|#{1,6}\s|[-*+]\s|\d+[.)]\s)`)

type wrapParser struct{}

func (wrapParser) Trigger() []byte { return []byte{'['} }

func (wrapParser) Open(_ ast.Node, r text.Reader, _ parser.Context) (ast.Node, parser.State) {
	line, _ := r.PeekLine()
	t := bytes.TrimLeft(line, " \t")
	tt := bytes.TrimRight(t, " \t\r\n")
	j := bytes.LastIndexByte(tt, '{')
	if j < 1 || !wrapRe.Match(tt) {
		return nil, parser.NoChildren
	}
	i := bytes.LastIndexByte(tt[:j], ']')
	if i < 0 || len(bytes.Trim(tt[i+1:j], " \t")) != 0 {
		return nil, parser.NoChildren
	}
	if _, n, ok := attrsAt(tt[j:]); !ok || n != len(tt)-j {
		return nil, parser.NoChildren
	}
	skip := len(line) - len(t) + 1
	for skip < len(line) && (line[skip] == ' ' || line[skip] == '\t') {
		skip++
	}
	r.Advance(skip)
	return &node{kind: kindWrap}, parser.HasChildren
}

func (wrapParser) Continue(ast.Node, text.Reader, parser.Context) parser.State {
	return parser.Close
}
func (wrapParser) Close(ast.Node, text.Reader, parser.Context) {}
func (wrapParser) CanInterruptParagraph() bool                 { return false }
func (wrapParser) CanAcceptIndentedLine() bool                 { return false }


type lineParser struct{}

func (lineParser) Trigger() []byte { return []byte{'{'} }

func (lineParser) Open(_ ast.Node, r text.Reader, _ parser.Context) (ast.Node, parser.State) {
	line, _ := r.PeekLine()
	t := bytes.TrimSpace(line)
	attrs, n, ok := attrsAt(t)
	if !ok || n != len(t) {
		return nil, parser.NoChildren
	}
	r.Advance(len(line))
	return &node{kind: kindLine, attrs: attrs}, parser.NoChildren
}

func (lineParser) Continue(ast.Node, text.Reader, parser.Context) parser.State {
	return parser.Close
}
func (lineParser) Close(ast.Node, text.Reader, parser.Context) {}
func (lineParser) CanInterruptParagraph() bool                 { return true }
func (lineParser) CanAcceptIndentedLine() bool                 { return false }


type opener struct {
	end, n int
	attrs  parser.Attributes
}

func closeIdx(b []byte) int {
	depth := 0
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case '\\':
			i++
		case '`':
			n := 1
			for i+n < len(b) && b[i+n] == '`' {
				n++
			}
			if j := bytes.Index(b[i+n:], bytes.Repeat([]byte("`"), n)); j >= 0 {
				i += n + j + n - 1
			} else {
				i += n - 1
			}
		case '[':
			depth++
		case ']':
			if depth--; depth == 0 {
				return i
			}
		}
	}
	return -1
}

type inline struct{}

func (inline) Trigger() []byte { return []byte{'[', ']', '{'} }

func (inline) Parse(parent ast.Node, r text.Reader, pc parser.Context) ast.Node {
	line, seg := r.PeekLine()
	pos := seg.Start
	st, _ := pc.Get(stackKey).([]opener)
	eat := func(n int) ast.Node {
		r.Advance(n)
		return ast.NewTextSegment(text.NewSegment(pos, pos))
	}
	if line[0] == '{' {
		attrs, n, ok := attrsAt(line)
		lines := parent.Lines()
		if !ok || lines.Len() == 0 || seg.Stop != lines.At(lines.Len()-1).Stop ||
			len(bytes.TrimSpace(line[n:])) != 0 {
			return nil
		}
		t := parent
		if p := parent.Parent(); p != nil && (p.Kind() == ast.KindBlockquote || p.Kind() == ast.KindListItem) {
			t = p
		}
		add(t, attrs)
		if x, ok := parent.LastChild().(*ast.Text); ok {
			x.Segment = x.Segment.TrimRightSpace(r.Source())
		}
		return eat(n)
	}

	if line[0] == '[' {
		end := closeIdx(line)
		if end < 0 {
			return nil
		}
		j := end + 1
		for j < len(line) && (line[j] == ' ' || line[j] == '\t') {
			j++
		}
		attrs, n, ok := attrsAt(line[j:])
		if !ok {
			return nil
		}
		pc.Set(stackKey, append(st, opener{pos + end, j - end - 1 + n, attrs}))
		return eat(1)
	}

	for len(st) > 0 && st[len(st)-1].end < pos {
		st = st[:len(st)-1]
	}
	if len(st) > 0 && st[len(st)-1].end == pos {
		o := st[len(st)-1]
		pc.Set(stackKey, st[:len(st)-1])
		add(parent, o.attrs)
		return eat(1 + o.n)
	}
	pc.Set(stackKey, st)

	j := 1
	for j < len(line) && (line[j] == ' ' || line[j] == '\t') {
		j++
	}
	if attrs, n, ok := attrsAt(line[j:]); ok {	for c := parent; c.Parent() != nil; c = c.Parent() {
			if w, ok := c.Parent().(*node); ok && w.kind == kindWrap {
				add(c, attrs)
				if t, ok := parent.LastChild().(*ast.Text); ok {
					t.Segment = t.Segment.TrimRightSpace(r.Source())
				}
				return eat(j + n)
			}
		}
	}
	return nil
}

type finish struct{}

func (finish) Transform(doc *ast.Document, _ text.Reader, _ parser.Context) {
	var ns []*node
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if v, ok := n.(*node); ok && entering {
			ns = append(ns, v)
		}
		return ast.WalkContinue, nil
	})
	for _, n := range ns {
		p := n.Parent()
		switch n.kind {
		case kindLine:
			if prev := n.PreviousSibling(); prev != nil {
				add(prev, n.attrs)
			}
		case kindWrap:
			for c := n.FirstChild(); c != nil; c = n.FirstChild() {
				n.RemoveChild(n, c)
				p.InsertBefore(p, n, c)
			}
		}
		p.RemoveChild(p, n)
	}
}
