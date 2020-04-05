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
	"sort"
	"strconv"

	"github.com/jsidew/covid/pkg/view"
	"github.com/spf13/cobra"
)

func init() {
	c := &topCmd{}

	cmd := &cobra.Command{
		Use:   "top",
		Short: "Ordered list of top countries that have the best control",
		Long:  `Ordered list of top countries that have the best control.`,
		RunE:  c.run,
	}

	cmd.Flags().BoolVarP(&c.worse, "worse", "w", c.worse, "Sort by worse")

	rootCmd.AddCommand(cmd)

}

type countryList []view.View

type topCmd struct {
	worse bool
}

func (c *topCmd) run(*cobra.Command, []string) error {
	countryNames, err := db.Countries()
	if err != nil {
		return err
	}
	var longest int
	countries := countryList{}
	for _, name := range countryNames {
		if l := len(name); l > longest {
			longest = l
		}
		country := &country{
			name: name,
			days: 7,
		}
		v := view.View{}
		err := country.FillView(&v)
		if err != nil {
			return err
		}
		countries = append(countries, v)
	}

	if c.worse {
		sort.Sort(sort.Reverse(countries))
	} else {
		sort.Sort(countries)
	}

	l := strconv.Itoa(longest)
	f1 := "%3s %-" + l + "s %5s %6s    %s\n"
	f2 := "%3d %-" + l + "s %3d %8d    %.2fx%.4f -- %.0f\n"
	fmt.Printf(f1, "#", "country", "score", "cases", "Rate x Rate²")
	printn("-", 3+1+longest+1+5+1+6+4+12)
	for i, country := range countries {
		if country.Current.Cases < 1 {
			continue
		}
		fmt.Printf(f2, i+1,
			country.Country,
			country.Status.Score,
			country.Current.Cases,
			country.Current.Rate,
			country.Comparison.RateOfRates,
			country.Recovery.PeakCases)
	}
	return nil
}

func printn(s string, n int) {
	var ss string
	for i := 0; i < n; i++ {
		ss += s
	}
	fmt.Println(ss)
}

func (c countryList) Len() int {
	return len(c)
}

func (c countryList) Less(i, j int) bool {
	if c[i].Status.Score < c[j].Status.Score {
		return true
	}
	if c[i].Status.Score > c[j].Status.Score {
		return false
	}

	score := c[i].Status.Score
	if score > 2 && score < 6 {
		return (c[i].Recovery.PeakCases / float64(c[i].Current.Cases)) <
			(c[j].Recovery.PeakCases / float64(c[j].Current.Cases))

	}

	return (c[i].Current.Rate < c[j].Current.Rate) ||
		(c[i].Current.Rate == c[j].Current.Rate && c[i].Comparison.RateOfRates < c[j].Comparison.RateOfRates)
}

func (c countryList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
