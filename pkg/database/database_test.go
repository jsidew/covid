package database_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jsidew/covid/pkg/database"
)

const fixfolder = "fixture"

var env *setting

func TestDB(t *testing.T) {
	defer setup().Teardown()
	db := database.New(env.ServerURL(), env.TmpDir(), 1*time.Second)
	db.Set("confirmed", "/confirmed.csv")
	db.Set("recovered", "/recovered.csv")
	db.Set("dead", "/deaths.csv")
	t.Run("LatestTime", func(t *testing.T) {
		d, err := db.Latest()
		require.NoError(t, err, "error")
		assert.Equal(t, time.Date(2020, time.March, 19, 0, 0, 0, 0, time.UTC), d, "time")
	})
	t.Run("ActiveCases", func(t *testing.T) {
		// (country string, d date, name string, subtractedNames ...string) (cases int, err error)
	})
	t.Run("Countries", func(t *testing.T) {})
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := strings.ReplaceAll(r.URL.EscapedPath(), "/", "")
	println(path)
	b, err := env.Fixture(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

type fixture struct {
	dir   string
	cache map[string][]byte
}

type setting struct {
	server *httptest.Server
	tmpdir string
	fix    *fixture
}

func setup() *setting {
	var err error
	env = &setting{
		fix: &fixture{},
	}
	env.fix.dir, err = filepath.Abs(fixfolder)
	panicif(err)
	env.tmpdir, err = ioutil.TempDir(env.fix.dir, "temp")
	panicif(err)
	env.server = httptest.NewServer(http.HandlerFunc(handler))
	return env
}

func (s *setting) ServerURL() string {
	return s.server.URL
}

func (s *setting) TmpDir() string {
	return s.tmpdir
}

func (s *setting) Fixture(name string) ([]byte, error) {
	if s.fix.cache == nil {
		s.fix.cache = map[string][]byte{}
	} else {
		c, ok := s.fix.cache[name]
		if ok && c != nil {
			return c, nil
		}
	}
	b, err := ioutil.ReadFile(filepath.Join(s.fix.dir, name))
	if err != nil {
		return nil, err
	}
	s.fix.cache[name] = b
	return b, nil
}

func (s *setting) Teardown() {
	s.server.Close()
	env = nil
	err := os.RemoveAll(s.tmpdir)
	panicif(err)
}

func panicif(err error) {
	if err != nil {
		panic(err)
	}
}
