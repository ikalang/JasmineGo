package route

import (
	"encoding/json"
	"traffic/util"
)

type Mark struct {
	Index    int
	Position util.IntPoint
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

type SubRoute struct {
	Type, MoveType                                        int
	Start, End                                            util.IntPoint
	MaxSpeed, MaxAcceleration, MaxDeceleration            int
	IsContinuousLock, IsContinuousLockWithNext, IsEndStop bool
	RefParams                                             [6]int32
	RefPoints                                             [2]util.IntPoint
}

func (sr SubRoute) String() string {
	data, err := json.Marshal(sr)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

const (
	Invalid = iota
	Straight
	QTurn
	Oblique
	UTurn
	STurn
	Linkage = 9
)

func (sr SubRoute) InOutDirection() (in, assist util.Direction) {
	switch sr.Type {
	case Straight:
		if sr.Start.X == sr.End.X {
			if sr.Start.Y < sr.End.Y {
				return util.YInc, util.YInc
			} else {
				return util.YDec, util.YDec
			}
		} else {
			if sr.Start.Y < sr.End.Y {
				return util.XInc, util.XInc
			} else {
				return util.XDec, util.XDec
			}
		}

	case QTurn:
		if sr.Start.X == sr.RefPoints[0].X {
			if sr.Start.Y < sr.RefPoints[0].Y {
				in = util.YInc
			} else {
				in = util.YDec
			}
			if sr.RefPoints[0].X < sr.End.X {
				assist = util.XInc
			} else {
				assist = util.XDec
			}
		} else {
			if sr.Start.X < sr.RefPoints[0].X {
				in = util.XInc
			} else {
				in = util.XDec
			}
			if sr.RefPoints[0].Y < sr.End.Y {
				assist = util.YInc
			} else {
				assist = util.YDec
			}
		}

	case Oblique:
		if abs(sr.Start.X-sr.End.X) == int(sr.RefParams[0]) {
			if sr.Start.Y < sr.End.Y {
				in = util.YInc
			} else {
				in = util.YDec
			}
			if sr.Start.X < sr.End.X {
				assist = util.XInc
			} else {
				assist = util.XDec
			}
		} else {
			if sr.Start.X < sr.End.X {
				in = util.XInc
			} else {
				in = util.XDec
			}
			if sr.Start.Y < sr.End.Y {
				assist = util.YInc
			} else {
				assist = util.YDec
			}
		}

	case UTurn:
		if sr.Start.X == sr.RefPoints[0].X {
			if sr.Start.Y < sr.RefPoints[1].X {
				in = util.YInc
			} else {
				in = util.YDec
			}
			if sr.RefPoints[0].X < sr.RefPoints[1].X {
				assist = util.XInc
			} else {
				assist = util.XDec
			}
		} else {
			if sr.Start.X < sr.RefPoints[0].X {
				in = util.XInc
			} else {
				in = util.XDec
			}
			if sr.RefPoints[0].Y < sr.RefPoints[1].Y {
				assist = util.YInc
			} else {
				assist = util.YDec
			}
		}

	case STurn:
		if sr.Start.X == sr.RefPoints[0].X {
			if sr.Start.Y < sr.RefPoints[0].Y {
				in = util.YInc
			} else {
				in = util.YDec
			}
			if sr.RefPoints[0].X < sr.RefPoints[1].X {
				assist = util.XInc
			} else {
				assist = util.XDec
			}
		} else {
			if sr.Start.X < sr.RefPoints[0].X {
				in = util.XInc
			} else {
				in = util.XDec
			}
			if sr.RefPoints[0].Y < sr.RefPoints[1].Y {
				assist = util.YInc
			} else {
				assist = util.YDec
			}
		}

	default:
		return util.DirErr, util.DirErr
	}

	return in, assist
}

func (sr SubRoute) HeadingOnSubRoute(p util.IntPoint, isAstern bool) (d util.Degree) {
	in, assist := sr.InOutDirection()

	switch sr.Type {
	case Straight:
		d = in.ToDegree()

	case QTurn:
		if p.Equal(sr.Start) {
			d = in.ToDegree()
		} else {
			d = assist.ToDegree()
		}

	case Oblique:
		d = in.ToDegree()

	case UTurn:
		if p.Equal(sr.Start) {
			d = in.ToDegree()
		} else {
			d = in.Opposite().ToDegree()
		}

	case STurn:
		d = in.ToDegree()

	default:
		return -1
	}

	if isAstern {
		return d.Opposite()
	} else {
		return d
	}
}
