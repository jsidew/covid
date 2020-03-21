package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var version = "0.0.3"

const (
	dateLayout      = "2006-01-02"
	refreshInterval = 12 * time.Hour
)

var (
	db = ".covid"

	argDays uint = 1

	argCountry string
	argSince   date
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	db = filepath.Join(home, db)

	err = os.MkdirAll(db, os.ModeDir|0700)
	if err != nil {
		panic(err)
	}

	flag.StringVar(&argCountry, "country", argCountry, "name of the country")
	flag.Var(&argSince, "since", "date to start the estimate - format: "+dateLayout)
	flag.UintVar(&argDays, "days", argDays, "estimate for the last n days - default to 1 day")

}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		fmt.Printf("covid v%s\n", version)
		os.Exit(0)
	}

	database := resources{}
	database.Set(db, "confirmed", "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Confirmed.csv")
	database.Set(db, "recovered", "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Recovered.csv")
	database.Set(db, "dead", "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Deaths.csv")

	err := database.Load("confirmed", "recovered", "dead")
	exitif(err)

	now, err := database.Get("confirmed").M.LastDate()
	exitif(err)

	if argSince.Time().IsZero() {
		argSince = now.AddDays(-int(argDays))
	}
	days := argSince.DaysSince(now)

	last, err := actives(argCountry, now, database.Get("confirmed").M, database.Get("recovered").M, database.Get("dead").M)
	exitif(err)

	start, err := actives(argCountry, argSince, database.Get("confirmed").M, database.Get("recovered").M, database.Get("dead").M)
	exitif(err)

	r := rate(start, last, days)
	f := forecast(last, r, 30)

	var growth string
	g := (f/last - 1) * 100
	if g > 0 {
		growth = "+"
	}
	growth = fmt.Sprintf("%s%.0f%%", growth, g)

	fmt.Printf("%s's active spread rate is at %.2f;", strings.ToTitle(argCountry), r)
	fmt.Printf(" now there are %.0f active cases [%s];\n", last, now)
	fmt.Printf("at the current rate, there will be %.0f active cases (%s) within the next 30 days", f, growth)
	if r < 1 {
		good := recession(last, r, 1)
		fmt.Printf(", and only 1 active case will be left after %.0f days", good)
	}
	fmt.Println(".")
}

func actives(country string, d date, confirmed matrix, subtracted ...matrix) (cases float64, err error) {
	var c int
	c, err = confirmed.Cases(country, d)
	if err != nil {
		return
	}
	for _, sub := range subtracted {
		r, err := sub.Cases(country, d)
		if err != nil {
			return 0, err
		}
		c -= r
	}
	cases = float64(c)
	return
}

func rate(start, last, days float64) float64 {
	return math.Pow(last/start, 1/days)
}

func forecast(current, rate, days float64) float64 {
	return current * math.Pow(rate, days)
}

func recession(current, rate, end float64) float64 {
	return math.Log(end/current) / math.Log(rate)
}

func exitif(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
