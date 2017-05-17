package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"traffic/util"
)

const (
	XYQ = iota
	XL
	XYL
	XQYL
)

type Function struct {
	Type       int
	Start, End util.FloatPoint
	sign       float64
	P          [6]float64 // f(x,y) = ax2 + by2 + cxy + dx + ey + f
}

func (f Function) String() string {
	data, err := json.Marshal(f)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (f Function) Evaluate(p util.FloatPoint) float64 {
	return f.P[0]*math.Pow(p.X, 2) + f.P[1]*math.Pow(p.Y, 2) + f.P[2]*p.X*p.Y +
		f.P[3]*p.X + f.P[4]*p.Y + f.P[5]
}

func (f Function) IsPointOnFunction(p util.FloatPoint) bool {
	return f.Evaluate(p) < 50
}

func (f Function) IsEndPoint(p util.FloatPoint) bool {
	return f.End.Equal(p)
}

func (f Function) IsStartPoint(p util.FloatPoint) bool {
	return f.Start.Equal(p)
}

func (f Function) IsCenterPoint(p util.FloatPoint) bool {
	return util.FloatInCloseInterval(p.X, f.Start.X, f.End.X, 0.01) &&
		util.FloatInCloseInterval(p.Y, f.Start.Y, f.End.Y, 0.01) &&
		f.IsPointOnFunction(p)
}

func (f Function) Verify() (bool, error) {
	if f.Evaluate(f.Start) != 0 {
		return false, fmt.Errorf("Start Evaluate is %d, not equal 0", f.Evaluate(f.Start))
	}

	if f.Evaluate(f.End) != 0 {
		return false, fmt.Errorf("End Evaluate is %d, not equal 0", f.Evaluate(f.End))
	}

	return true, nil
}

func (f *Function) recognizeType() {
	if !util.FloatEqual(f.P[1], 0) {
		f.Type = XYQ
	} else if util.FloatEqual(f.P[2]*f.Start.X+f.P[4], 0) {
		f.Type = XL
	} else if util.FloatEqual(f.P[0], 0) && util.FloatEqual(f.P[2], 0) {
		f.Type = XYL
	} else {
		f.Type = XQYL
	}
}

func (f *Function) recognizeSign() {
	switch f.Type {
	case XYQ:
		f.sign = float64(f.SectionSignDX())

	case XL:
		if f.Start.Y < f.End.Y {
			f.sign = 1
		} else {
			f.sign = -1
		}
	}

}

func (f Function) SignXYQdy(p util.FloatPoint) (bool, int) {
	x1 := f.XSolverXYQ(1, p.Y)
	x2 := f.XSolverXYQ(-1, p.Y)

	if util.FloatEqual(x1, p.X) {
		return true, 1
	} else if util.FloatEqual(x2, p.X) {
		return true, -1
	} else {
		return false, 0
	}
}

func (f Function) SignXYQdx(p util.FloatPoint) (bool, int) {
	y1 := f.YSolverXYQ(1, p.X)
	y2 := f.YSolverXYQ(-1, p.X)

	if util.FloatEqual(y1, p.Y) {
		return true, 1
	} else if util.FloatEqual(y2, p.Y) {
		return true, -1
	} else {
		log.Print(y1, y2, p.Y)
		panic("")
		return false, 0
	}
}

func (f Function) SectionSignDX() int {
	switch f.Type {
	case XYQ:
		resStart, SignStart := f.SignXYQdx(f.Start)
		resEnd, SignEnd := f.SignXYQdx(f.End)
		if !resStart || !resEnd {
			panic("")
		}

		if resStart || resEnd {
			if SignStart == SignEnd {
				return SignStart
			} else if SignStart == 0 {
				return SignEnd
			} else if SignEnd == 0 {
				return SignStart
			}
		}

	case XL:
		if f.Start.Y < f.End.Y {
			f.sign = 1
		} else {
			f.sign = -1
		}
	}

	return 0
}

func (f Function) SectionSignDY() int {

	switch f.Type {
	case XYQ:
		resStart, SignStart := f.SignXYQdy(f.Start)
		resEnd, SignEnd := f.SignXYQdy(f.End)
		if !resStart || !resEnd {
			return 0
			panic("")
		}

		if resStart || resEnd {
			if SignStart == SignEnd {
				return SignStart
			} else if SignStart == 0 {
				return SignEnd
			} else if SignEnd == 0 {
				return SignStart
			}
		}

	case XL:
		if f.Start.X < f.End.X {
			f.sign = 1
		} else {
			f.sign = -1
		}
	}

	return 0
}

func (f Function) XSolverXYQ(sign, Y float64) float64 {
	P := f.P
	B := 2 * P[0]
	C := P[2]*Y + P[3]
	D := C*C - 4*P[0]*(P[1]*Y*Y+P[4]*Y+P[5])

	if util.FloatEqual(D, 0) {
		return -C / B
	} else if D < 0 {
		return math.NaN()
	}

	return (-C + math.Copysign(math.Sqrt(D), sign)) / B
}

func (f Function) YSolverXYQ(sign, X float64) float64 {
	P := f.P
	B := 2 * P[1]
	C := P[2]*X + P[4]
	D := C*C - 4*P[1]*(P[0]*X*X+P[3]*X+P[5])

	if util.FloatEqual(D, 0) {
		return -C / B
	} else if D < 0 {
		return math.NaN()
	}

	Dsq := math.Sqrt(D)
	return (-C + Dsq*float64(sign)) / B
}

func (f Function) XInflectionPoint() (_ bool, pt util.FloatPoint) {
	P := f.P
	U := -4 * P[0] * P[1]
	S := math.Pow(P[2], 2) + U

	T := P[2]*P[4] - 2*P[1]*P[3]
	T2 := math.Pow(T, 2)
	C := P[2]*P[3]*P[4] - P[1]*P[3]*P[3] - P[2]*P[2]*P[5]
	D := T2 - S*C/P[0]

	if D < 0 {
		return false, pt
	}

	//sectionSign := f.SectionSignDX()
	//log.Print("X section ",sectionSign)
	Dsq := math.Sqrt(D)

	x1 := (-T + Dsq) / S
	if util.FloatInCloseInterval(x1, f.Start.X, f.End.X, math.SmallestNonzeroFloat64) {
		y1 := f.YSolverXYQ(1, x1)
		if !math.IsNaN(y1) {
			log.Print(1, x1, y1)
		}
	}

	x2 := (-T - Dsq) / S
	if util.FloatInCloseInterval(x2, f.Start.X, f.End.X, math.SmallestNonzeroFloat64) {
		y2 := f.YSolverXYQ(-1, x2)
		if !math.IsNaN(y2) {
			log.Print(-1, x2, y2)
		}
	}

	return false, pt
}

func (f Function) YInflectionPoint() (_ bool, pt util.FloatPoint) {
	P := f.P
	U := -4 * P[0] * P[1]
	S := math.Pow(P[2], 2) + U

	T := P[2]*P[3] - 2*P[0]*P[4]
	T2 := math.Pow(T, 2)
	C := P[2]*P[3]*P[4] - P[0]*P[4]*P[4] - P[2]*P[2]*P[5]
	D := T2 - S*C/P[1]

	if D < 0 {
		return false, pt
	}

	Dsq := math.Sqrt(D)
	sectionSign := f.SectionSignDY()

	pt.Y = (-T + float64(sectionSign)*Dsq) / S
	if !util.FloatInCloseInterval(pt.Y, f.Start.Y, f.End.Y, math.SmallestNonzeroFloat64) {
		return false, pt
	}

	pt.X = f.XSolverXYQ(float64(sectionSign), pt.Y)
	if math.IsNaN(pt.X) {
		return false, pt
	}

	return true, pt
}

type QuadraticCurve struct {
	StartPoint, CenterPoint, EndPoint util.FloatPoint
	RA, RB, PHI                       float64 // RA RB unit m
}

type QuadraticCurveArray []QuadraticCurve

func (q *QuadraticCurve) standardize() {
	standardizeShift := util.FloatPoint{-7, 0}

	q.StartPoint = q.StartPoint.Shift(standardizeShift).Scale(100)
	q.CenterPoint = q.CenterPoint.Shift(standardizeShift).Scale(100)
	q.EndPoint = q.EndPoint.Shift(standardizeShift).Scale(100)
	q.RA *= 100
	q.RB *= 100
}

type QuadraticCurveTrack struct {
	Name        string
	Front, Rear QuadraticCurveArray
}

type Track struct {
	Front, Rear []Function
}

func (t Track) String() string {
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (t Track) ToCCode() {
	fmt.Printf("t->fa.n = %d;\n", len(t.Front))
	fmt.Println("t->fa.f = malloc(t->fa.n * sizeof(func));")

	for i := 0; i < len(t.Front); i++ {
		f := t.Front[i]
		fmt.Printf("t->fa.f[%d] = (func) {\n", i)
		fmt.Printf(".stpt = (point) {%f, %f},\n", f.Start.X, f.Start.Y)
		fmt.Printf(".edpt = (point) {%f, %f},\n", f.End.X, f.End.Y)
		fmt.Printf(".p = {%f, %f, %f, %f, %f, %f}\n", f.P[0], f.P[1], f.P[2], f.P[3], f.P[4], f.P[5])
		fmt.Print("};\n")
	}

	fmt.Printf("t->ra.n = %d;\n", len(t.Rear))
	fmt.Println("t->ra.f = malloc(t->ra.n * sizeof(func));")

	for i := 0; i < len(t.Rear); i++ {
		f := t.Rear[i]
		fmt.Printf("t->ra.f[%d] = (func) {\n", i)
		fmt.Printf(".stpt = (point) {%f, %f},\n", f.Start.X, f.Start.Y)
		fmt.Printf(".edpt = (point) {%f, %f},\n", f.End.X, f.End.Y)
		fmt.Printf(".p = {%f, %f, %f, %f, %f, %f}\n", f.P[0], f.P[1], f.P[2], f.P[3], f.P[4], f.P[5])
		fmt.Print("};\n")
	}
}

func GetTrack(moveTypeID int) (t Track, err error) {
	var qct QuadraticCurveTrack
	err = qct.LoadFromJSONFile("track/data/QTurn4to7.json")
	if err != nil {
		return t, err
	}

	t = qct.ToTrack()
	log.Print(t)
	return t, nil
}

type OBBsize struct {
	Type                      int
	Front, Rear, Inner, Outer float64
}

func GetOBBSize(moveTypeID int) (OBBsize, error) {
	var t OBBsize
	return t, nil
}

func calcLineFunc(f *Function) {
	f.P[0] = 0
	f.P[1] = 0
	f.P[2] = 0

	if f.Start.X == f.End.X {
		f.P[3] = -1
		f.P[4] = 0
		f.P[5] = 0
	} else {
		k := (f.End.Y - f.Start.Y) / (f.End.X - f.Start.X)
		b := f.End.Y - k*f.End.X
		f.P[3] = k
		f.P[4] = -1
		f.P[5] = b
	}
}

func calcEllipseFunc(f *Function, qc QuadraticCurve) {
	sinPHI := math.Sin(qc.PHI)
	cosPHI := math.Cos(qc.PHI)

	f.P[0] = math.Pow(qc.RA, 2)*math.Pow(sinPHI, 2) + math.Pow(qc.RB, 2)*math.Pow(cosPHI, 2)
	f.P[1] = math.Pow(qc.RA, 2)*math.Pow(cosPHI, 2) + math.Pow(qc.RB, 2)*math.Pow(sinPHI, 2)
	f.P[2] = 2 * (math.Pow(qc.RA, 2) - math.Pow(qc.RB, 2)) * sinPHI * cosPHI
	f.P[3] = -2*f.P[0]*qc.CenterPoint.X - f.P[2]*qc.CenterPoint.Y
	f.P[4] = -2*f.P[1]*qc.CenterPoint.Y - f.P[2]*qc.CenterPoint.X
	f.P[5] = f.P[0]*math.Pow(qc.CenterPoint.X, 2) + +f.P[1]*math.Pow(qc.CenterPoint.Y, 2) + f.P[2]*qc.CenterPoint.X*qc.CenterPoint.Y - math.Pow(qc.RA, 2)*math.Pow(qc.RB, 2)
}

func calcHyperbolaFunc(f *Function, qc QuadraticCurve) {
	sinPHI := math.Sin(qc.PHI)
	cosPHI := math.Cos(qc.PHI)

	f.P[0] = math.Pow(qc.RB, 2)*math.Pow(cosPHI, 2) - math.Pow(qc.RA, 2)*math.Pow(sinPHI, 2)
	f.P[1] = math.Pow(qc.RB, 2)*math.Pow(sinPHI, 2) - math.Pow(qc.RA, 2)*math.Pow(cosPHI, 2)
	f.P[2] = -2 * (math.Pow(qc.RA, 2) + math.Pow(qc.RB, 2)) * sinPHI * cosPHI
	f.P[3] = -2*f.P[0]*qc.CenterPoint.X - f.P[2]*qc.CenterPoint.Y
	f.P[4] = -2*f.P[1]*qc.CenterPoint.Y - f.P[2]*qc.CenterPoint.X
	f.P[5] = f.P[0]*math.Pow(qc.CenterPoint.X, 2) + f.P[1]*math.Pow(qc.CenterPoint.Y, 2) + f.P[2]*qc.CenterPoint.X*qc.CenterPoint.Y - math.Pow(qc.RA, 2)*math.Pow(qc.RB, 2)
}

func (q QuadraticCurve) toFunction() (f Function) {
	f.Start = q.StartPoint
	f.End = q.EndPoint

	if q.RA == 0 && q.RB == 0 {
		calcLineFunc(&f)
	} else if q.RA < 0 || q.RB < 0 {
		calcHyperbolaFunc(&f, q)
	} else {
		calcEllipseFunc(&f, q)
	}

	a := math.Log10(math.Abs(f.P[0]))
	b := math.Log10(math.Abs(f.P[1]))
	c := -int((a + b) * 0.5)

	factor := math.Pow(10, float64(c))

	if factor != 0 {
		for i := 0; i < 6; i++ {
			f.P[i] *= factor
		}
	}

	f.recognizeType()
	return f
}

func (qca QuadraticCurveArray) ToFunctionArrayAndSplit() (fArray []Function) {
	for _, qc := range qca {
		fArray = append(fArray, qc.toFunction())

		f := fArray[len(fArray)-1]

		resX, XinflectionPt := f.XInflectionPoint()
		resY, YinflectionPt := f.YInflectionPoint()
		if resX || resY {
			var infPt util.FloatPoint
			if resX {
				infPt = XinflectionPt
			} else {
				infPt = YinflectionPt
			}
			log.Print(resX, resY, infPt)
			fArray = append(fArray, (fArray)[len(fArray)-1])

			fArray[len(fArray)-2].End = infPt
			fArray[len(fArray)-1].Start = infPt
			fArray[len(fArray)-1].End = qc.EndPoint
		}
	}

	for i := 0; i < len(fArray); i++ {
		fArray[i].recognizeSign()
	}

	return fArray
}

func (qct *QuadraticCurveTrack) standardize() {
	for i := 0; i < len(qct.Front); i++ {
		f := &qct.Front[i]
		f.standardize()
	}

	for i := 0; i < len(qct.Rear); i++ {
		f := &qct.Rear[i]
		f.standardize()
	}
}

func (qct *QuadraticCurveTrack) LoadFromJSONFile(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		log.Print(err)
		return err
	}

	data, err := ioutil.ReadAll(f)
	err = json.Unmarshal(data, &qct)
	if err != nil {
		return err
	}

	qct.standardize()
	return nil
}

func (qct QuadraticCurveTrack) ToTrack() (t Track) {
	log.Print("Front")
	t.Front = qct.Front.ToFunctionArrayAndSplit()
	log.Print("Rear")
	t.Rear = qct.Rear.ToFunctionArrayAndSplit()
	return t
}

func (f Function) dxXYQ(current util.FloatPoint, interval float64) (x0 float64) {
	x0 = current.X + math.Copysign(1.0, f.End.X-f.Start.X)

	if math.Abs(x0-current.X) > math.Abs(f.End.X-current.X) {
		return f.End.X
	}

	p := f.P

	biv := 1 / p[1]
	sbiv := math.Pow(biv, 2)
	K := 0.5 * p[2] * sbiv
	H := 1 + p[2]*K - p[0]*biv
	I := p[2]*p[4]*sbiv + p[2]*biv*current.Y - 2*current.X - p[3]*biv
	J := 0.5 * (p[4] + 2*p[1]*current.Y) * sbiv
	P := math.Pow(current.X, 2) + math.Pow(current.Y, 2) - math.Pow(interval, 2) +
		0.5*math.Pow(p[4], 2)*sbiv - p[5]*biv + p[4]*biv*current.Y
	L := math.Pow(p[2], 2) - 4*p[0]*p[1]
	M := p[2]*p[4] - 2*p[1]*p[3]

	for i := 0; i < 2; i++ {
		d := math.Pow(p[2]*x0+p[4], 2) - 4*p[1]*(p[0]*math.Pow(x0, 2)+p[3]*x0+p[5])
		if d < 0 {
			return math.MaxFloat64
		}

		sd := math.Sqrt(d)
		fx := H*math.Pow(x0, 2) + I*x0 - f.sign*(K*x0+J)*sd + P
		dfx := 2*H*x0 + I - f.sign*(K*sd+(K*x0+J)*(L*x0+M)/sd)
		x0 -= fx / dfx
	}

	return x0
}

func (f Function) dxXQYL(current, center util.FloatPoint, interval float64) (x0 float64) {
	x0 = current.X + 5

	p := f.P
	for i := 0; i < 2; i++ {
		C := 1 / (p[2]*x0 + p[4])
		A := (p[0]*math.Pow(x0, 2) + p[3]*x0 + p[5]) * C

		fx := math.Pow(x0-center.X, 2) + math.Pow(A+center.Y, 2) + math.Pow(interval, 2)
		dA := (p[0]*(2-p[2])*math.Pow(x0, 2) + 2*p[0]*p[4]*x0 + p[3]*p[4] - p[2]*p[5]) * math.Pow(C, 2)
		dfx := 2*(x0-center.X) + 2*(A+center.Y)*dA
		x0 -= fx / dfx
	}

	return x0
}

func (f Function) dxXYL(current, center util.FloatPoint, interval float64) (x0 float64) {
	x0 = current.X
	p := f.P

	d_e := p[3] / p[4]
	f_e := p[5] / p[4]
	for i := 0; i < 2; i++ {
		fx := math.Pow(x0-center.X, 2) + math.Pow(d_e*x0+f_e+center.Y, 2) - math.Pow(interval, 2)
		dfx := 2 * ((1+math.Pow(d_e, 2))*x0 + d_e*(f_e+center.Y) - center.X)
		x0 -= fx / dfx
	}

	return x0
}

const (
	NotFount = iota
	StartPoint
	CenterPoint
	EndPoint
)

func (f Function) NextPoint(current util.FloatPoint, interval float64) (next util.FloatPoint, res int) {
	defer func() {
		log.Print(current, next, res)
	}()
	p := f.P

	if current.Equal(util.FloatPointZero) {
		return f.Start, StartPoint
	}

	switch f.Type {
	case XYQ:
		next.X = f.dxXYQ(current, interval)
		if !util.FloatInCloseInterval(next.X, f.Start.X, f.End.X, 0.01) {
			return f.End, EndPoint
		}
		next.Y = f.YSolverXYQ(f.sign, next.X)

	case XL:
		next = util.FloatPoint{f.Start.X, current.Y + f.sign*interval}
		if !util.FloatInCloseInterval(next.Y, f.Start.Y, f.End.Y, 0.01) {
			return f.End, EndPoint
		}

	case XQYL:
		next.X = f.dxXQYL(current, current, interval)
		if !util.FloatInCloseInterval(next.X, f.Start.X, f.End.X, 0.01) {
			return f.End, EndPoint
		}

		C := p[2]*next.X + p[4]
		next.Y = -(p[0]*math.Pow(next.X, 2) + p[3]*next.X + p[5]) / C

	case XYL:
		next.X = current.X + math.Copysign(interval/math.Sqrt(1+math.Pow(p[3]/p[4], 2)), f.End.X-f.Start.X)
		if !util.FloatInCloseInterval(next.X, f.Start.X, f.End.X, 0.01) {
			return f.End, EndPoint
		}
		next.Y = -(p[3]*next.X + p[5]) / p[4]

	default:
		return next, NotFount
	}

	if next.Equal(current) {
		panic("")
	}

	if next.Equal(f.End) {
		return f.End, EndPoint
	}

	return next, CenterPoint
}

func (f Function) NextPointRef(current, ref util.FloatPoint, interval float64) (next util.FloatPoint, res int) {
	p := f.P

	if current.Equal(util.FloatPointZero) {
		return f.Start, StartPoint
	}

	switch f.Type {
	case XYQ:
		next.X = f.dxXYQ(ref, interval)
		if !util.FloatInCloseInterval(next.X, f.Start.X, f.End.X, 0.01) {
			return f.End, EndPoint
		}
		next.Y = f.YSolverXYQ(f.sign, next.X)

	case XL:
		next.X = f.Start.X
		ssub := math.Pow(interval, 2) - math.Pow(ref.X-current.X, 2)
		if ssub < 0 {
			log.Print(ssub, current, ref)
			//panic("ssub <0")
		}
		next.Y = ref.Y + math.Copysign(math.Sqrt(ssub), current.Y-ref.Y)
		if !util.FloatInCloseInterval(next.Y, f.Start.Y, f.End.Y, 0.01) {
			return f.End, EndPoint
		}

	case XYL:
		next.X = f.dxXYL(current, ref, interval)
		next.Y = -(p[3]*next.X + p[5]) / p[4]
		if !util.FloatInCloseInterval(next.Y, f.Start.Y, f.End.Y, 0.01) {
			return f.End, EndPoint
		}

	case XQYL:
		next.X = f.dxXQYL(current, ref, interval)
		if !util.FloatInCloseInterval(next.X, f.Start.X, f.End.X, 0.01) {
			return f.End, EndPoint
		}

		C := p[2]*next.X + p[4]
		next.Y = -(p[0]*math.Pow(next.X, 2) + p[3]*next.X + p[5]) / C

	default:
		return next, NotFount
	}

	return next, CenterPoint
}
