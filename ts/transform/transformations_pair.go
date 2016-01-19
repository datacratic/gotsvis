// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

import "math"

type DividePair struct{}

func (d *DividePair) Name() string {
	return "DividePair"
}

func (d *DividePair) TransformPair(f float64, s float64) float64 {
	if s != 0 {
		return f / s
	}
	return math.NaN()
}

type IsLarger struct{}

func (il *IsLarger) Name() string {
	return "IsOver"
}

func (il *IsLarger) TransformPair(f float64, s float64) float64 {
	if f > s {
		return 1
	}
	return 0
}
