// Package view provides ways to textually describe computed outputs through template views.
package view

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/jsidew/covid/internal/errors"
)

const (
	// Resolving status
	Resolving uint8 = iota + 1

	// ResolvingSlowly status
	ResolvingSlowly

	// UnderControl status
	UnderControl

	// BarelyUnderControl status
	BarelyUnderControl

	// HardToControl status
	HardToControl

	// LoosingControl status
	LoosingControl

	// OutOfControl status
	OutOfControl
)

func init() {
	errors.Prefix = "view"
}

// TemplateName is a name of a template
type TemplateName string

// View of a template.
// Use View's fields to pass information to your template.
type View struct {
	Country string
	Updated time.Time

	Status struct {
		Score     uint8
		Resolving bool
		Improving bool
	}

	Current struct {
		Rate  float64
		Cases int
	}

	Comparison struct {
		Rate        float64
		RateOfRates float64
	}

	Recovery struct {
		DaysTo1     float64
		DaysToStart float64
		DaysToPeak  float64
		PeakCases   float64
	}

	Forecast struct {
		Growth string
		Cases  float64
		Days   int
	}

	tpl *template.Template
}

/*
New template view created from a selected template in the template directory dir.
If the template directory has no template files (*.tpl),
a default template will be created with content from Template.
If selected is empty, the default template (Template) will be used.

The following custom template functions can be used:
	- print LANG ARGUMENTS, like fmt.Print, but formatted according to local LANG (see doc for golang.org/x/text/message);
	- printf LANG FORMAT ARGUMENTS, like fmt.Printf, but formatted according to local LANG (see doc for golang.org/x/text/message);
	- fmtdate LAYOUT TIME, like t.Format(LAYOUT), where t is the TIME object.
*/
func New(dir string, selected TemplateName) (*View, error) {
	// create default template file if doesn't exist.
	if p := filepath.Join(dir, Name.String()+Extension); !exists(p) {
		err := ioutil.WriteFile(p, []byte(Template), 0644)
		if err != nil {
			return nil, errors.W(err)
		}
	}

	root, err := template.New("root").Funcs(funcMap).
		ParseGlob(filepath.Join(dir, "*"+Extension))
	if err != nil {
		return nil, errors.W(err)
	}

	name := strings.TrimSpace(selected.String())
	if name == "" {
		name = Name.String()
	}
	if !strings.HasSuffix(name, Extension) {
		name += Extension
	}
	t := root.Lookup(name)
	if t == nil {
		return nil, errors.F(`template "%s" doesn't exist%s`, name, root.DefinedTemplates())
	}

	return &View{tpl: t}, nil
}

/*
Execute a selected template, wrapping around text/template.Execute.
*/
func (v *View) Execute(w io.Writer) error {
	err := v.tpl.Execute(w, v)
	if err != nil {
		return errors.W(err)
	}
	return nil
}

func (n TemplateName) String() string {
	return string(n)
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
