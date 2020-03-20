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

// TODO sum (US, UK, etc.)

var version = "0.0.1"

const (
	dateLayout      = "2006-01-02"
	refreshInterval = 12 * time.Hour
)

var (
	db = ".covid"

	country string
	from    date
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

	flag.StringVar(&country, "country", country, "name of the country")
	flag.Var(&from, "from", "date to start the estimate from - format: "+dateLayout)

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
	database.Set(db, "deaths", "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Deaths.csv")

	err := database.Load("confirmed", "recovered")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for _, country := range strings.Split(country, ",") {
		start, last, days, err := actives(country, database.Get("confirmed").M, database.Get("recovered").M)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Printf(
			"%s's spread rate ratio is at %.2f, from %.0f to %.0f active cases, over %.0f days.\n",
			strings.Title(country), ratio(start, last, days), start, last, days)
	}
}

func actives(country string, confirmed, recovered matrix) (start, last, days float64, err error) {
	c_start, err := confirmed.Cases(country, from)
	if err != nil {
		return
	}
	c_end, latest, err := confirmed.CasesLatest(country)
	if err != nil {
		return
	}
	r_start, err := recovered.Cases(country, from)
	if err != nil {
		return
	}
	r_end, err := recovered.Cases(country, latest)
	if err != nil {
		return
	}
	start = float64(c_start - r_start)
	last = float64(c_end - r_end)
	days = latest.DaysSince(from)
	return
}

func ratio(start, last, days float64) float64 {
	return math.Pow(last/start, 1/days)
}
