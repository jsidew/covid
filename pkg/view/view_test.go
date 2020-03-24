package view_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsidew/covid/pkg/view"
)

const tempfolderpfx = "temp-"

var env *setting

func Test(t *testing.T) {
	defer setup(t).Teardown()
	defer env.TmpList()

	t.Run("default template", func(t *testing.T) {
		_, err := view.New(env.TmpDir(), "")
		assert.NoError(t, err, "error")
		assert.True(t, env.TmpExists("default.tpl"), "default template wasn't created")
	})

	t.Run("template not found", func(t *testing.T) {
		_, err := view.New(env.TmpDir(), "phantom")
		assert.EqualError(t, err, `view: template "phantom.tpl" doesn't exist; defined templates are: "default.tpl"`, "error")
	})

	t.Run("custom template", func(t *testing.T) {
		env.TmpCreate("custom.tpl", nil)
		_, err := view.New(env.TmpDir(), "custom")
		assert.NoError(t, err, "error")
	})

	t.Run("execute", func(t *testing.T) {
		b := strings.Builder{}
		defer b.Reset()
		env.TmpCreate("custom.tpl", []byte(`Cases: {{ print "en" .Current.Cases }}
Rate: {{ printf "en" "%.4f" .Current.Rate  }}
Date: {{ fmtdate "2 Jan 2006" .Updated }}`))
		env.TmpCreate("test.tpl", nil)
		v, err := view.New(env.TmpDir(), "custom")
		require.NoError(t, err, "New error")
		v.Current.Cases = 1435678
		v.Current.Rate = 12343.45945857
		v.Updated = time.Date(2020, 03, 23, 11, 41, 55, 0, time.UTC)
		err = v.Execute(&b)
		require.NoError(t, err, "View.Execute error")
		rows := strings.Split(b.String(), "\n")
		require.Len(t, rows, 3, "view rows")
		assert.Equal(t, `Cases: 1,435,678`, rows[0], "final view: print")
		assert.Equal(t, `Rate: 12,343.4595`, rows[1], "final view: printf")
		assert.Equal(t, `Date: 23 Mar 2020`, rows[2], "final view: fmtdate")
	})

}

type setting struct {
	t      *testing.T
	tmpdir string
}

func setup(t *testing.T) *setting {
	var err error
	env = &setting{t: t}
	dir, err := os.Getwd()
	panicif(err)
	env.tmpdir, err = ioutil.TempDir(dir, tempfolderpfx)
	panicif(err)
	return env
}

func (s *setting) TmpDir() string {
	return s.tmpdir
}

func (s *setting) TmpList() {
	err := filepath.Walk(s.TmpDir(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		s.t.Log(path)
		return nil
	})
	panicif(err)
}

func (s *setting) TmpExists(name string) bool {
	info, err := os.Stat(filepath.Join(s.TmpDir(), name))
	return !os.IsNotExist(err) && info.Mode().IsRegular()
}

func (s *setting) TmpCreate(name string, data []byte) {
	err := ioutil.WriteFile(filepath.Join(s.TmpDir(), name), data, 0644)
	panicif(err)
}

func (s *setting) Teardown() {
	env = nil
	err := os.RemoveAll(s.tmpdir)
	panicif(err)
}

func panicif(err error) {
	if err != nil {
		panic(err)
	}
}
