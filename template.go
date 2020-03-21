package main

var defaultTpl = `{{ .Country }}:
{{-      if lt .Rate 0.95 }} resolving
{{- else if lt .Rate 1.00 }} resolving slowly
{{- else if lt .Rate 1.05 }} under control
{{- else if lt .Rate 1.10 }} under control
{{- else if lt .Rate 1.15 }} barely under control
{{- else if le .Rate 1.20 }} hard to control
{{- else }} out of control
{{- end }}. Active cases of #Covid_19 are compounding daily by {{ printf "en" "%.2f" .Rate }}; as of {{ fmtdate "2 Jan 2006" .UpdateDate }}, there are {{ print "en" .ActiveCases }} active cases. Projections: {{ print "en" .ForecastCases }} cases ({{ .ForecastGrowth }}) in {{ print "en" .ForecastDays }} days
{{- if lt .Rate 1.0 -}}
; only 1 active case left in {{ print "en" .RecessionDays }} days
{{- end -}}
. [source: https://github.com/jsidew/covid]
`
