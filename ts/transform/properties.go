// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

import "github.com/datacratic/gotsvis/ts"

func Any(ts *ts.TimeSeries, predicate func(float64) bool) bool {
	it := ts.Iterator()
	for val, ok := it.Next(); ok; val, ok = it.Next() {
		if predicate(val) {
			return true
		}
	}
	return false
}

func None(ts *ts.TimeSeries, predicate func(float64) bool) bool {
	it := ts.Iterator()
	for val, ok := it.Next(); ok; val, ok = it.Next() {
		if predicate(val) {
			return false
		}
	}
	return true
}
