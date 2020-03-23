package database_test

import (
	"fmt"
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

const (
	tempfolderpfx = "temp-"
	fixfolder     = "fixture"
)

var env *setting

func TestDB(t *testing.T) {
	defer setup().Teardown()

	const sleepfor = 50 * time.Millisecond

	for _, dur := range []time.Duration{
		1 * time.Second,
		300 * time.Millisecond,
		60 * time.Millisecond,
		10 * time.Millisecond,
	} {
		desc := fmt.Sprintf("with %s expiration", dur)
		t.Run(desc, func(t *testing.T) {
			db := database.New(env.ServerURL(), env.TmpDir(), dur)
			db.Set("confirmed", "/confirmed.csv")
			db.Set("recovered", "/recovered.csv")
			db.Set("dead", "/deaths.csv")
			t.Run("LatestTime", func(t *testing.T) {
				latest, err := db.Latest()
				require.NoError(t, err, "error")
				assert.Equal(t, date(2020, time.March, 19), latest, "time")
			})
			t.Run("ActiveCases", func(t *testing.T) {
				cases, err := db.ActiveCases("italy", date(2020, time.March, 3), "confirmed", "recovered", "dead")
				require.NoError(t, err, "error")
				assert.Equal(t, 2263, cases, "active cases")
			})
			t.Run("Countries", func(t *testing.T) {
				countries, err := db.Countries()
				require.NoError(t, err, "error")
				require.Len(t, countries, 155, "total count")
				assert.Equal(t, "Afghanistan", countries[0])
				assert.Equal(t, "Zambia", countries[len(countries)-1])
				assert.Equal(t, "Kazakhstan", countries[77])
			})
		})
		time.Sleep(sleepfor)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	path := strings.ReplaceAll(r.URL.EscapedPath(), "/", "")
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
	env.tmpdir, err = ioutil.TempDir(env.fix.dir, tempfolderpfx)
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

func date(year int, m time.Month, day int) time.Time {
	return time.Date(year, m, day, 0, 0, 0, 0, time.UTC)
}

func panicif(err error) {
	if err != nil {
		panic(err)
	}
}
