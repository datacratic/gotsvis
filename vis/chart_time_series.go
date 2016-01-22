// Copyright (c) 2014 Datacratic. All rights reserved.

package vis

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"math"
	"strconv"

	"github.com/datacratic/gotsvis/ts"
)

func timeString(ts *ts.TimeSeries) (string, error) {
	s := bytes.NewBufferString("[ 'x', ")

	cursor := ts.Start()
	step := ts.Step()
	// TODO: change to use start, end, step as iterator instead of a data slice
	for _, _ = range ts.Data() {
		if err := s.WriteByte('\''); err != nil {
			return "", err
		}
		if _, err := s.WriteString(cursor.Format("20060102 15:04:05")); err != nil {
			return "", err
		}
		if _, err := s.WriteString("',"); err != nil {
			return "", err
		}
		cursor = cursor.Add(step)
	}
	if s.Len() > 0 {
		s.Truncate(s.Len() - 1)
	}
	if err := s.WriteByte(']'); err != nil {
		return "", err
	}
	return s.String(), nil
}

func valueString(ts *ts.TimeSeries) (string, error) {
	s := bytes.NewBufferString("[ '")
	if _, err := s.WriteString(ts.Key()); err != nil {
		return "", err
	}
	if _, err := s.WriteString("', "); err != nil {
		return "", err
	}

	for _, v := range ts.Data() {
		if math.IsInf(v, 0) {
			if _, err := s.WriteString("NaN"); err != nil {
				return "", err
			}
		} else if _, err := s.WriteString(strconv.FormatFloat(v, 'E', 2, 64)); err != nil {
			return "", err
		}
		if err := s.WriteByte(','); err != nil {
			return "", err
		}
	}
	if s.Len() > 0 {
		s.Truncate(s.Len() - 1)
	}
	if err := s.WriteByte(']'); err != nil {
		return "", err
	}
	return s.String(), nil
}

func ChartSingle(ts *ts.TimeSeries) (template.JS, error) {
	t, err := timeString(ts)
	if err != nil {
		return "", err
	}
	fmt.Println(ts.Key())
	fmt.Println(ts.Start(), ts.End())
	s := bytes.NewBufferString(t)
	if err := s.WriteByte(','); err != nil {
		return "", err
	}

	v, err := valueString(ts)
	if err != nil {
		return "", err
	}
	if _, err := s.WriteString(v); err != nil {
		return "", err
	}
	return template.JS(s.String()), nil
}

func ChartSlice(tss ts.TimeSeriesSlice) (template.JS, error) {
	if tss == nil {
		return template.JS("[]"), nil
	}
	start := tss.Start()
	end := tss.End()
	step, ok := tss.Step()
	if !ok {
		return "", errors.New("time series step is not equal")
	}
	fmt.Println(tss.Key())
	fmt.Println(tss.Start(), tss.End())
	dummyTS, err := ts.NewTimeSeriesOfTimeRange("dummyTS", start, end, step, 0)
	if err != nil {
		return "", err
	}

	t, err := timeString(dummyTS)
	if err != nil {
		return "", err
	}
	s := bytes.NewBufferString(t)

	for i, _ := range tss {
		if err := s.WriteByte(','); err != nil {
			return "", err
		}

		v, err := valueString(&tss[i])
		if err != nil {
			return "", err
		}
		if _, err := s.WriteString(v); err != nil {
			return "", err
		}
	}
	return template.JS(s.String()), nil
}

func Chart(series interface{}) (template.JS, error) {
	switch t := series.(type) {
	case *ts.TimeSeries:
		return ChartSingle(t)
	case ts.TimeSeries:
		return ChartSingle(&t)
	case ts.TimeSeriesSlice:
		return ChartSlice(t)
	default:
		return "", fmt.Errorf("unknown type '%T'", series)
	}
}

func TimeSeriesTagJS(series interface{}) (template.JS, error) {
	switch t := series.(type) {
	case *ts.TimeSeries:
		return template.JS(t.Key()), nil
	case ts.TimeSeries:
		return template.JS(t.Key()), nil
	case ts.TimeSeriesSlice:
		return template.JS(t.Key()), nil
	default:
		return "", fmt.Errorf("unknown type '%T'", series)
	}
}
