// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

import (
	"fmt"
	"math"
	"time"

	. "github.com/datacratic/gotsvis/ts"
)

func SubSeries(ts *TimeSeries, start, end time.Time) *TimeSeries {
	if ts == nil {
		return nil
	}

	if start.IsZero() || start.Before(ts.Start()) {
		start = ts.Start()
	}
	if end.IsZero() || ts.End().Before(end) {
		end = ts.End()
	}
	key := fmt.Sprintf("SubSeries(%s,%s)(%s)", start.Format(time.RFC3339), end.Format(time.RFC3339), ts.Key())
	step := ts.Step()

	newTs, err := NewTimeSeriesOfTimeRange(key, start, end, step, math.NaN())
	if err != nil {
		return nil
	}

	for cursor := start; cursor.Before(end); cursor = cursor.Add(step) {
		v, _ := ts.GetAt(cursor)
		newTs.SetAt(cursor, v)
		if int(step) == 0 {
			break
		}
	}
	return newTs
}
