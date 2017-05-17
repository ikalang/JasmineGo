package util

import (
	"math"
)

func FloatEqual(a, b float64) bool {
	if math.Abs(a) < 0.008 && math.Abs(b) < 0.008 {
		return true
	}

	return math.Abs((a-b)/math.Max(a, b)) < 0.005
}

func FloatEqualTolerance(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}

func FloatInOpenInterval(a, b, c, tolerance float64) bool {
	return a > math.Min(b, c)-tolerance && a < math.Max(b, c)+tolerance
}

func FloatInCloseInterval(a, b, c, tolerance float64) bool {
	return !(a < math.Min(b, c)-tolerance || a > math.Max(b, c)+tolerance)
}
