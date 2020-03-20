package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const formatFromCSV = "1/2/06"

type matrix [][]string

func (m matrix) Cases(country string, d date) (int, error) {
	for _, row := range m[1:] {
		if !strings.EqualFold(country, row[1]) {
			continue
		}
		for j, cell := range row[4:] {
			t, err := time.Parse(formatFromCSV, m[0][j+4])
			if err != nil {
				return 0, err
			}
			if !t.Equal(d.Time()) {
				continue
			}
			n, err := strconv.Atoi(cell)
			if err != nil {
				return 0, err
			}
			return n, nil
		}
	}

	return 0, nil
}

func (m matrix) LastDate() (date, error) {
	t, err := time.Parse(formatFromCSV, m[0][len(m[0])-1])
	if err != nil {
		return date{}, err
	}
	return date(t), nil
}

func (m matrix) Validate() error {
	if len(m) < 2 {
		return errors.New("results should have at least 2 rows")
	}
	if len(m[0]) < 5 {
		return errors.New("results should have at least 5 columns")
	}
	return nil
}
