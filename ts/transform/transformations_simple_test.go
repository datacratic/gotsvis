// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/datacratic/gotsvis/ts"
)

type TestSeries struct {
	Key   string
	Start time.Time
	End   time.Time
	Step  time.Duration
	Data  []float64
}

func (ts *TestSeries) String() string {
	s := bytes.NewBufferString("")
	s.WriteString(ts.Key)

	s.WriteString(" Start: ")
	s.WriteString(ts.Start.String())

	s.WriteString(" End: ")
	s.WriteString(ts.End.String())

	s.WriteString(" Step: ")
	s.WriteString(ts.Step.String())

	s.WriteString(" Length: ")
	s.WriteString(strconv.Itoa(len(ts.Data)))

	s.WriteString(" ")

	for _, v := range ts.Data {
		s.WriteString(strconv.FormatFloat(v, 'f', 2, 64))
		s.WriteByte(',')
	}
	if s.Len() > 0 {
		s.Truncate(s.Len() - 1)
	}
	return s.String()
}

func checkTimeSeries(t *testing.T, got *ts.TimeSeries, exp *TestSeries) {
	if got == nil {
		t.Errorf("FAIL(TimeSeries): can't be nil")
		return
	}
	if exp == nil {
		t.Errorf("FAIL(TimeSeries): can't be nil")
		return
	}
	checkKey(t, got.Key(), exp.Key)
	checkStart(t, got.Start(), exp.Start)
	checkEnd(t, got.End(), exp.End)
	checkStep(t, got.Step(), exp.Step)
	checkData(t, got.Data(), exp.Data)
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

func TestTransforms(t *testing.T) {
	start := time.Date(2016, time.Month(1), 14, 10, 0, 0, 0, time.UTC)
	end := start.Add(6 * time.Minute)
	step := time.Minute
	NaN := math.NaN()

	ts1, err := ts.NewTimeSeriesOfTimeRange("ts1", start, end, step, 1)
	checkErr(t, err)
	if ts1 == nil {
		t.Errorf("FAIL(ts1): can't be nil, if we want to continue with other tests")
		return
	}

	ts1NaN, err := ts.NewTimeSeriesOfData("ts1NaN", start, step,
		[]float64{1, NaN, NaN, 1, 1, NaN})
	checkErr(t, err)
	if ts1NaN == nil {
		t.Errorf("FAIL(ts1NaN): can't be nil, if we want to continue with other tests")
		return
	}

	tsRaise, err := ts.NewTimeSeriesOfData("tsRaise", start, step,
		[]float64{NaN, 1, 1, 0, 1, NaN, 2, 2})
	checkErr(t, err)
	if tsRaise == nil {
		t.Errorf("FAIL(tsRaise): can't be nil, if we want to continue with other tests")
		return
	}

	tss := []struct {
		Got *ts.TimeSeries
		Exp *TestSeries
	}{
		{
			Got: ts1.Transform(&CumulativeSum{}),
			Exp: &TestSeries{
				Key:   "CumulativeSum(ts1)",
				Start: start,
				End:   end,
				Step:  step,
				Data:  []float64{1, 2, 3, 4, 5, 6},
			},
		},
		{
			Got: ts1NaN.Transform(&CumulativeSum{}),
			Exp: &TestSeries{
				Key:   "CumulativeSum(ts1NaN)",
				Start: start,
				End:   end,
				Step:  step,
				Data:  []float64{1, 1, 1, 2, 3, 3},
			},
		},
		{
			Got: ts1.Transform(&DivideBy{2}),
			Exp: &TestSeries{
				Key:   (&DivideBy{2}).Name() + "(ts1)",
				Start: start,
				End:   end,
				Step:  step,
				Data:  []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
			},
		},
		{
			Got: ts1NaN.Transform(&DivideBy{2}),
			Exp: &TestSeries{
				Key:   (&DivideBy{2}).Name() + "(ts1NaN)",
				Start: start,
				End:   end,
				Step:  step,
				Data:  []float64{0.5, NaN, NaN, 0.5, 0.5, NaN},
			},
		},
		{
			Got: ts1.Transform(&MultiplyBy{2}),
			Exp: &TestSeries{
				Key:   (&MultiplyBy{2}).Name() + "(ts1)",
				Start: start,
				End:   end,
				Step:  step,
				Data:  []float64{2, 2, 2, 2, 2, 2},
			},
		},
		{
			Got: ts1NaN.Transform(&MultiplyBy{2}),
			Exp: &TestSeries{
				Key:   (&MultiplyBy{2}).Name() + "(ts1NaN)",
				Start: start,
				End:   end,
				Step:  step,
				Data:  []float64{2, NaN, NaN, 2, 2, NaN},
			},
		},
		{
			Got: tsRaise.Transform(&MarkRaise{}),
			Exp: &TestSeries{
				Key:   (&MarkRaise{}).Name() + "(tsRaise)",
				Start: start,
				End:   start.Add(8 * step),
				Step:  step,
				Data:  []float64{NaN, 0, 0, 0, 1, NaN, 1, 0},
			},
		},
	}

	for _, pair := range tss {
		fmt.Printf("%s\n%s\n\n", pair.Got, pair.Exp)
		checkTimeSeries(t, pair.Got, pair.Exp)
	}
	//t.Errorf("here")
}
