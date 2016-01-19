// Copyright (c) 2014 Datacratic. All rights reserved.

package ts

import (
	"bytes"
	"html/template"
	"math"
	"strconv"
	"time"
)

type TimeSeriesSlice []TimeSeries

func (tss TimeSeriesSlice) FilterKeep(filter Filter) TimeSeriesSlice {
	keepTSS := make(TimeSeriesSlice, 0)
	for _, ts := range tss {
		keep := false
		for _, v := range ts.data {
			keep = filter.Filter(v)
			if keep {
				keepTSS = append(keepTSS, ts)
				break
			}
		}
	}

	return keepTSS
}

func (tss TimeSeriesSlice) GetKey(key string) *TimeSeries {
	for _, ts := range tss {
		if ts.key == key {
			return &ts
		}
	}
	return nil
}

func (tss TimeSeriesSlice) getStartEnd() (start, end time.Time) {
	if len(tss) > 0 {
		start = tss[0].start
		end = tss[0].End()
	}

	for i, _ := range tss {
		if start.After(tss[i].start) {
			start = tss[i].start
		}
		if end.Before(tss[i].End()) {
			end = tss[i].End()
		}
	}
	return
}

func (tss TimeSeriesSlice) checkEqualStep() (time.Duration, bool) {
	var step time.Duration
	if len(tss) > 0 {
		step = tss[0].step
	}
	for i := 1; i < len(tss); i++ {
		if step != tss[i].step {
			return 0, false
		}
	}
	return step, true
}

func (tss TimeSeriesSlice) Start() time.Time {
	var start time.Time
	if len(tss) > 0 {
		start = tss[0].start
	}

	for i, _ := range tss {
		if start.After(tss[i].start) {
			start = tss[i].start
		}
	}
	return start
}

func (tss TimeSeriesSlice) End() time.Time {
	var end time.Time
	if len(tss) > 0 {
		end = tss[0].End()
	}

	for i, _ := range tss {
		if end.Before(tss[i].End()) {
			end = tss[i].End()
		}
	}
	return end
}

func (tss TimeSeriesSlice) Step() (time.Duration, bool) {
	return tss.checkEqualStep()
}

func (tss TimeSeriesSlice) TransformSlice(transform TranformSlice) *TimeSeries {

	step, ok := tss.checkEqualStep()
	if !ok {
		return nil
	}
	start, end := tss.getStartEnd()

	key := "Sum("
	for i, _ := range tss {
		key += tss[i].key
		if i < len(tss)-1 {
			key += ","
		}
	}
	key += ")"

	size := end.Sub(start) / step
	result := &TimeSeries{
		key:   transform.Name() + "(" + key + ")",
		start: start,
		step:  step,
		data:  make([]float64, size),
	}
	for i, _ := range result.data {
		result.data[i] = math.NaN()
	}

	// TODO: this can be optimized to use a different iteration type.
	// since now we use random access and indexing must be determined,
	// it may be a costly iteration method.
	// A better approach may be to use multiple cursors that are
	// eventually aligned.
	cursor := start
	slice := make([]float64, 0, len(tss))
	for i, _ := range result.data {
		slice = slice[:0]
		for j, _ := range tss {
			// slow part
			if v, ok := tss[j].GetAt(cursor); ok {
				slice = append(slice, v)
			}
		}
		if len(slice) > 0 {
			result.data[i] = transform.TransformSlice(slice)
		}
		cursor = cursor.Add(step)
	}
	return result
}

func (tss TimeSeriesSlice) TransformEach(transform Transform) (TimeSeriesSlice, error) {
	return nil, nil
}

func (tss TimeSeriesSlice) C3Data() template.JS {
	step, ok := tss.checkEqualStep()
	if !ok {
		panic("step sizes don't match")
	}
	start, end := tss.getStartEnd()

	s := bytes.NewBufferString("")
	for i, _ := range tss {
		if _, err := s.WriteString("['"); err != nil {
			panic(err)
		}
		if _, err := s.WriteString(tss[i].key); err != nil {
			panic(err)
		}
		if _, err := s.WriteString("', "); err != nil {
			panic(err)
		}
		for cursor := start; cursor.Before(end); cursor = cursor.Add(step) {
			if v, ok := tss[i].GetAt(cursor); !ok {
				if _, err := s.WriteString("NaN"); err != nil {
					panic(err)
				}
			} else {
				if _, err := s.WriteString(strconv.FormatFloat(v, 'E', 2, 64)); err != nil {
					panic(err)
				}
			}
			if err := s.WriteByte(','); err != nil {
				panic(err)
			}
		}
		if s.Len() > 0 {
			s.Truncate(s.Len() - 1)
		}
		if _, err := s.WriteString("]"); err != nil {
			panic(err)
		}
		if i < len(tss)-1 {
			if err := s.WriteByte(','); err != nil {
				panic(err)
			}
		}
	}
	return template.JS(s.String())
}

func (tss TimeSeriesSlice) C3Time() template.JS {
	step, ok := tss.checkEqualStep()
	if !ok {
		panic("step sizes don't match")
	}
	start, end := tss.getStartEnd()

	s := bytes.NewBufferString("[ 'x', ")
	for cursor := start; cursor.Before(end); cursor = cursor.Add(step) {
		if err := s.WriteByte('\''); err != nil {
			panic(err)
		}
		if _, err := s.WriteString(cursor.Format("20060102 15:04:05")); err != nil {
			panic(err)
		}
		if _, err := s.WriteString("',"); err != nil {
			panic(err)
		}
	}
	if s.Len() > 0 {
		s.Truncate(s.Len() - 1)
	}
	if err := s.WriteByte(']'); err != nil {
		panic(err)
	}
	return template.JS(s.String())
}

func (tss TimeSeriesSlice) JSTag() template.JS {
	return template.JS(tss.Key())
}

func (tss TimeSeriesSlice) Key() string {
	s := ""
	for i, _ := range tss {
		s += tss[i].key
		if i < len(tss)-1 {
			s += ","
		}
	}
	return s
}

type TranformSlice interface {
	Name() string
	TransformSlice([]float64) float64
}

type Filter interface {
	Filter(float64) bool
}
