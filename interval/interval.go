package interval

import "math"

type Interval struct {
	Min float64
	Max float64
}

func NewInterval() Interval {
	return Interval{
		Min: math.Inf(-1),
		Max: math.Inf(1),
	}
}

func (i Interval) Size() float64 {
	return i.Max - i.Min
}

func (i Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

func (i Interval) Sorrounds(x float64) bool {
	return i.Min < x && x < i.Max
}

var interval_empty Interval = Interval{Min: math.Inf(1), Max: math.Inf(-1)}
var interval_universe Interval = Interval{Min: math.Inf(-1), Max: math.Inf(1)}
