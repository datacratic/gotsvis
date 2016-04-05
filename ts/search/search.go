// Copyright (c) 2014 Datacratic. All rights reserved.

package search

import (
	"math"
	"time"

	"github.com/datacratic/gotsvis/ts"
)

var NaN = math.NaN()

func First(ts *ts.TimeSeries, predicate func(float64) bool) (time.Time, float64, bool) {
	if ts == nil {
		return time.Time{}, NaN, false
	}
	it := ts.IteratorTimeValue()
	for t, v, ok := it.Next(); ok; t, v, ok = it.Next() {
		if predicate(v) {
			return t, v, true
		}
	}
	return time.Time{}, 0, false
}
