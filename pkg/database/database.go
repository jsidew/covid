package database

import (
	"fmt"
	"time"
)

type DB struct {
	origin, cachedir string

	expiration time.Duration
	first      string
	resources  resources
}

func New(origin, cachedir string, cacheExpiration time.Duration) *DB {
	return &DB{origin: origin, cachedir: cachedir, expiration: cacheExpiration}
}

func (db *DB) Set(name, endpoint string) {
	if db.resources == nil {
		db.resources = resources{}
	}
	url := fmt.Sprintf("%s/%s", db.origin, endpoint)
	db.resources.Set(db.cachedir, name, url, db.expiration)
	if db.first == "" {
		db.first = name
	}
}

func (db *DB) Latest() (time.Time, error) {
	r, err := db.resources.Get(db.first)
	if err != nil {
		return time.Time{}, err
	}
	return r.Latest()
}

func (db *DB) ActiveCases(country string, t time.Time, name string, subtracted ...string) (int, error) {
	r, err := db.resources.Get(name)
	if err != nil {
		return 0, err
	}

	var c int
	c, err = r.Cases(country, t)
	if err != nil {
		return 0, err
	}
	for _, sub := range subtracted {
		r, err := db.resources.Get(sub)
		if err != nil {
			return 0, err
		}
		s, err := r.Cases(country, t)
		if err != nil {
			return 0, err
		}
		c -= s
	}

	return c, nil
}
