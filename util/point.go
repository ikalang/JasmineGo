package util

import "math"

type FloatPoint struct {
	X, Y float64
}

var FloatPointZero = FloatPoint{0.0, 0.0}
var InvalidFloatPoint = FloatPoint{math.NaN(), math.NaN()}

func (f FloatPoint) Equal(q FloatPoint) bool {
	return FloatEqual(f.X, q.X) && FloatEqual(f.Y, q.Y)
}

func (f FloatPoint) CenterPoint(q FloatPoint) FloatPoint {
	return FloatPoint{(f.X + q.X) * 0.5, (f.Y + q.Y) * 0.5}
}

func (f FloatPoint) Distance(q FloatPoint) float64 {
	return math.Sqrt(math.Pow(f.X-q.X, 2) + math.Pow(f.Y-q.Y, 2))
}

func (f FloatPoint) DegreeTo(q FloatPoint) (d Degree) {
	if FloatEqual(f.X, q.X) {
		if f.Y > q.Y {
			d = 90
		} else {
			d = 270
		}
	} else {
		k := Slope((f.Y - q.Y) / (f.X - q.X))
		d = k.ToDegree()

		if f.X < q.X {
			d += 180
		}
	}

	return d
}

func (f FloatPoint) ToIntPoint() IntPoint {
	return IntPoint{int(f.X), int(f.Y)}
}

type IntPoint struct {
	X, Y int
}

func (i IntPoint) Equal(q IntPoint) bool {
	return i.X == q.X && i.Y == q.Y
}

func (f FloatPoint) Scale(scale float64) FloatPoint {
	return FloatPoint{f.X * scale, f.Y * scale}
}

func (f FloatPoint) Shift(shift FloatPoint) FloatPoint {
	return FloatPoint{f.X + shift.X, f.Y + shift.Y}
}

func (i IntPoint) ToFloatPoint() FloatPoint {
	return FloatPoint{float64(i.X), float64(i.Y)}
}

func (f FloatPoint) SymmetryXAxis() FloatPoint {
	return FloatPoint{f.X, -f.Y}
}

func (f FloatPoint) SymmetryYAxis() FloatPoint {
	return FloatPoint{-f.X, f.Y}
}

func (f FloatPoint) SymmetryOrigin() FloatPoint {
	return FloatPoint{-f.X, -f.Y}
}
