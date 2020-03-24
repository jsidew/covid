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
{{-      if eq .Status 1 }} resolving
{{- else if eq .Status 2 }} resolving slowly
{{- else if eq .Status 3 }} under control
{{- else if eq .Status 4 }} barely under control
{{- else if eq .Status 5 }} loosing control
{{- else if eq .Status 6 }} hard to control
{{- else }} out of control
{{- end }}. Active cases of #Covid_19 are compounding daily by {{ printf "en" "%.2f" .Current.Rate }}; as of {{ fmtdate "2 Jan 2006" .Updated }}, there are {{ print "en" .Current.Cases }} active cases. Projections: {{ printf "en" "%.0f" .Forecast.Cases }} cases ({{ .Forecast.Growth }}) in {{ print "en" .Forecast.Days }} days
{{- if lt .Current.Rate 1.0 -}}
; only 1 active case left in {{ print "en" .Recovery.DaysTo1 }} days
{{- end -}}
. [source: https://github.com/jsidew/covid]
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
