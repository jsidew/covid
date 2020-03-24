package calc_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jsidew/covid/pkg/calc"
)

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

func Test(t *testing.T) {
	var (
		past, current float64 = 12543, 35456
		per1, per2    float64 = 13, 21
	)
	rate := calc.Rate(past, current, per1)
	final := calc.Forecast(current, rate, per2)
	assert.Equal(t, rate, calc.Rate(current, final, per2))
	assert.Equal(t, current, math.Round(calc.Forecast(past, rate, per1)))
	assert.Equal(t, per1, math.Round(calc.Period(past, current, rate)))
	assert.Equal(t, per2, math.Round(calc.Period(current, final, rate)))
}

func TestForecast3D(t *testing.T) {
	var a, p, r, r2 float64 = 735, 23, 1.11, 0.995

	r3 := r
	n := a
	for i := 1; i <= int(p); i++ {
		r3 *= r2
		n *= r3
	}

	f := calc.Forecast3D(a, r, r2, p)
	assert.Equal(t, fmt.Sprintf("%.4f", n), fmt.Sprintf("%.4f", f))
}

/*
The following tests don't check assertions.
They just print the output of the tested functions.
*/

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
