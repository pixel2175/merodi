package jinja

import (
	"fmt"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/nodes"
	"github.com/nikolalohinski/gonja/v2/parser"
	"github.com/nikolalohinski/gonja/v2/tokens"
	. "github.com/pixel2175/errs"
)

type Page struct {
	Title, Lang     string
	Styles, Scripts []string
	Metas           []Meta
	Set             bool
}

type Meta struct {
	Name    string
	Content string
}

type Tag struct {
	loc *tokens.Token
	run func(r *exec.Renderer)
}

func (t *Tag) Position() *tokens.Token { return t.loc }
func (t *Tag) String() string          { return "tag" }
func (t *Tag) Execute(r *exec.Renderer, _ *nodes.ControlStructureBlock) error {
	t.run(r)
	return nil
}

func (self *JinjaEngine) fields() map[string]*string {
	return map[string]*string{
		"title": &self.Page.Title, "lang": &self.Page.Lang,
	}
}

func parseArgs(args *parser.Parser, allowed map[string]bool) (exprs map[string]nodes.Expression, err error) {
	defer Catch(&err)
	exprs = map[string]nodes.Expression{}
	for !args.End() {
		key := args.Match(tokens.Name)
		if key == nil {
			TryE(args.Error("expected argument name", nil))
		}
		if !allowed[key.Val] {
			TryE(args.Error(fmt.Sprintf("argument %q not exists", key.Val), key))
		}
		if args.Match(tokens.Assign) == nil {
			TryE(args.Error("expected '='", key))
		}
		exprs[key.Val] = TryV(args.ParseExpression())
	}
	return
}

func (self *JinjaEngine) parseDocument(p *parser.Parser, args *parser.Parser) (cs nodes.ControlStructure, err error) {
	defer Catch(&err)
	dst := self.fields()
	allowed := map[string]bool{}
	for k := range dst {
		allowed[k] = true
	}
	exprs := TryV(parseArgs(args, allowed))
	return &Tag{p.Current(), func(r *exec.Renderer) {
		for k, e := range exprs {
			*dst[k] = r.Eval(e).String()
		}
		self.Page.Set = true
	}}, nil
}

func (self *JinjaEngine) parseMeta(p *parser.Parser, args *parser.Parser) (cs nodes.ControlStructure, err error) {
	defer Catch(&err)
	exprs := TryV(parseArgs(args, map[string]bool{"name": true, "content": true}))
	return &Tag{p.Current(), func(r *exec.Renderer) {
		m := Meta{}
		if e, ok := exprs["name"]; ok {
			m.Name = r.Eval(e).String()
		}
		if e, ok := exprs["content"]; ok {
			m.Content = r.Eval(e).String()
		}
		self.Page.Metas = append(self.Page.Metas, m)
	}}, nil
}

func parseList(list func() *[]string) parser.ControlStructureParser {
	return func(p *parser.Parser, args *parser.Parser) (cs nodes.ControlStructure, err error) {
		defer Catch(&err)
		e := TryV(args.ParseExpression())
		return &Tag{p.Current(), func(r *exec.Renderer) {
			l := list()
			*l = append(*l, r.Eval(e).String())
		}}, nil
	}
}

func (self *JinjaEngine) registerDocument() (err error) {
	defer Catch(&err)
	cs := gonja.DefaultEnvironment.ControlStructures
	tags := map[string]parser.ControlStructureParser{
		"document": self.parseDocument,
		"meta":     self.parseMeta,
		"style":    parseList(func() *[]string { return &self.Page.Styles }),
		"script":   parseList(func() *[]string { return &self.Page.Scripts }),
	}
	for name, p := range tags {
		if cs.Exists(name) {
			TryE(cs.Replace(name, p))
			continue
		}
		TryE(cs.Register(name, p))
	}
	return
}
