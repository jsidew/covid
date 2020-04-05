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
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/jsidew/covid/pkg/view"
)

const (
	dateLayout = "2006-01-02"
	fcastDays  = 30
)

type date time.Time

func init() {
	c := &statusCmd{
		days: 7,
	}
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
	flags.Uint8VarP(&c.days, "days", "d", c.days, "estimate for the last n days, define either this or --since")
	flags.Uint8VarP(&c.compareDays, "compareDays", "c", c.compareDays, "coparison estimate for the last n days, define either this or --compareSince (default is twice --days)")
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

	if len(args) > 0 {
		c.country = strings.TrimSpace(args[0])
	}

	country := &country{
		name:        c.country,
		days:        c.days,
		compareDays: c.compareDays,
		now:         c.now,
		since:       c.since,
		compare:     c.compare,
	}

	err = country.FillView(v)
	if err != nil {
		return err
	}

	err = v.Execute(os.Stdout)

	return err
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
