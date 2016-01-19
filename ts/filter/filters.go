// Copyright (c) 2014 Datacratic. All rights reserved.

package filter

import "math"

type HasNotNaN struct{}

func (hnn HasNotNaN) Filter(val float64) bool {
	if math.IsNaN(val) {
		return false
	}
	return true
}

type HasLarger struct {
	Value float64
}

func (hl *HasLarger) Filter(val float64) bool {
	return hl.Value < val
}

type HasLargerOrEqual struct {
	Value float64
}

func (hloe *HasLargerOrEqual) Filter(val float64) bool {
	return hloe.Value <= val
}
