package calc_test

import (
	"testing"

	"github.com/jsidew/covid/pkg/calc"
)

/*
These tests don't check assertions.
They just print the output of the tested functions.
*/

var table = []struct {
	first, second, third float64
}{
	{0, 0, 0},
	{0, 123, 20},
	{123, 0, 20},
	{123, 1234, 0},
	{123, 1234, 1},
	{1234, 123, 1},
	{1234, 123, 20},
	{1, 1, 1},
	{1234, 1, 123},
	{1, 123, 1234},
	{-4, 12, 34},
	{4, -12, 34},
	{4, 12, -34},
	{12, 4, 34},
}

func TestRate(t *testing.T) {
	t.Log("past, current, period")
	for i, test := range table {
		r := calc.Rate(test.first, test.second, test.third)
		t.Logf("%d: %.2f, %.2f, %.2f = %f", i, test.first, test.second, test.third, r)
	}
}

func TestForecast(t *testing.T) {
	t.Log("current, rate, period")
	for i, test := range table {
		f := calc.Forecast(test.first, test.second, test.third)
		t.Logf("%d: %.2f, %.2f, %.2f = %f", i, test.first, test.second, test.third, f)
	}
}

func TestPeriod(t *testing.T) {
	t.Log("current, final, rate")
	for i, test := range table {
		d := calc.Period(test.first, test.second, test.third)
		t.Logf("%d: %.2f, %.2f, %.2f = %f", i, test.first, test.second, test.third, d)
	}
}
