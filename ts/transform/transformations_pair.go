// Copyright (c) 2014 Datacratic. All rights reserved.

package transform

import (
	"fmt"
	"math"
)

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

type And struct {
	FirstPredicate      func(float64) bool
	FirstPredicateName  string
	SecondPredicate     func(float64) bool
	SecondPredicateName string
}

func (and *And) TransformPair(f float64, s float64) float64 {
	if and.FirstPredicate(f) && and.SecondPredicate(s) {
		return 1
	}
	return 0
}

func (and *And) Name() string {
	return fmt.Sprintf("(%s)And(%s)", and.FirstPredicateName, and.SecondPredicateName)
}

type Or struct {
	FirstPredicate      func(float64) bool
	FirstPredicateName  string
	SecondPredicate     func(float64) bool
	SecondPredicateName string
}

func (or *Or) TransformPair(f float64, s float64) float64 {
	if or.FirstPredicate(f) || or.SecondPredicate(s) {
		return 1
	}
	return 0
}

func (or *Or) Name() string {
	return fmt.Sprintf("(%s)Or(%s)", or.FirstPredicateName, or.SecondPredicateName)
}
