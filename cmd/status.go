/*
Copyright © 2020 Jacopo Salvestrini <jsidew@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jsidew/covid/pkg/calc"
	"github.com/jsidew/covid/pkg/view"
)

const (
	dateLayout = "2006-01-02"
	fcastDays  = 30
)

type date time.Time

func init() {
	c := &statusCmd{}
	cmd := &cobra.Command{
		Use:   "status [COUNTRY]",
		Short: "Prints a tweet-long message about COVID-19 situation of the selected COUNTRY",
		Long: `Prints a tweet-long message about COVID-19 situation of the selected COUNTRY.

COUNTRY is one of the countries with cases as listed with the command 'covid countries';
to print the status of the whole world, either set COUNTRY to "world" or leave it empty.`,
		RunE: c.run,
		Args: cobra.MaximumNArgs(1),
	}
	flags := cmd.Flags()
	flags.Uint8VarP(&c.days, "days", "d", 7, "estimate for the last n days")
	flags.VarP(&c.since, "since", "s", "when to start the estimate with format: "+dateLayout)
	rootCmd.AddCommand(cmd)
}

type statusCmd struct {
	now, since date
	days       uint8
	country    string
}

func (c *statusCmd) run(_ *cobra.Command, args []string) error {
	view, err := view.New(profile, view.TemplateName(os.Getenv("COVID_TPL")))
	if err != nil {
		return err
	}

	if t, err := db.Latest(); err != nil {
		return err
	} else {
		c.now = date(t)
	}

	if len(args) > 0 {
		c.country = strings.TrimSpace(args[0])
	} else {
		c.country = "world"
	}
	if c.since.Time().IsZero() {
		c.since = c.now.AddDays(-int(c.days))
	}
	days := c.since.DaysSince(c.now)

	start, last, err := c.cases()
	if err != nil {
		return err
	}

	r := calc.Rate(float64(start), float64(last), float64(days))
	f := calc.Forecast(float64(last), r, fcastDays)
	good := calc.Period(float64(last), 1, r)

	var growth string
	g := (f/float64(last) - 1) * 100
	if g > 0 {
		growth = "+"
	}
	growth = fmt.Sprintf("%s%.0f%%", growth, g)

	err = view.Execute(os.Stdout, &struct {
		Country, ForecastGrowth string

		Rate float64

		ActiveCases, ForecastCases, ForecastDays, RecessionDays int

		UpdateDate time.Time
	}{
		Country:        strings.ToTitle(c.country),
		Rate:           r,
		ActiveCases:    int(last),
		UpdateDate:     c.now.Time(),
		ForecastCases:  int(f),
		ForecastGrowth: growth,
		ForecastDays:   fcastDays,
		RecessionDays:  int(good),
	})
	return err
}

func (c *statusCmd) cases() (start, last int, err error) {
	country := c.country
	if strings.EqualFold(country, "world") {
		country = ""
	}
	last, err = db.ActiveCases(country, c.now.Time())
	if err != nil {
		return
	}
	start, err = db.ActiveCases(country, c.since.Time())
	return
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

func (d date) Type() string {
	return "date"
}

func (d *date) Set(s string) error {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return err
	}
	*d = date(t)
	return nil
}
