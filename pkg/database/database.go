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
