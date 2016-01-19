// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

type Sum struct {
}

func (s *Sum) Name() string {
	return "Sum"
}

func (s *Sum) TransformSlice(vals []float64) float64 {
	var sum float64
	for _, v := range vals {
		sum += v
	}
	return sum
}
