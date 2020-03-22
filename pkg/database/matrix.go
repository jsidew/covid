package database

import (
	"encoding/csv"
	"errors"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
)

const formatFromCSV = "1/2/06" // month/day/year

type matrix [][]string

func newMatrix(r io.Reader) (matrix, error) {
	results, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, err
	}
	if len(results) < 2 {
		return nil, errors.New("results should have at least 2 rows")
	}
	if len(results[0]) < 5 {
		return nil, errors.New("results should have at least 5 columns")
	}
	m := append(matrix{}, results[0])
	for _, row := range results[1:] {
		_, err := strconv.ParseFloat(row[3], 64)
		if err != nil {
			continue
		}
		m = append(m, row)
	}
	return m, nil
}

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
