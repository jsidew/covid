package calc

import "math"

// Rate at which a current number is growing, coming from a past number in a given period (days, weeks, years, ...).
func Rate(past, current, period float64) float64 {
	return math.Pow(current/past, 1/period)
}

// Forecast of how much a current number will grow at a given rate after a certain period (days, weeks, years, ...).
func Forecast(current, rate, period float64) float64 {
	return current * math.Pow(rate, period)
}

// Period (days, weeks, years, ...) that a current number will take to reach a final number at a given growth rate.
func Period(current, final, rate float64) float64 {
	return math.Log(final/current) / math.Log(rate)
}
