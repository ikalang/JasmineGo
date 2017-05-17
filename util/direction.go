package util

import "math"

type Direction int
type Degree float64 //   [0-360)
type Slope float64
type Rad float64

const (
	DirErr = iota
	XInc
	XDec
	YInc
	YDec
)

func (d Direction) Opposite() Direction {
	switch d {
	case XInc:
		return XDec
	case XDec:
		return XInc
	case YInc:
		return YDec
	case YDec:
		return YInc
	default:
		return DirErr
	}
}

func (d Direction) String() string {
	switch d {
	case XInc:
		return "XInc"
	case XDec:
		return "XDec"
	case YInc:
		return "YInc"
	case YDec:
		return "YDec"
	default:
		return "Direction Error"
	}
}

func (d Direction) ToDegree() Degree {
	switch d {
	case XInc:
		return 0
	case XDec:
		return 180
	case YInc:
		return 90
	case YDec:
		return 270
	default:
		return -1
	}
}

func (d Degree) Opposite() Degree {
	switch d {
	case 0:
		return 180
	case 180:
		return 0
	case 90:
		return 270
	case 270:
		return 90
	default:
		return -1
	}
}

func (d Degree) ToRad() Rad {
	return Rad(d * math.Pi / 180)
}

func (d Degree) ToSlope() Slope {
	return Slope(math.Tan(float64(d.ToRad())))
}

func (d Degree) Equal(q Degree) bool {
	return FloatEqual(float64(d), float64(q))
}

func (d Degree) SymmetryXAxis() Degree {
	return 360 - d
}

func (s Slope) ToDegree() Degree {
	return Degree(math.Mod(180*math.Atan(float64(s))/math.Pi, 360))
}
