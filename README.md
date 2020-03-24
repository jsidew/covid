# covid

`covid` is an attempt that wants to illustrate quickly and effectively how well (or badly) different countries are controlling the COVID-19 spread. For that I developed what I call the Virus Control Scale, ranging from 1 to 7, where 7 is the worst.

## Acknowledgments

* Inspired by body.io's podcast
   "[The Sky Is Falling: Grab a Corona and Relax](https://body.io/the-sky-is-falling-grab-a-corona-and-relax/)";
* data taken from [CSSEGISandData/COVID-19](https://github.com/CSSEGISandData/COVID-19).


## Installation

Requires [GO v1.13](https://golang.org/dl/)+

```
$ go get -u github.com/jsidew/covid
```


## Usage

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

## Template

TODO

## Virus Control Scale (VCS) Algorithm Explained

TODO
