package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"sort"
)

const formatFromCSV = "1/2/06" // month/day/year

type matrix [][]string

func (m matrix) Cases(country string, d date) (int, error) {

	const coldatefrom = 4

	// find column index
	var colix int
	for i := range m[0][coldatefrom:] {
		t, err := time.Parse(formatFromCSV, m[0][i+coldatefrom])
		if err != nil {
			return 0, err
		}
		if !t.Equal(d.Time()) {
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

func (m matrix) Countries() []string {
	c := map[string]struct{}{}

	for _, row := range m[1:] {
		c[strings.TrimSpace(row[1])] = struct{}{}
	}

	list := []string{}
	for k := range c {
		list = append(list, k)
	}
	sort.Strings(list)

	return list
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

// Clean up from unwanted rows
func (m matrix) Clean() matrix {
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
