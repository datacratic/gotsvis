// Copyright (c) 2014 Datacratic. All rights reserved.

package ts

import (
	"fmt"
	"math"
)

type TimeSeriesPair struct {
	First  *TimeSeries
	Second *TimeSeries
}

func (tsp *TimeSeriesPair) TransformPair(t TransformPair) *TimeSeries {

	if !tsp.First.IsEqualStep(tsp.Second) {
		return nil
	}
	step := tsp.First.step

	start := tsp.First.start
	if start.After(tsp.Second.start) {
		start = tsp.Second.start
	}
	end := tsp.First.End()
	if end.Before(tsp.Second.End()) {
		end = tsp.Second.End()
	}

	size := end.Sub(start) / step
	result := &TimeSeries{
		key:   fmt.Sprintf("%s(%s,%s)", t.Name(), tsp.First.key, tsp.Second.key),
		start: start,
		step:  step,
		data:  make([]float64, size),
	}
	for i, _ := range result.data {
		result.data[i] = math.NaN()
	}

	cursor := start
	for i, _ := range result.data {
		f, ok1 := tsp.First.GetAt(cursor)
		s, ok2 := tsp.Second.GetAt(cursor)
		cursor = cursor.Add(step)
		if !ok1 && !ok2 {
			continue
		}
		result.data[i] = t.TransformPair(f, s)
	}
	return result
}

type TransformPair interface {
	Name() string
	TransformPair(float64, float64) float64
}
