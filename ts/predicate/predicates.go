package predicate

import "math"

var GT = Greater
var EQ = Equal
var NE = NotEqual
var GE = GreaterOrEqual
var LT = Lesser
var LE = LesserOrEqual
var EQNAN = math.IsNaN
var NENAN = NotEqualNaN

var Greater = func(comp float64) func(float64) bool {
	return func(val float64) bool {
		return val > comp
	}
}

var Equal = func(comp float64) func(float64) bool {
	return func(val float64) bool {
		return val == comp
	}
}

var NotEqual = func(comp float64) func(float64) bool {
	return func(val float64) bool {
		return val != comp
	}
}

var EqualNaN = func(val float64) bool {
	return math.IsNaN(val)
}

var NotEqualNaN = func(val float64) bool {
	return !math.IsNaN(val)
}

var GreaterOrEqual = func(comp float64) func(float64) bool {
	return func(val float64) bool {
		return val >= comp
	}
}

var Lesser = func(comp float64) func(float64) bool {
	return func(val float64) bool {
		return val < comp
	}
}

var LesserOrEqual = func(comp float64) func(float64) bool {
	return func(val float64) bool {
		return val <= comp
	}
}
