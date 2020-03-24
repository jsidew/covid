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
	flags.Uint8VarP(&c.days, "days", "d", 7, "estimate for the last n days, define either this or --since")
	flags.Uint8VarP(&c.compareDays, "compareDays", "c", 0, "coparison estimate for the last n days, define either this or --compareSince (default is twice --days)")
	flags.VarP(&c.since, "since", "s", "when to start the estimate with format: "+dateLayout+", define either this or --days")
	flags.VarP(&c.compare, "compareSince", "a", "when to start the comparison estimate with format: "+dateLayout+", define either this or --compareDays")
	rootCmd.AddCommand(cmd)
}

type statusCmd struct {
	now, since, compare date
	days, compareDays   uint8
	country             string
}

func (c *statusCmd) run(_ *cobra.Command, args []string) error {
	v, err := view.New(profile, view.TemplateName(os.Getenv("COVID_TPL")))
	if err != nil {
		return err
	}

	err = c.set(args)
	if err != nil {
		return err
	}

	pre, start, last, err := c.cases()
	if err != nil {
		return err
	}

	r := calc.Rate(float64(start), float64(last), float64(c.days))
	f := calc.Forecast(float64(last), r, fcastDays)
	good := calc.Period(float64(last), 1, r)

	var growth string
	g := (f/float64(last) - 1) * 100
	if g > 0 {
		growth = "+"
	}
	growth = fmt.Sprintf("%s%.0f%%", growth, g)

	v.Country = strings.ToTitle(c.country)
	v.Updated = c.now.Time()
	v.Current.Rate = r
	v.Current.Cases = int(last)
	v.Recovery.DaysTo1 = good
	v.Forecast.Cases = f
	v.Forecast.Days = fcastDays
	v.Forecast.Growth = growth
	{
		r2 := calc.Rate(float64(pre), float64(last), float64(c.compareDays))
		r3 := calc.Rate(r2, r, float64(c.compareDays-c.days))
		recovery := calc.Period(r, 0.94, r3)
		peak := calc.Period(r, 1, r3)
		peakCases := calc.Forecast3D(float64(last), r, r3, peak)

		v.Comparison.Rate = r2
		v.Comparison.RateOfRates = r3
		v.Recovery.DaysToStart = recovery
		v.Recovery.DaysToPeak = peak
		v.Recovery.PeakCases = peakCases

		improving := r3 < 0.998
		status := view.OutOfControl

		if r < 0.94 {
			status = view.Resolving
		} else if r < 0.99 {
			status = view.ResolvingSlowly
		} else if r < 1.05 || (r < 1.09 && improving) {
			status = view.UnderControl
		} else if (r < 1.09 && !improving) || (r < 1.14 && improving) {
			status = view.BarelyUnderControl
		} else if r < 1.14 && !improving {
			status = view.LoosingControl
		} else if improving {
			status = view.HardToControl
		}

		v.Status.Score = status
		v.Status.Resolving = status == view.Resolving || status == view.ResolvingSlowly
		v.Status.Improving = improving && !v.Status.Resolving

	}

	err = v.Execute(os.Stdout)

	return err
}

func (c *statusCmd) set(args []string) error {
	// setting dates
	if t, err := db.Latest(); err != nil {
		return err
	} else {
		c.now = date(t)
	}
	if c.since.Time().IsZero() && c.days > 0 {
		c.since = c.now.AddDays(-int(c.days))
	}
	if !c.since.Time().IsZero() && c.days == 0 {
		c.days = c.now.DaysSince(c.since)
	}
	if c.compare.Time().IsZero() && c.compareDays == 0 {
		c.compareDays = c.days * 2
	}
	if c.compare.Time().IsZero() && c.compareDays > 0 {
		c.compare = c.now.AddDays(-int(c.compareDays))
	}
	if !c.compare.Time().IsZero() && c.compareDays == 0 {
		c.compareDays = c.now.DaysSince(c.compare)
	}

	// setting country
	if len(args) > 0 {
		c.country = strings.TrimSpace(args[0])
	} else {
		c.country = "world"
	}

	return nil
}

func (c *statusCmd) cases() (pre, start, last int, err error) {
	country := c.country
	if strings.EqualFold(country, "world") {
		country = ""
	}
	last, err = db.ActiveCases(country, c.now.Time())
	if err != nil {
		return
	}
	start, err = db.ActiveCases(country, c.since.Time())
	if err != nil {
		return
	}
	if !c.compare.Time().IsZero() {
		pre, err = db.ActiveCases(country, c.compare.Time())
	}
	return
}

func (d date) Time() time.Time {
	return time.Time(d)
}

func (d date) DaysSince(s date) uint8 {
	return uint8(d.Time().Sub(s.Time()).Minutes() / 60 / 24)
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
