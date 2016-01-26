// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

import (
	"fmt"
	"math"
	"time"

	. "github.com/datacratic/gotsvis/ts"
)

func Summarize(ts *TimeSeries, step time.Duration) *TimeSeries {

	key := fmt.Sprintf("Summarize(%v)(%s)", step, ts.Key())
	start := ts.Start().Truncate(step)
	end := ts.End().Truncate(step)
	if !ts.End().Equal(end) {
		end = end.Add(step)
	}

	newTs, err := NewTimeSeriesOfTimeRange(key, start, end, step, math.NaN())
	if err != nil {
		return nil
	}

	var sum float64
	newC := start
	for oldC, oldE, oldS := ts.Start(), ts.End(), ts.Step(); oldC.Before(oldE); oldC = oldC.Add(oldS) {
		if newC.Equal(oldC.Truncate(step)) {
			v, ok := ts.GetAt(oldC)
			if !ok || math.IsNaN(v) {
				continue
			}
			sum += v
		} else {
			newTs.SetAt(newC, sum)
			newC = newC.Add(step)

			v, ok := ts.GetAt(oldC)
			if !ok || math.IsNaN(v) {
				sum = 0
				continue
			}
			sum = v
		}
	}
	newTs.SetAt(newC, sum)

	return newTs
}
