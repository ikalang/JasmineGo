package legume

import (
	"bytes"
	"container/ring"
	"fmt"
	"log"
	"traffic/track"
	"traffic/util"
)

const gLEGUME_RINGBUFFER_SIZE = 5000

type bean struct {
	index int
	OBB
	isStop bool
}

type Legume struct {
	n                     int
	ringBuf               *ring.Ring
	self, claimed, trying *ring.Ring
}

func (l *Legume) Init(n int) {
	l.ringBuf = ring.New(n)
	l.self = l.ringBuf
	l.claimed = l.ringBuf
	l.trying = l.ringBuf
}

func (l *Legume) Reset() {
	l.claimed = l.self.Next()
	l.trying = l.claimed
}

func (l *Legume) String() string {
	var buf bytes.Buffer
	for r := l.self; r != l.trying; r = r.Next() {
		buf.Write([]byte(r.Value.(bean).OBB.String()))
		buf.WriteByte('\n')
	}
	return buf.String()
}

func (l *Legume) AppendBean(o OBB) {
	b := bean{OBB: o}

	if l.self == l.trying {
		b.index = 0
		l.trying.Value = b
		l.trying = l.trying.Next()
		return
	}

	tail := l.trying.Prev().Value.(bean)
	if l.self == l.trying || !tail.Center.Equal(b.Center) || !tail.deg.Equal(b.deg) ||
		tail.XHalfLength != b.XHalfLength || tail.YHalfLength != b.YHalfLength {
		b.index = tail.index + 1
		b.isStop = false
		l.trying.Value = b
		l.trying = l.trying.Next()
	}
}

func (l *Legume) GrowCenter(c util.FloatPoint, xh, yh float64, deg util.Degree) {
	l.AppendBean(CreateOBB(c, xh, yh, deg))
}

func (l *Legume) GrowFrontRear(f, r util.FloatPoint, xh, yh float64) {
	l.GrowCenter(f.CenterPoint(r), xh, yh, f.DegreeTo(r))
}

func (l *Legume) IsOverlapWithOBB(start, end *ring.Ring, o OBB) bool {
	for r := start; r != end; r = r.Next() {
		if o.IsOverlap(r.Value.(bean).OBB) {
			return true
		}
	}
	return false
}

func (l *Legume) IsOverlapWithLegume(start, end *ring.Ring, q *Legume, qStart, qEnd *ring.Ring) (bool, *ring.Ring, *ring.Ring) {
	for r := start; r != end; r = r.Next() {
		for j := qStart; j != qEnd; j = j.Next() {
			if r.Value.(bean).OBB.IsOverlap(j.Value.(bean).OBB) {
				return true, r, j
			}
		}
	}
	return false, nil, nil
}

func (l *Legume) GrowStraightSlice(start, end util.IntPoint) {
	if start.Equal(end) {
		return
	}

	f := track.Function{Start: start.ToFloatPoint(), End: end.ToFloatPoint()}
	if start.X == end.X {
		f.P[3] = 1
		f.P[5] = start.ToFloatPoint().X
	} else {
		f.P[3] = float64(end.Y-start.Y) / float64(end.X-start.X)
		f.P[4] = -1
		f.P[5] = start.ToFloatPoint().Y - f.P[3]*start.ToFloatPoint().X
	}
}

const (
	gFactorLinearDX    = 100
	gFactorQuadraticDX = 10
	gAGVWheelbase      = 1200
)

func (l *Legume) GrowAlongTrack(t track.Track, xh, yh float64) error {
	var f, r util.FloatPoint
	rn := 0
	for fn := 0; fn < len(t.Front); {
		log.Printf("fn %d", fn)
		nextFront, res := t.Front[fn].NextPoint(f, gFactorQuadraticDX)
		if res == track.NotFount {
			return fmt.Errorf("Can't find next front %d", fn)
		}

		for rn < len(t.Rear) {
			nextRear, res := t.Rear[rn].NextPointRef(r, nextFront, gAGVWheelbase)
			if res == track.NotFount {
				rn++
				nextRear = util.FloatPointZero
				continue
			}

			isBreak := nextFront.Distance(nextRear) < gAGVWheelbase+10
			if isBreak {
				l.GrowFrontRear(nextFront, nextRear, xh, yh)
			}

			if res == track.EndPoint {
				rn++
				r = util.FloatPointZero
			} else {
				r = nextRear
			}

			if isBreak {
				break
			}
		}

		if res == track.EndPoint {
			fn++
			f = util.FloatPointZero
		} else {
			f = nextFront
		}
	}

	return nil
}

var gBaseLegume map[int]*Legume = make(map[int]*Legume)

const (
	STARIGHT_HALF_HEIGHT = 850
	STARIGHT_HALF_WEIGHT = 170
)

func stemLegume(moveTypeID int) *Legume {
	key := moveTypeID*100 + int(util.XInc)*10 + int(util.YInc)

	l, ok := gBaseLegume[key]
	if ok {
		return l
	}

	t, err := track.GetTrack(moveTypeID)
	if err != nil {
		return nil
	}

	l = &Legume{}
	l.Init(500)
	l.GrowFrontRear(t.Front[0].Start, t.Rear[0].Start, STARIGHT_HALF_HEIGHT, STARIGHT_HALF_WEIGHT)

	if l.GrowAlongTrack(t, 0, 0) != nil {
		return nil
	}

	os, err := track.GetOBBSize(moveTypeID)
	if err != nil {
		return nil
	}

	for i := l.self; i != l.trying; i = i.Next() {
		b := i.Value.(bean)
		b.Center = b.OBB.Centroid(os.Front, os.Rear, os.Inner, os.Outer)
		b.XHalfLength = 0.5 * (os.Front + os.Rear)
		b.YHalfLength = 0.5 * (os.Inner + os.Outer)
	}

	l.GrowFrontRear(t.Front[len(t.Front)-1].End, t.Rear[len(t.Rear)-1].End, STARIGHT_HALF_HEIGHT, STARIGHT_HALF_WEIGHT)

	gBaseLegume[key] = l
	return gBaseLegume[key]
}

func BaseLegume(moveTypeID int, d1, d2 util.Direction) *Legume {
	key := moveTypeID*100 + int(d1)*10 + int(d2)

	l, ok := gBaseLegume[key]
	if ok {
		return l
	}

	sl := stemLegume(moveTypeID)
	if sl == nil {
		return nil
	}

	if d1 == util.XInc && d2 == util.YInc {
		return sl
	}

	l.Init(500)

	for i := l.self; i != l.trying; i = i.Next() {
		l.AppendBean(i.Value.(bean).Transform(d1, d2))
	}

	gBaseLegume[key] = l
	return gBaseLegume[key]
}
