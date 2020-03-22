package database

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

const filext = ".csv"

type resources map[string]resource

type resource struct {
	name, url, filepath string

	expire time.Duration
	mx     matrix
}

func (r resources) Set(db, name, url string, expire time.Duration) {
	r[name] = resource{
		name: name, url: url,
		filepath: filepath.Join(db, name+filext),
		expire:   expire,
	}
}

func (r resources) Get(name string) (matrix, error) {
	res, ok := r[name]
	if !ok {
		return nil, fmt.Errorf("unkown resource name `%s`", name)
	}
	if res.mx != nil {
		return res.mx, nil
	}
	m, err := res.open()
	if err != nil {
		return nil, err
	}
	r[name] = res.cloneWithM(m)
	return r[name].mx, nil
}

func (r resource) cloneWithM(m matrix) resource {
	return resource{mx: m, name: r.name, url: r.url, filepath: r.filepath}
}

func (r resource) open() (matrix, error) {
	f, err := os.Open(r.filepath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		goto update
	}
	defer f.Close()

	if info, err := f.Stat(); err != nil {
		return nil, err
	} else if info.ModTime().Before(time.Now().Add(-r.expire)) {
		f.Close()
		goto update
	}
	goto parse

update:
	f, err = os.Create(r.filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = r.update(f)
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

parse:
	results, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	m := matrix(results)
	err = m.validate()
	if err != nil {
		return nil, err
	}
	m = m.cleancopy()

	return m, err
}

func (r resource) update(w io.Writer) error {

	var body io.Reader
	{
		resp, err := http.Get(r.url)
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