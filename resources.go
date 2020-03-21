package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type resources map[string]resource

type resource struct {
	Name, URL, FilePath string

	M matrix
}

func (r resources) Set(db, name, url string) {
	r[name] = resource{Name: name, URL: url, FilePath: filepath.Join(db, name+".csv")}
}

func (r resources) Get(name string) resource {
	return r[name]
}

func (r resources) Load(names ...string) error {
	for _, name := range names {
		res, ok := r[name]
		if !ok {
			return fmt.Errorf("unkown resource name `%s`", name)
		}
		m, err := res.open()
		if err != nil {
			return err
		}
		r[name] = res.cloneWithM(m)
	}
	return nil
}

func (r resource) cloneWithM(m matrix) resource {
	return resource{M: m, Name: r.Name, URL: r.URL, FilePath: r.FilePath}
}

func (r resource) open() (matrix, error) {
	f, err := os.Open(r.FilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		goto update
	}
	defer f.Close()

	if info, err := f.Stat(); err != nil {
		return nil, err
	} else if info.ModTime().Before(time.Now().Add(-refreshInterval)) {
		f.Close()
		goto update
	}
	goto parse

update:
	f, err = os.Create(r.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = r.update(f)
	if err != nil {
		return nil, err
	}

parse:
	results, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	m := matrix(results)
	err = m.Validate()
	if err != nil {
		return nil, err
	}
	m = m.Clean()

	return m, err
}

func (r resource) update(w io.Writer) error {

	var body io.Reader
	{
		resp, err := http.Get(r.URL)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
		}
		defer resp.Body.Close()
		body = resp.Body
	}

	_, err := bufio.NewReader(body).WriteTo(w)

	return err
}
