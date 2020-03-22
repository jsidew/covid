package database

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const formatFromCSV = "1/2/06" // month/day/year

type matrix [][]string

func (m matrix) Cases(country string, t time.Time) (int, error) {

	const coldatefrom = 4

	// find column index
	var colix int
	for i := range m[0][coldatefrom:] {
		u, err := time.Parse(formatFromCSV, m[0][i+coldatefrom])
		if err != nil {
			return 0, err
		}
		if !u.Equal(t) {
			continue
		}
		colix = i + coldatefrom
		break
	}

	// aggregate count by country
	var sum int
	for _, row := range m[1:] {
		if country != "" && !strings.EqualFold(country, strings.TrimSpace(row[1])) {
			continue
		}
		val := strings.TrimSpace(row[colix])
		if val == "" {
			continue
		}
		n, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		sum += n
	}

	return sum, nil
}

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

func (m matrix) cleancopy() matrix {
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
