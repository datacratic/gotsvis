// Copyright (c) 2014 Datacratic. All rights reserved.

package property

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/datacratic/gotsvis/ts"
	. "github.com/datacratic/gotsvis/ts/predicate"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("FAIL(error): %s", err)
	}
}

func TestProperties(t *testing.T) {

	start := time.Date(2016, time.Month(1), 21, 0, 0, 0, 0, time.UTC)
	step := time.Hour
	NaN := math.NaN()

	ts1, err := ts.NewTimeSeriesOfData("ts1", start, step, []float64{1, 2, 3, 4, 5})
	checkErr(t, err)

	tsNaN, err := ts.NewTimeSeriesOfData("tsNaN", start, step, []float64{1, 2, NaN, 4, 5})
	checkErr(t, err)

	tests := []struct {
		f    func(*ts.TimeSeries, func(float64) bool) bool
		ts   *ts.TimeSeries
		pred func(float64) bool
		exp  bool
	}{
		{
			f:    Any,
			ts:   ts1,
			pred: EQ(5),
			exp:  true,
		},
		{
			f:    Any,
			ts:   ts1,
			pred: GT(5),
			exp:  false,
		},
		{
			f:    None,
			ts:   ts1,
			pred: GT(5),
			exp:  true,
		},
		{
			f:    None,
			ts:   ts1,
			pred: EQ(5),
			exp:  false,
		},
		{
			f:    Any,
			ts:   tsNaN,
			pred: EQNAN,
			exp:  true,
		},
		{
			f:    None,
			ts:   tsNaN,
			pred: EQNAN,
			exp:  false,
		},
		{
			f:    Any,
			ts:   tsNaN,
			pred: EQ(5),
			exp:  true,
		},
		{
			f:    Any,
			ts:   tsNaN,
			pred: GT(5),
			exp:  false,
		},
		{
			f:    None,
			ts:   tsNaN,
			pred: GT(5),
			exp:  true,
		},
		{
			f:    None,
			ts:   tsNaN,
			pred: EQ(5),
			exp:  false,
		},
		{
			f:    None,
			ts:   nil,
			pred: EQ(5),
			exp:  false,
		},
	}

	for i, test := range tests {
		got := test.f(test.ts, test.pred)
		if got != test.exp {
			fmt.Printf("FAIL(%d): got '%v' != '%v' exp\n", i, got, test.exp)
			t.Fail()
		}
	}

}
