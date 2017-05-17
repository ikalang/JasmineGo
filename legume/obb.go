package legume

import (
	"math"
	"traffic/util"
	"encoding/json"
)

type AABB struct {
	Center                   util.FloatPoint
	XHalfLength, YHalfLength float64
}

func (a AABB) IsOverlap(b AABB) bool {
	return a.Center.X+a.XHalfLength > b.Center.X-b.XHalfLength &&
		a.Center.Y+a.YHalfLength > b.Center.Y-b.YHalfLength &&
		b.Center.X+b.XHalfLength > a.Center.X-a.XHalfLength &&
		b.Center.Y+b.YHalfLength > a.Center.Y-a.YHalfLength
}

//func (a AABB) Vertexes() [4]util.FloatPoint {
//
//}

type Vector struct {
	X, Y float64
}

func (v Vector) Scale(f float64) Vector {
	return Vector{v.X * f, v.Y * f}
}

func (v Vector) Projection(axis Vector) float64 {
	return math.Abs(v.X*axis.X + v.Y*axis.Y)
}

type OBB struct {
	AABB
	deg          util.Degree
	xAxis, yAxis Vector
}

func (o OBB) String() string {
	data, err := json.Marshal(o)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func CreateOBB(c util.FloatPoint, xh, yh float64, deg util.Degree) (o OBB) {
	o.Center = c
	o.XHalfLength = xh
	o.YHalfLength = yh
	o.deg = deg
	o.xAxis = Vector{math.Cos(float64(deg)), math.Sin(float64(deg))}
	o.yAxis = Vector{-o.xAxis.Y, o.xAxis.X}
	return o
}

func (o OBB) Vertexes() [4]util.FloatPoint {
	o.xAxis.Scale(o.XHalfLength)
	o.yAxis.Scale(o.YHalfLength)
	return [4]util.FloatPoint{
		{o.Center.X + o.xAxis.X + o.yAxis.X, o.Center.Y + o.xAxis.Y + o.yAxis.Y},
		{o.Center.X + o.xAxis.X - o.yAxis.X, o.Center.Y + o.xAxis.Y - o.yAxis.Y},
		{o.Center.X - o.xAxis.X + o.yAxis.X, o.Center.Y - o.xAxis.Y + o.yAxis.Y},
		{o.Center.X - o.xAxis.X - o.yAxis.X, o.Center.Y - o.xAxis.Y - o.yAxis.Y},
	}
}

func (o OBB) Centroid(a, b, c, d float64) util.FloatPoint {
	return util.FloatPoint{
		o.Center.X + 0.5*((a-b)*o.xAxis.X+(c-d)*o.yAxis.X),
		o.Center.Y + 0.5*((a-b)*o.xAxis.Y+(c-d)*o.yAxis.Y),
	}
}

func (o OBB) IsOverlap(b OBB) bool {
	//Separate Axis Testing
	vectorOfCenters := Vector{o.Center.X - b.Center.X, o.Center.Y - b.Center.Y}

	if !(b.xAxis.Projection(o.xAxis)*b.XHalfLength+b.yAxis.Projection(o.xAxis)*b.YHalfLength+
		o.xAxis.Projection(o.xAxis)*o.XHalfLength > vectorOfCenters.Projection(o.xAxis)) {
		return false
	}

	if !(b.xAxis.Projection(o.yAxis)*b.XHalfLength+b.yAxis.Projection(o.yAxis)*b.YHalfLength+
		o.yAxis.Projection(o.yAxis)*o.YHalfLength > vectorOfCenters.Projection(o.yAxis)) {
		return false
	}

	if !(o.xAxis.Projection(b.xAxis)*o.XHalfLength+o.yAxis.Projection(b.xAxis)*o.YHalfLength+
		b.xAxis.Projection(b.xAxis)*b.YHalfLength > vectorOfCenters.Projection(b.xAxis)) {
		return false
	}

	if !(o.xAxis.Projection(b.yAxis)*o.XHalfLength+o.yAxis.Projection(b.yAxis)*o.YHalfLength+
		b.yAxis.Projection(b.yAxis)*b.YHalfLength > vectorOfCenters.Projection(b.yAxis)) {
		return false
	}

	return true
}

func (o OBB) SymmetryXAxis() (r OBB) {
	r.Center = o.Center.SymmetryXAxis()
	r.deg = o.deg.SymmetryXAxis()
	return r
}

func (o OBB) Transform(d1, d2 util.Direction) (r OBB) {
	r.Center = o.Center.SymmetryXAxis()
	r.deg = o.deg.SymmetryXAxis()
	return r
}
