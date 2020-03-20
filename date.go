package main

import (
	"time"
)

type date time.Time

func (d date) Time() time.Time {
	return time.Time(d)
}

func (d date) DaysSince(s date) float64 {
	return -d.Time().Sub(s.Time()).Minutes() / 60 / 24
}

func (d date) AddDays(days int) date {
	return date(d.Time().Add(time.Duration(days) * 24 * time.Hour))
}

func (d date) String() string {
	return d.Time().Format(dateLayout)
}

func (d *date) Set(s string) error {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	*d = date(t)
	return nil
}
