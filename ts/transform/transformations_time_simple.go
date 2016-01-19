// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

import (
	"math"
	"time"
)

type Summarize struct {
	NewStep time.Duration

	cursor time.Time
	sum    float64
}

func (sum *Summarize) Name() string {
	return "Summarize(" + sum.NewStep.String() + ")"
}
func (sum *Summarize) Start() time.Time {
	return time.Time{}
}
func (sum *Summarize) End() time.Time {
	return time.Time{}
}
func (sum *Summarize) Step() time.Duration {
	return sum.NewStep
}

func (sum *Summarize) TimeTransform(t time.Time, f float64) (rt time.Time, rv float64) {
	trunc := t.Truncate(sum.NewStep)
	if sum.cursor.IsZero() {
		sum.cursor = trunc
	}

	if trunc.Equal(sum.cursor) {
		if !math.IsNaN(f) {
			sum.sum += f
		}
		rt = sum.cursor
		rv = sum.sum
	} else {
		rt = sum.cursor
		rv = sum.sum

		sum.cursor = sum.cursor.Add(sum.NewStep)
		if !math.IsNaN(f) {
			sum.sum = f
		} else {
			sum.sum = 0
		}
	}
	return
}
