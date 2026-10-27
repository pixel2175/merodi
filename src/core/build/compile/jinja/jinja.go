package jinja

import (
	"errors"
	"fmt"
	. "github.com/pixel2175/errs"
	"strings"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/config"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
	"github.com/nikolalohinski/gonja/v2/parser"
)

type JinjaEngine struct {
	Title  string
	Loader loaders.Loader
	Config *config.Config
	Data   map[string]any
	Page   *Page
}

func (self *JinjaEngine) Init(templates string) (err error) {
	defer Catch(&err)
	self.Loader = TryV(loaders.NewFileSystemLoader(templates))
	self.Config = gonja.DefaultConfig
	return
}

func (self *JinjaEngine) JinjaHandler(html_content *string) (err error) {
	defer Catch(&err)
	self.Page = &Page{Title: self.Title, Lang: "en"}

	TryE(self.registerDocument())

	shiftedLoader := TryV(loaders.NewShiftedLoader(
		"root",
		strings.NewReader(*html_content),
		self.Loader,
	))

	tpl, terr := exec.NewTemplate(
		"root",
		self.Config,
		shiftedLoader,
		gonja.DefaultEnvironment,
	)
	if terr != nil {
		var syntaxErr *parser.SyntaxError
		if errors.As(terr, &syntaxErr) {
			terr = fmt.Errorf(
				"Jinja syntax error:\n- Line=%d, Column=%d\n- %s",
				syntaxErr.Line,
				syntaxErr.Column,
				syntaxErr.Message,
			)
		}
		TryE(terr)
	}

	ctx := exec.NewContext(self.Data)
	ctx.Set("document", self.Page)

	*html_content = TryV(tpl.ExecuteToString(ctx))
	return
}
