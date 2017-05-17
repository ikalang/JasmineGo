package track

import (
	"log"
	"math"
	"testing"
	"traffic/util"
)

func TestQuadraticCurveTrack_ToTrack(t *testing.T) {
	var qct QuadraticCurveTrack
	if err := qct.LoadFromJSONFile("data/QTurn4to7.json"); err != nil {
		t.Error(err)
		return
	}

	//tk := qct.ToTrack()
	//log.Print(tk)
}

func TestFunction_YInflectionPoint(t *testing.T) {
	type Function_YInflectionPointCase struct {
		Function
		res bool
		pt  util.FloatPoint
	}

	cases := []Function_YInflectionPointCase{
	//{Function: Function{Start: util.FloatPoint{2, math.Sqrt(3)}, End: util.FloatPoint{2, -math.Sqrt(3)}, P: [6]float64{1, -1, 0, 0, 0, -1}}, res: true, pt: util.FloatPoint{1, 0}},
	//{Function: Function{Start: util.FloatPoint{-2, math.Sqrt(3)}, End: util.FloatPoint{2, -math.Sqrt(3)}, P: [6]float64{1, -1, 0, 0, 0, -1}}, res: false, pt: util.FloatPoint{0, 0}},
	//{Function: Function{Start: util.FloatPoint{0, -1}, End: util.FloatPoint{0, 1}, P: [6]float64{1, 1, 0, 0, 0, -1}}, res: true, pt: util.FloatPoint{1, 0}},
	}

	for _, c := range cases {
		res, pt := c.YInflectionPoint()
		if res != c.res || !pt.Equal(c.pt) {
			t.Errorf("wrong %v %v", res, pt)
		}
	}

}

func TestFunction_SignXYQdx_Circle(t *testing.T) {
	log.Print("TestFunction_SectionSignDX_Circle")
	f := Function{P: [6]float64{1, 1, 0, 0, 0, -1}}

	cases := []util.FloatPoint{
		{0, -1},
		{0.1, -math.Sqrt(1 - math.Pow(0.1, 2))},
		{0.5, -math.Sqrt(1 - math.Pow(0.5, 2))},
		{1, 0},
		{0.1, math.Sqrt(1 - math.Pow(0.1, 2))},
		{0, 1},
	}

	for _, c := range cases {
		res, s := f.SignXYQdx(c)
		log.Print("  ", c, res, s)
	}

}

func TestFunction_SignXYQdx_Hyperbola1(t *testing.T) {
	log.Print("TestFunction_SignXYQdx_Hyperbola1")
	f := Function{P: [6]float64{1, -1, 0, 0, 0, -1}}

	cases := []util.FloatPoint{
		{1, 0},
		{2, -math.Sqrt(math.Pow(2, 2) - 1)},
		{2, math.Sqrt(math.Pow(2, 2) - 1)},
		{-1, 0},
		{-2, -math.Sqrt(math.Pow(2, 2) - 1)},
		{-2, math.Sqrt(math.Pow(2, 2) - 1)},
	}

	for _, c := range cases {
		res, s := f.SignXYQdx(c)
		log.Print("  ", c, res, s)
	}

}

func TestFunction_SignXYQdx_Hyperbola2(t *testing.T) {
	log.Print("TestFunction_SignXYQdx_Hyperbola2")
	f := Function{P: [6]float64{-1, 1, 0, 0, 0, -1}}

	cases := []util.FloatPoint{
		{0, 1},
		{-math.Sqrt(math.Pow(2, 2) - 1), 2},
		{math.Sqrt(math.Pow(2, 2) - 1), 2},
		{0, -1},
		{-math.Sqrt(math.Pow(2, 2) - 1), -2},
		{math.Sqrt(math.Pow(2, 2) - 1), -2},
	}

	for _, c := range cases {
		res, s := f.SignXYQdx(c)
		log.Print("  ", c, res, s)
	}

}
