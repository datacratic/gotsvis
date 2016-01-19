// Copyright (c) 2014 Datacratic. All rights reserved.

package ts

import (
	"errors"
	"fmt"
	"math"
)

type TimeSeriesPair struct {
	First  *TimeSeries
	Second *TimeSeries
}

func (tsp *TimeSeriesPair) TransformPair(t TransformPair) (*TimeSeries, error) {

	if !tsp.First.IsEqualStep(tsp.Second) {
		return nil, errors.New("step sizes don't match")
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

	key := fmt.Sprintf("%s(%s,%s)", t.Name(), tsp.First.key, tsp.Second.key)
	result, err := NewTimeSeriesOfTimeRange(key, start, end, step, math.NaN())
	if err != nil {
		return nil, err
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
	return result, nil
}

type TransformPair interface {
	Name() string
	TransformPair(float64, float64) float64
}
