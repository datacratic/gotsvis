// Copyright (c) 2014 Datacratic. All rights reserved.

package ts

import (
	"fmt"
	"math"
	"testing"
	"time"
)

var NaN = math.NaN()

func checkTimeSeries(t *testing.T, got, exp *TimeSeries) {
	if got == nil {
		t.Errorf("FAIL(TimeSeries): can't be nil")
		return
	}
	if exp == nil {
		t.Errorf("FAIL(TimeSeries): can't be nil")
		return
	}
	checkKey(t, got.key, exp.key)
	checkStart(t, got.start, exp.start)
	checkEnd(t, got.End(), exp.End())
	checkStep(t, got.step, exp.step)
	checkData(t, got.data, exp.data)
}

func checkData(t *testing.T, got, exp []float64) {
	if len(got) != len(exp) {
		t.Errorf("FAIL(data): length '%d' != '%d':\ngot:\n\t%v,\nexpected:\n\t%v",
			len(got), len(exp), got, exp)
		return
	}

	for i, g := range got {
		if g != exp[i] && (!math.IsNaN(g) || !math.IsNaN(exp[i])) {
			t.Errorf("FAIL(data): at index: '%d', '%f' != '%f':\ngot:\n\t%v,\nexpected:\n\t%v",
				i, g, exp[i], got, exp)
		}
	}
}

func checkKey(t *testing.T, got, exp string) {
	if got != exp {
		t.Errorf("FAIL(key): got: '%s', expected '%s'", got, exp)
	}
}

func checkStart(t *testing.T, got, exp time.Time) {
	if got != exp {
		t.Errorf("FAIL(start): got: '%s', expected '%s'", got, exp)
	}
}

func checkEnd(t *testing.T, got, exp time.Time) {
	if got != exp {
		t.Errorf("FAIL(end): got: '%s', expected '%s'", got, exp)
	}
}

func checkStep(t *testing.T, got, exp time.Duration) {
	if got != exp {
		t.Errorf("FAIL(duration): got: '%s', expected '%s'", got, exp)
	}
}

func checkLengthDataEqual(t *testing.T, data []float64, exp int) {
	if len(data) != exp {
		t.Errorf("FAIL(data length): got: '%d', expected '%d'", len(data), exp)
	}
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("FAIL(error): %s", err)
	}
}

func TestTimeSeriesConstructors(t *testing.T) {
	start := time.Date(2016, time.Month(1), 13, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Minute)
	step := time.Minute

	ts0, err := NewTimeSeriesOfTimeRange("test0", start, start, step, 0)
	checkErr(t, err)

	ts1, err := NewTimeSeriesOfTimeRange("test1", start, end, step, 0)
	checkErr(t, err)

	ts2, err := NewTimeSeriesOfLength("test2", start, step, 0, 0)
	checkErr(t, err)

	ts3, err := NewTimeSeriesOfData("test3", start, step, []float64{0})
	checkErr(t, err)

	ts4, err := NewTimeSeriesOfLength("test4", start, step, 1, 0)
	checkErr(t, err)

	data5 := []float64{1, math.NaN(), 0, 2, math.NaN()}
	ts5, _ := NewTimeSeriesOfData("test5", start, step, data5)

	ts6, err := NewTimeSeriesOfData("test6", start, step, []float64{})
	checkErr(t, err)

	tss := []struct {
		Got *TimeSeries
		Exp *TimeSeries
	}{
		{
			Got: ts0,
			Exp: &TimeSeries{
				key:   "test0",
				start: start,
				step:  step,
				data:  []float64{},
			},
		},
		{
			Got: ts1,
			Exp: &TimeSeries{
				key:   "test1",
				start: start,
				step:  step,
				data:  []float64{0},
			},
		},
		{
			Got: ts2,
			Exp: &TimeSeries{
				key:   "test2",
				start: start,
				step:  step,
				data:  []float64{},
			},
		},
		{
			Got: ts3,
			Exp: &TimeSeries{
				key:   "test3",
				start: start,
				step:  step,
				data:  []float64{0},
			},
		},
		{
			Got: ts4,
			Exp: &TimeSeries{
				key:   "test4",
				start: start,
				step:  step,
				data:  []float64{0},
			},
		},
		{
			Got: ts5,
			Exp: &TimeSeries{
				key:   "test5",
				start: start,
				step:  step,
				data:  data5,
			},
		},
		{
			Got: ts6,
			Exp: &TimeSeries{
				key:   "test6",
				start: start,
				step:  step,
				data:  []float64{},
			},
		},
	}

	for _, pair := range tss {
		fmt.Printf("%s\n%s\n\n", pair.Got, pair.Exp)
		checkTimeSeries(t, pair.Got, pair.Exp)
	}
}

func TestTimeSeriesExtend(t *testing.T) {
	start := time.Date(2016, time.Month(1), 25, 10, 0, 0, 0, time.UTC)
	step := time.Minute

	ts0, err := NewTimeSeriesOfData("test0", start, step, []float64{1, 2, 3})
	checkErr(t, err)
	ts1 := ts0.Copy()
	ts1.ExtendBy(3 * time.Minute)

	ts2 := ts0.Copy()
	ts2.ExtendTo(time.Date(2016, time.Month(1), 25, 10, 5, 0, 0, time.UTC))

	ts3 := ts0.Copy()
	ts3.ExtendWith([]float64{4, 5, 6})

	tss := []struct {
		Got *TimeSeries
		Exp *TimeSeries
	}{
		{
			Got: ts1,
			Exp: &TimeSeries{
				key:   "test0",
				start: start,
				step:  step,
				data:  []float64{1, 2, 3, NaN, NaN, NaN},
			},
		},
		{
			Got: ts2,
			Exp: &TimeSeries{
				key:   "test0",
				start: start,
				step:  step,
				data:  []float64{1, 2, 3, NaN, NaN, NaN},
			},
		},
		{
			Got: ts3,
			Exp: &TimeSeries{
				key:   "test0",
				start: start,
				step:  step,
				data:  []float64{1, 2, 3, 4, 5, 6},
			},
		},
	}

	for _, pair := range tss {
		fmt.Printf("%s\n%s\n\n", pair.Got, pair.Exp)
		checkTimeSeries(t, pair.Got, pair.Exp)
	}
}
