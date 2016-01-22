// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/datacratic/gotsvis/ts"
)

func TestSummarize(t *testing.T) {
	start := time.Date(2016, time.Month(1), 21, 0, 0, 0, 0, time.UTC)
	end := start.Add(36 * time.Hour)
	step := time.Hour
	NaN := math.NaN()

	ts1, err := ts.NewTimeSeriesOfTimeRange("ts1", start, end, step, 1)
	checkErr(t, err)
	if ts1 == nil {
		t.Errorf("FAIL(ts1): can't be nil, if we want to continue with other tests")
		return
	}

	ts1NaN, err := ts.NewTimeSeriesOfData("ts1NaN", start, step,
		[]float64{
			1, NaN, 1, NaN, 1, NaN, 1, NaN, 1, NaN, 1, NaN,
			1, NaN, 1, NaN, 1, NaN, 1, NaN, 1, NaN, 1, NaN,
			1, NaN, 1, NaN, 1, NaN, 1, NaN, 1, NaN, 1, NaN,
		})
	checkErr(t, err)
	if ts1NaN == nil {
		t.Errorf("FAIL(ts1NaN): can't be nil, if we want to continue with other tests")
		return
	}

	tss := []struct {
		Got *ts.TimeSeries
		Exp *TestSeries
	}{
		{
			Got: Summarize(ts1, 24*time.Hour),
			Exp: &TestSeries{
				Key:   "Summarize(24h0m0s)(ts1)",
				Start: start,
				End:   start.Add(2 * 24 * time.Hour),
				Step:  24 * time.Hour,
				Data:  []float64{24, 12},
			},
		},
		{
			Got: Summarize(ts1NaN, 24*time.Hour),
			Exp: &TestSeries{
				Key:   "Summarize(24h0m0s)(ts1NaN)",
				Start: start,
				End:   start.Add(2 * 24 * time.Hour),
				Step:  24 * time.Hour,
				Data:  []float64{12, 6},
			},
		},
	}

	for _, pair := range tss {
		fmt.Printf("%s\n%s\n\n", pair.Got, pair.Exp)
		checkTimeSeries(t, pair.Got, pair.Exp)
	}
	//t.Errorf("here")
}
