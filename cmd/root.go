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

// Package cmd provides the command API for covid.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/jsidew/covid/pkg/calc"
	"github.com/jsidew/covid/pkg/database"
	"github.com/jsidew/covid/pkg/view"
)

var version = "1.1.0-alpha.1"

const (
	cacheExpire = 8 * time.Hour
	dbOrigin    = "https://raw.githubusercontent.com/bumbeishvili/covid19-daily-data"
)

var (
	profile = ".covid"
	rootCmd = &cobra.Command{
		Version: version,
		Use:     "covid",
		Short:   "Worldwide situation of COVID-19",
		Long:    `covid is a simple tool to undertand the COVID-19 current situation of countries around the world.`,
	}

	db *database.DB
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.covid.yaml)")
}

func initConfig() {
	home, err := homedir.Dir()
	exitif(err)

	profile = filepath.Join(home, profile)
	err = os.MkdirAll(profile, os.ModeDir|0700)
	exitif(err)

	db = database.New(dbOrigin, profile, cacheExpire)
	db.Set("confirmed", "/master/time_series_19-covid-Confirmed.csv")
	db.Set("recovered", "/master/time_series_19-covid-Recovered.csv")
	db.Set("dead", "/master/time_series_19-covid-Deaths.csv")
}

func exitif(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type country struct {
	name string

	days, compareDays uint8

	now, since, compare date
}

func (c *country) FillView(v *view.View) error {
	err := c.set()
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

	v.Country = strings.ToTitle(c.name)
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
	return nil
}

func (c *country) set() error {
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

	// setting country name
	if c.name == "" {
		c.name = "world"
	}

	return nil
}

func (c *country) cases() (pre, start, last int, err error) {
	country := c.name
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
