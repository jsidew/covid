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

TODO

## Virus Control Scale (VCS) Algorithm Explained

TODO
