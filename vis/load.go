// Copyright (c) 2014 Datacratic. All rights reserved.

package vis

import (
	"html/template"
	"strings"
)

func EscapeJQ(s string) string {
	s = strings.Replace(s, ".", "\\.", -1)
	s = strings.Replace(s, ",", "\\,", -1)
	s = strings.Replace(s, "(", "\\(", -1)
	s = strings.Replace(s, ")", "\\)", -1)
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	s = strings.Replace(s, ":", "\\:", -1)
	return s
}

func LoadTemplates() (*template.Template, error) {
	funcMap := make(template.FuncMap)
	funcMap["jq"] = EscapeJQ
	funcMap["Chart"] = Chart
	funcMap["ChartTS"] = ChartSingle
	funcMap["ChartTSS"] = ChartSlice
	funcMap["TimeSeriesTagJS"] = TimeSeriesTagJS

	templates := template.New("").Funcs(funcMap)
	for _, tmpl := range GetTMPL() {
		if _, err := templates.Parse(tmpl); err != nil {
			return nil, err
		}
	}
	return templates, nil
}
