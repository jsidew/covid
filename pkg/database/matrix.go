package database

import (
	"errors"
	"strconv"
	"time"
)

const formatFromCSV = "1/2/06" // month/day/year

type matrix [][]string

func (m matrix) Latest() (time.Time, error) {
	t, err := time.Parse(formatFromCSV, m[0][len(m[0])-1])
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func (m matrix) validate() error {
	if len(m) < 2 {
		return errors.New("results should have at least 2 rows")
	}
	if len(m[0]) < 5 {
		return errors.New("results should have at least 5 columns")
	}
	return nil
}

// Clean up from unwanted rows
func (m matrix) clean() matrix {
	new := append(matrix{}, m[0])
	for _, row := range m[1:] {
		_, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			continue
		}
		new = append(new, row)
	}
	return new
}
