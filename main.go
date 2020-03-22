package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/message"
)

var version = "0.1.0"

const (
	dateLayout  = "2006-01-02"
	dfltTplName = "default.tpl"
	fcastDays   = 30
	cacheExpire = 8 * time.Hour
	httpTimeout = 10 * time.Second
)

var (
	argDays uint = 7

	argCountry string
	argSince   date

	tpl template.Template

	db = ".covid"

	funcMap = template.FuncMap{
		"printf": func(lang string, format string, a ...interface{}) string {
			p := message.NewPrinter(message.MatchLanguage(lang))
			return p.Sprintf(format, a...)
		},
		"print": func(lang string, a ...interface{}) string {
			p := message.NewPrinter(message.MatchLanguage(lang))
			return p.Sprint(a...)
		},
		"fmtdate": func(layout string, d date) string {
			return d.Time().Format(layout)
		},
	}
)

func init() {

	// config folder
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	db = filepath.Join(home, db)
	err = os.MkdirAll(db, os.ModeDir|0700)
	if err != nil {
		panic(err)
	}

	// command flags
	flag.StringVar(&argCountry, "country", argCountry, "name of the country")
	flag.Var(&argSince, "since", "date to start the estimate - format: "+dateLayout)
	flag.UintVar(&argDays, "days", argDays, "estimate for the last n days")

	// global settings
	http.DefaultClient.Timeout = httpTimeout

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

	if flag.Arg(0) == "countries" {
		err := database.Load("confirmed")
		exitif(err)
		for _, country := range database.Get("confirmed").M.Countries() {
			fmt.Println(country)
		}
		os.Exit(0)
	}

	err := loadTemplate(dfltTplName, &tpl)
	exitif(err)

	err = database.Load("confirmed", "recovered", "dead")
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
	f := forecast(last, r, fcastDays)
	good := recession(last, r, 1)

	var growth string
	g := (f/last - 1) * 100
	if g > 0 {
		growth = "+"
	}
	growth = fmt.Sprintf("%s%.0f%%", growth, g)

	country := "WORLD"
	if argCountry != "" {
		country = strings.ToTitle(argCountry)
	}

	err = tpl.Execute(os.Stdout, &struct {
		Country, ForecastGrowth string

		Rate float64

		ActiveCases, ForecastCases, ForecastDays, RecessionDays int

		UpdateDate date
	}{
		Country:        country,
		Rate:           r,
		ActiveCases:    int(last),
		UpdateDate:     now,
		ForecastCases:  int(f),
		ForecastGrowth: growth,
		ForecastDays:   fcastDays,
		RecessionDays:  int(good),
	})
	exitif(err)
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

func loadTemplate(name string, tpl *template.Template) error {
	err := func(name string) error {
		f, err := os.OpenFile(filepath.Join(db, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.WriteString(defaultTpl)
		return err
	}(name)
	if err != nil && !os.IsExist(err) {
		return err
	}

	t, err := template.New("root").Funcs(funcMap).ParseGlob(filepath.Join(db, "*.tpl"))
	if err != nil {
		return err
	}

	if v := strings.TrimSpace(os.Getenv("COVID_TPL")); v != "" {
		if !strings.HasSuffix(v, "tpl") {
			v += ".tpl"
		}
		name = v
	}
	t = t.Lookup(name)
	if t == nil {
		return fmt.Errorf("template \"%s\" doesn't exist%s", name, t.DefinedTemplates())
	}

	*tpl = *t

	return nil
}

func exitif(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
