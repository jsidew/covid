// Package database provides access to covid data sources from the web.
package database

import (
	"fmt"
	"time"
)

// EndpointName is a unique name identifying a web ednpoint and its related resource.
type EndpointName string

/*
DB is the database.
The resources are lazy-loaded from 2 caches, either:
    1. in-memory, for the runtime;
    2. the file-system, as files saved in the specified directory.
If the cache period has expired, or the files don't exist already,
the resources are taken from the web, and then stored in the caches.
*/
type DB struct {
	origin, cachedir string

	expiration time.Duration
	first      EndpointName
	resources  resources
}

/*
New database.
origin is the base URL common to all the resources (e.g. https://raw.githubusercontent.com/CSSEGISandData/COVID-1).
cachedir is the full path of the directory where the resources are cached.
cacheExpiration is the period after which the cache is refreshed (from endpoints under origin).
*/
func New(origin, cachedir string, cacheExpiration time.Duration) *DB {
	return &DB{origin: origin, cachedir: cachedir, expiration: cacheExpiration}
}

// Set a new named endpoint to a web resource.
func (db *DB) Set(n EndpointName, endpoint string) {
	if db.resources == nil {
		db.resources = resources{}
	}
	url := fmt.Sprintf("%s/%s", db.origin, endpoint)
	db.resources.Set(db.cachedir, string(n), url, db.expiration)
	if db.first == "" {
		db.first = n
	}
}

// Latest update time.
func (db *DB) Latest() (time.Time, error) {
	r, err := db.resources.Get(db.first.String())
	if err != nil {
		return time.Time{}, err
	}
	return r.Latest()
}

/*
ActiveCases affected, selected by country and time.
The active cases is a difference between total cases (subtrahend)
and other cases (minuends) that are no more considered active.
*/
func (db *DB) ActiveCases(country string, t time.Time, subtrahend EndpointName, minuends ...EndpointName) (int, error) {
	r, err := db.resources.Get(subtrahend.String())
	if err != nil {
		return 0, err
	}

	var c int
	c, err = r.Cases(country, t)
	if err != nil {
		return 0, err
	}
	for _, min := range minuends {
		r, err := db.resources.Get(min.String())
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

// Countries listed in the resources, sorted by their name.
func (db *DB) Countries() ([]string, error) {
	r, err := db.resources.Get(db.first.String())
	if err != nil {
		return nil, err
	}
	return r.Countries(), nil
}

func (e EndpointName) String() string {
	return string(e)
}
