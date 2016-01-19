// Copyright (c) 2014 Datacratic. All rights reserved.

package vis

import (
	"html/template"
	"path/filepath"
	"strings"
)

var TemplateDir = "/home/michael/mygo/src/gotsvis/vis/"

func EscapeJQ(s string) string {
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, ",", "\\,", -1)
	s = strings.Replace(s, "(", "\\(", -1)
	s = strings.Replace(s, ")", "\\)", -1)
	return s
}

func LoadTemplates() *template.Template {
	pattern := filepath.Join(TemplateDir, "*.tmpl")
	funcMap := make(template.FuncMap)
	funcMap["jq"] = EscapeJQ
	funcMap["Chart"] = Chart
	funcMap["ChartTS"] = ChartSingle
	funcMap["ChartTSS"] = ChartSlice
	funcMap["TimeSeriesTagJS"] = TimeSeriesTagJS

	templates := template.Must(template.New("").Funcs(funcMap).ParseGlob(pattern))
	return templates
}
