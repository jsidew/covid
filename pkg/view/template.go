package view

import (
	"text/template"
	"time"

	"golang.org/x/text/message"
)

const (
	// Extension used for template files
	Extension = ".tpl"

	// Name of the default template file, without extension.
	Name TemplateName = "default"

	// Template default content.
	Template = `{{ .Country }}:
{{-      if eq .Status.Score 1 }} resolving
{{- else if eq .Status.Score 2 }} resolving slowly
{{- else if eq .Status.Score 3 }} under control
{{- else if eq .Status.Score 4 }} barely under control
{{- else if eq .Status.Score 5 }} loosing control
{{- else if eq .Status.Score 6 }} hard to control
{{- else }} out of control
{{- end -}}
. #Covid_19 active cases {{ if lt .Current.Rate 1.0 }}dropping{{ else }}growing{{ end }} daily by {{ printf "en" "%.2f" .Current.Rate }}
{{- if .Status.Improving -}}
, w/dim factor of {{ printf "en" "%.3f" .Comparison.RateOfRates }}
{{- end -}}
. {{ print "en" .Current.Cases }} active cases, as of {{ fmtdate "2 Jan 2006" .Updated }}. Projection:
{{- if .Status.Improving }} recovering will start in {{ printf "en" "%.0f" .Recovery.DaysToStart }} days with a peak of {{ printf "en" "%.0f" .Recovery.PeakCases }} cases before it
{{- else }} {{ printf "en" "%.0f" .Forecast.Cases }} cases in {{ print "en" .Forecast.Days }} days
{{- end -}}
{{- if .Status.Resolving -}}
; only 1 active case left in {{ printf "en" "%.0f" .Recovery.DaysTo1 }} days
{{- end -}}
. [src: https://github.com/jsidew/covid]
`
)

var funcMap = template.FuncMap{
	"printf": func(lang string, format string, a ...interface{}) string {
		return message.NewPrinter(message.MatchLanguage(lang)).Sprintf(format, a...)
	},
	"print": func(lang string, a ...interface{}) string {
		return message.NewPrinter(message.MatchLanguage(lang)).Sprint(a...)
	},
	"fmtdate": func(layout string, t time.Time) string {
		return t.Format(layout)
	},
}
