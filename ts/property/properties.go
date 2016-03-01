// Copyright (c) 2014 Datacratic. All rights reserved.

package property

import "github.com/datacratic/gotsvis/ts"

func Any(ts *ts.TimeSeries, predicate func(float64) bool) bool {
	if ts == nil {
		return false
	}
	it := ts.Iterator()
	for val, ok := it.Next(); ok; val, ok = it.Next() {
		if predicate(val) {
			return true
		}
	}
	return false
}

func None(ts *ts.TimeSeries, predicate func(float64) bool) bool {
	if ts == nil {
		return false
	}
	it := ts.Iterator()
	for val, ok := it.Next(); ok; val, ok = it.Next() {
		if predicate(val) {
			return false
		}
	}
	return true
}

func Last(ts *ts.TimeSeries, predicate func(float64) bool) bool {
	if ts == nil {
		return false
	}
	last, ok := ts.Iterator().Last()
	if !ok {
		return false
	}
	return predicate(last)
}
