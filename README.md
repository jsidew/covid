# covid

`covid` is a tool to educate people (governments included - some I think they are still underestimating) about how bad it can get if the virus spread stays uncontrolled. You can use this tool to generate messages to be shared on social media platforms. You can customise your message, also in your language of choice. Contributions and ideas are very welcome. In other words, the purpose is to use numbers to help monitoring quickly and effectively how well (or badly) different countries are controlling the COVID-19 spread.

[LinkedIn Article about Covid]()

`covid` is written in [Golang](https://golang.org/), and works around what I call the Virus Control Scale, ranging from 1 to 7, where 7 is the worst.

## Contents
* [Acknowledgments](#acknowledgments)
* [Installation](#installation)
* [Usage](#usage)
* [Customise](#customise)
* [Virus Control Scale (VCS) Algorithm Explained](#virus-control-scale-vcs-algorithm-explained)

## Acknowledgments

Inspired by body.io's podcast
   "[The Sky Is Falling: Grab a Corona and Relax](https://body.io/the-sky-is-falling-grab-a-corona-and-relax/)".

Data source: [bumbeishvili/covid19-daily-data](https://github.com/bumbeishvili/covid19-daily-data).

## Installation

Requires [GO v1.13](https://golang.org/dl/)+

```
$ go get -u github.com/jsidew/covid
```

## Usage

### Example: Defaults

This case shows Italy's status. By default the spread rate is calculated over the last 7 days (1.10 below), while the control rate (the "dim factor" below) is calculated comparing the rate of the last 14 days with the previous one. "Recovering" is when the spread rate is below 1; i.e. when the active cases are decreasing. The peak is the maximum number of cases that could be reached, with such spread and control rates.
```
$ covid status italy
ITALY: barely under control.
#Covid_19 active cases growing daily by 1.10, w/dim factor of 0.997.
57,521 active cases, as of 25 Mar 2020.
Projection: recovering will start in 52 days with a peak of 270,215 cases before it.
```

### Example: Custom Days 

You can change the last days to consider for the spread rate and the last days for the control rate. As you notice, the output is slightly different.
```
$ covid status italy -d 14 -c 28
ITALY: barely under control.
#Covid_19 active cases growing daily by 1.13, w/dim factor of 0.996.
57,521 active cases, as of 25 Mar 2020.
Projection: recovering will start in 48 days with a peak of 368,528 cases before it.
```

### Example: Set a Date

Instead of setting days you can set a date, for example if after that date there was an important decision from a country's government.
```
$ covid status 'Czech Republic' -a '2020-03-12'
CZECH REPUBLIC: hard to control.
#Covid_19 active cases growing daily by 1.20, w/dim factor of 0.994.
1,638 active cases, as of 25 Mar 2020.
Projection: recovering will start in 38 days with a peak of 18,999 cases before it.
```

### Example: Out of Control

When the status is "out of control", the control rate is not showed and the projection doesn't consider it and it's made entirely from the spread rate (1.36 below).
```
$ covid status usa
USA: out of control.
#Covid_19 active cases growing daily by 1.36.
66,791 active cases, as of 25 Mar 2020.
Projection: 714,778,230 cases in 30 days.
```

### Example: Recovering

Also when the status is "recovering", the control rate is not considered. It's also shown how many days there will be left only 1 active case.
```
$ covid status china
CHINA: resolving.
#Covid_19 active cases dropping daily by 0.91.
4,290 active cases, as of 25 Mar 2020.
Projection: 281 cases in 30 days; only 1 active case left in 92 days.
```


### Help

```
$ covid help
covid is a simple tool to undertand the COVID-19 current situation of countries around the world.

Usage:
  covid [command]

Available Commands:
  countries   List names of the countries with COVID-19 cases
  help        Help about any command
  status      Prints a tweet-long message about COVID-19 situation of the selected COUNTRY
  version     Prints covid's version

Flags:
  -h, --help   help for covid

Use "covid [command] --help" for more information about a command.
```

## Customise

You can customise the message output from `covid status`. Under your home folder (`~` or `$HOME` under Unix OSs like Mac and Linux; `%userprofile%` for Windows) there is a folder that is created the first time you run `covid`. The folder is called `.covid`. Under this folder, there is a file named `default.tpl`, which you can change as you like, following [Golang's text/template](https://pkg.go.dev/text/template?tab=doc) syntax. You can also create multiple templates (with exptension `.tpl`) and choose the one you prefer with the environment variable `COVID_TPL`:
```
$ COVID_TPL=test covid status
Hello WORLD: 7@1.1540755042768218x1.0018865162982553
```

### Parameters

Supported parameters are
* `.Country`, _string_, the name of the country;
* `.Updated`, _time.Time_, the date when the data source was last updated;
* `.Status.Score`, _uint8_, the score (from 1 to 7) of the VCS;
* `.Status.Resolving`, _bool_, if the situation is resolving;
* `.Status.Improving`, _bool_, if the situation is improving (note that it's not resolving, the spread is still growing, but less day by day);
* `.Current.Rate`, _float64_, the current spread rate;
* `.Current.Cases`, _int_, the latest number of active cases as stored in the data source;
* `.Comparison.Rate`, _float64_, the spread rate from greater time-span;
* `.Comparison.RateOfRates`, _float64_, the rate of growth between `.Current.Cases` and `.Comparison.Rate` (if the current rate is lower than the rate from a greater time span, than the situation is improving and under control);
* `.Recovery.DaysTo1`, _float64_, number of days necessary to have only 1 active case left;
* `.Recovery.DaysToStart`, _float64_, number of days necessary to start recovering, considering `.Comparison.RateOfRates`;
* `.Recovery.DaysToPeak`, _float64_, number of days necessary to reach a peak of active cases (right before resolution starts), considering `.Comparison.RateOfRates`;
* `.Recovery.PeakCases`, _float64_, peak number of active cases, considering `.Comparison.RateOfRates`;
* `.Forecast.Cases`, _float64_, number of cases that will be reached after `.Forecast.Days` at `.Current.Rate`;
* `.Forecast.Days`, _int_, number of days considered to reach `.Forecast.Cases`.

### Functions

Supported function, other than Golang standard ones, are
* `printf`(_lang string, format string, a ...interface{}_), format a list of values according to a [formatting syntax](https://pkg.go.dev/golang.org/x/text/message?tab=doc) and language (e.g. "en", "it", etc.);
* `print`(_lang string, a ...interface{}_), format a list of values according to a [language](https://pkg.go.dev/golang.org/x/text/message?tab=doc) (e.g. "en", "it", etc.);
* `fmtdate`(_layout string, t time.Time_), format a time, like `.Updated` (see previous paragraph), according to [Time.Format](https://pkg.go.dev/time?tab=doc#Time.Format);

## Virus Control Scale (VCS) Algorithm Explained

### Variables
* `c`: current number of (active) cases
* `d`: number of cases in the past
* `p`: the period (e.g. 5 days) between `c` and `d`
* `e`: number of cases in a more remote past
* `q`: the period between `c` and `e`

```
r = (c/d)^(1/p)
s = (c/e)^(1/q)
x = (r/s)^(1/(q-p))
```
* `r`: spread rate
* `s`: spread rate from earlier cases
* `x`: rate of rates, or control rate

### The Scale

| x >= 0.998 | r < of: | x < 0.998 |
|------------|---------|-----------|
| (1) | 0.94 | (1) |
| (2) | 0.99 | (2) |
| (3) | 1.05 | (3) |
| (4) | 1.09 | (3) |
| (5) | 1.14 | (4) |
| (7) | 1.20 | (6) |
| (7) | r >= 1.20 | (6) |

### Attributes

1. Resolving
2. Resolving Slowly
3. Under Control
4. Barely Under Control
5. Loosing Control
6. Hard to Control
7. Out of Control

### Forecast with Moving Rate

* `f`: cases in the future 
* `z`: period in the future between `c` and `f`

```

f = c * r^z * x^((z*z+z)/2)
```
