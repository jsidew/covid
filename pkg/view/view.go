// Package view provides ways to textually describe computed outputs through template views.
package view

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/jsidew/covid/internal/errors"
)

func init() {
	errors.Prefix = "view"
}

// TemplateName is a name of a template
type TemplateName string

// View of a template.
type View struct {
	tpl *template.Template
}

/*
New template view created from a selected template in the template directory dir.
If the template directory has no template files (*.tpl),
a default template will be created with content from Template.
If selected is empty, the default template (Template) will be used.
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
The following custom template functions can be used:
	- print LANG ARGUMENTS, like fmt.Print, but formatted according to local LANG (see doc for golang.org/x/text/message);
	- printf LANG FORMAT ARGUMENTS, like fmt.Printf, but formatted according to local LANG (see doc for golang.org/x/text/message);
	- fmtdate LAYOUT TIME, like t.Format(LAYOUT), where t is the TIME object.
*/
func (v *View) Execute(w io.Writer, data interface{}) error {
	err := v.tpl.Execute(w, data)
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
