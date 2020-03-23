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
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	"github.com/jsidew/covid/pkg/database"
)

var version = "0.3.0"

const (
	cacheExpire = 8 * time.Hour
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

	db = database.New("https://raw.githubusercontent.com/CSSEGISandData/COVID-19", profile, cacheExpire)
	db.Set("confirmed", "/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Confirmed.csv")
	db.Set("recovered", "/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Recovered.csv")
	db.Set("dead", "/master/csse_covid_19_data/csse_covid_19_time_series/time_series_19-covid-Deaths.csv")
}

func exitif(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
