// Copyright (c) 2014 Datacratic. All rights reserved.

package ts

func ones(ts *TimeSeries) {
	for i, _ := range ts.data {
		ts.data[i] = 1.0
	}
}

/*
func TestDTSSAligned(t *testing.T) {
	start := time.Date(2015, time.Month(12), 29, 15, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	ts1 := DenseTimeSeries{
		Key:   "ts1",
		Start: start,
		End:   end,
		Step:  time.Minute,
	}
	ts1.Init()
	ones(&ts1)
	ts2 := ts1
	tss := DTSS{ts1, ts2}

	sum := Sum{}
	sumSeries, err := tss.TransformSlice(&sum)
	if err != nil {
		t.Fatal(err)
	}

	if !start.Equal(sumSeries.Start) {
		t.Fatal("start not good")
	}
	if !end.Equal(sumSeries.End) {
		t.Fatal("end not good")
	}
	if time.Minute != sumSeries.Step {
		t.Fatal("step not good")
	}

	for _, v := range sumSeries.Data {
		if *v != 2.0 {
			t.Fatal("should all be 2.0")
		}
	}
}

func TestDTSSNoOverlap(t *testing.T) {
	start := time.Date(2015, time.Month(12), 29, 15, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	ts1 := DenseTimeSeries{
		Key:   "ts1",
		Start: start,
		End:   end,
		Step:  time.Minute,
	}
	ts1.Init()
	ones(&ts1)

	ts2 := ts1
	ts2.Start = ts2.Start.Add(time.Hour)
	ts2.End = ts2.End.Add(time.Hour)

	tss := DTSS{ts1, ts2}

	sum := Sum{}
	sumSeries, err := tss.TransformSlice(&sum)
	if err != nil {
		t.Fatal(err)
	}

	if time.Minute != sumSeries.Step {
		t.Fatal("step not good")
	}
	if !ts1.Start.Equal(sumSeries.Start) {
		t.Fatal("start not good")
	}
	if !ts2.End.Equal(sumSeries.End) {
		t.Fatal("end not good")
	}

	for _, v := range sumSeries.Data {
		if *v != 1.0 {
			t.Fatal("should all be 1.0")
		}
	}
}

func TestDTSSWithOverlap(t *testing.T) {
	start := time.Date(2015, time.Month(12), 29, 15, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	ts1 := DenseTimeSeries{
		Key:   "ts1",
		Start: start,
		End:   end,
		Step:  time.Minute,
	}
	ts1.Init()
	ones(&ts1)

	ts2 := ts1
	ts2.Start = ts2.Start.Add(30 * time.Minute)
	ts2.End = ts2.End.Add(30 * time.Minute)

	tss := DTSS{ts1, ts2}

	sum := Sum{}
	sumSeries, err := tss.TransformSlice(&sum)
	if err != nil {
		t.Fatal(err)
	}

	if time.Minute != sumSeries.Step {
		t.Fatal("step not good")
	}
	if !ts1.Start.Equal(sumSeries.Start) {
		t.Fatal("start not good")
	}
	if !ts2.End.Equal(sumSeries.End) {
		t.Fatal("end not good")
	}

	cursor := sumSeries.Start
	for _, v := range sumSeries.Data {
		_, ok1 := ts1.GetAt(cursor)
		_, ok2 := ts2.GetAt(cursor)
		if !ok1 || !ok2 {
			if *v != 1.0 {
				t.Fatal("should be 1.0")
			}
		} else {
			if *v != 2.0 {
				t.Fatal("should be 2.0")
			}
		}

		if *v != 1.0 {
			if (cursor.Equal(ts2.Start) || cursor.After(ts2.Start)) &&
				cursor.Before(ts1.End) && *v != 2.0 {
				t.Fatal("should be 2.0")
				continue
			}
		}
		cursor = cursor.Add(sumSeries.Step)
	}
}

func TestC3DTSS(t *testing.T) {
	start := time.Date(2015, time.Month(12), 29, 15, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	ts1 := DenseTimeSeries{
		Key:   "ts1",
		Start: start,
		End:   end,
		Step:  time.Minute,
	}
	ts1.Init()
	ones(&ts1)
	ts2 := ts1
	tss := DTSS{ts1, ts2}
	fmt.Println(tss)
	fmt.Println(tss.C3Time())
	fmt.Println(tss.C3Data())

	pair := &TimeSeriesPair{
		First:  &ts1,
		Second: &ts2,
	}

	ans, _ := pair.TransformPair(&DividePair{})
	fmt.Println(ans)
	fmt.Println(ans.C3Data())
	tss = DTSS{*ans}
	fmt.Println(tss.C3Data())

	div := ans.Transform(&DivideBy{By: 10})
	fmt.Println(div)
	fmt.Println(div.C3Data())

	ctor := NewDenseTimeSeries("target", div.Start, div.End, div.Step, 2.0)
	fmt.Println(ctor)
	fmt.Println(ctor.C3Data())

	hSum := ctor.TimeTransform(&Summarize{
		To: &DenseTimeSeries{
			Start: ctor.Start,
			End:   ctor.End,
			Step:  24 * time.Hour,
		},
	})
	fmt.Println(hSum)
	fmt.Println(hSum.C3Data())
}
*/
