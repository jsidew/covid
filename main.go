package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jsidew/covid/pkg/database"
	"github.com/jsidew/covid/pkg/view"
)

var version = "0.2.0"

const (
	dateLayout  = "2006-01-02"
	fcastDays   = 30
	cacheExpire = 8 * time.Hour
)

type date time.Time

var (
	argDays uint = 7

	argCountry string
	argSince   date

	profile = ".covid"
)

func init() {

	// config folder
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	profile = filepath.Join(home, profile)
	err = os.MkdirAll(profile, os.ModeDir|0700)
	if err != nil {
		panic(err)
	}

	// command flags
	flag.StringVar(&argCountry, "country", argCountry, "name of the country")
	flag.Var(&argSince, "since", "date to start the estimate - format: "+dateLayout)
	flag.UintVar(&argDays, "days", argDays, "estimate for the last n days")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		fmt.Printf("covid v%s\n", version)
		os.Exit(0)
	}

	db := database.New("https://raw.githubusercontent.com/CSSEGISandData/COVID-19", profile, cacheExpire)
	db.Set("confirmed", "/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Confirmed.csv")
	db.Set("recovered", "/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Recovered.csv")
	db.Set("dead", "/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Deaths.csv")

	if flag.Arg(0) == "countries" {
		countries, err := db.Countries()
		exitif(err)
		for _, country := range countries {
			fmt.Println(country)
		}
		os.Exit(0)
	}

	view, err := view.New(profile, view.TemplateName(os.Getenv("COVID_TPL")))
	exitif(err)

	var now date
	{
		l, err := db.Latest()
		exitif(err)
		now = date(l)
	}

	if argSince.Time().IsZero() {
		argSince = now.AddDays(-int(argDays))
	}
	days := argSince.DaysSince(now)

	last, err := db.ActiveCases(argCountry, now.Time(), "confirmed", "recovered", "dead")
	exitif(err)

	start, err := db.ActiveCases(argCountry, argSince.Time(), "confirmed", "recovered", "dead")
	exitif(err)

	r := rate(float64(start), float64(last), float64(days))
	f := forecast(float64(last), r, fcastDays)
	good := recession(float64(last), r, 1)

	var growth string
	g := (f/float64(last) - 1) * 100
	if g > 0 {
		growth = "+"
	}
	growth = fmt.Sprintf("%s%.0f%%", growth, g)

	country := "WORLD"
	if argCountry != "" {
		country = strings.ToTitle(argCountry)
	}

	err = view.Execute(os.Stdout, &struct {
		Country, ForecastGrowth string

		Rate float64

		ActiveCases, ForecastCases, ForecastDays, RecessionDays int

		UpdateDate time.Time
	}{
		Country:        country,
		Rate:           r,
		ActiveCases:    int(last),
		UpdateDate:     now.Time(),
		ForecastCases:  int(f),
		ForecastGrowth: growth,
		ForecastDays:   fcastDays,
		RecessionDays:  int(good),
	})
	exitif(err)
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

func (d date) Time() time.Time {
	return time.Time(d)
}

func (d date) DaysSince(s date) float64 {
	return -d.Time().Sub(s.Time()).Minutes() / 60 / 24
}

func (d date) AddDays(days int) date {
	return date(d.Time().Add(time.Duration(days) * 24 * time.Hour))
}

func (d date) String() string {
	return d.Time().Format(dateLayout)
}

func (d *date) Set(s string) error {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	*d = date(t)
	return nil
}

func exitif(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
