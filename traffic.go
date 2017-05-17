package traffic

import (
	"encoding/json"
	"traffic/route"
)

const AGVWheelbase = 1200
const ToleranceParallel = 6
const ToleranceVertical = 11

type AGVStatus struct {
	ID, Heading, Speed, MotionStatus, Priority int
	Pos                                        route.IntPoint
	ControlCode                                int
	Target                                     route.Mark
}

func (as AGVStatus) String() string {
	data, err := json.Marshal(as)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

type AGVCommand struct {
	Type, Heading                     int
	Target                            route.IntPoint
	MaxStraightSpeed, MaxSpecialSpeed int
	IsNeedAccurateStop                bool
	CommandRefs                       [7]int
}

const (
	Undefined = iota
	Idle
	Run
	Fault
	Pause
	Manual
)

const (
	MSUncertain = iota
	MSStop
	MSStraight
	MSQTurn
	MSOblique
	MSUTurn
	MSSTurn
)

type AGVResult struct {
	ID, RunStatus int
	Command       AGVCommand
	ErrorCode     int
}

type AGV struct {
	ID, Heading, Speed, MotionStatus, ControlCode, Priority int
	Position                                                route.FloatPoint
	Routes                                                  []route.SubRoute
	Orientation                                             route.Direction
	Current, CommandRM, Claim, Trying, Target               route.Mark
	Command                                                 AGVCommand
	RunStatus, TryFailedCount, ErrorCode                    int
}

func (v AGV) IsArrivalWithTolerance(p route.IntPoint) bool {
	if v.MotionStatus != MSStop && v.MotionStatus != MSStraight {
		return false
	}

	if v.Orientation == route.XInc || v.Orientation == route.XDec {
		return route.FloatEqualTolerance(v.Position.X, p.ToFloatPoint().X, ToleranceParallel) &&
			route.FloatEqualTolerance(v.Position.Y, p.ToFloatPoint().Y, ToleranceVertical)
	} else if v.Orientation == route.YInc || v.Orientation == route.YDec {
		return route.FloatEqualTolerance(v.Position.X, p.ToFloatPoint().X, ToleranceVertical) &&
			route.FloatEqualTolerance(v.Position.Y, p.ToFloatPoint().Y, ToleranceParallel)
	} else {
		return false
	}
}

func (v AGV) IsAGVOnLineWithTolerance(start, end route.IntPoint) bool {
	if v.Orientation == route.XInc || v.Orientation == route.XDec && start.Y == end.Y {
		return route.FloatEqualTolerance(start.ToFloatPoint().Y, v.Position.Y, ToleranceVertical)
	} else if v.Orientation == route.YInc || v.Orientation == route.YDec && start.X == end.X {
		return route.FloatEqualTolerance(start.ToFloatPoint().X, v.Position.X, ToleranceVertical)
	} else {
		return false
	}
}

func (v AGV) IsAGVOnSegmentWithTolerance(start, end route.IntPoint) bool {
	if v.Orientation == route.XInc || v.Orientation == route.XDec && start.Y == end.Y {
		return route.FloatEqualTolerance(start.ToFloatPoint().Y, v.Position.Y, ToleranceVertical) &&
			route.FloatInOpenInterval(v.Position.X, start.ToFloatPoint().X, end.ToFloatPoint().X, ToleranceParallel)
	} else if v.Orientation == route.YInc || v.Orientation == route.YDec && start.X == end.X {
		return route.FloatEqualTolerance(start.ToFloatPoint().X, v.Position.X, ToleranceVertical) &&
			route.FloatInOpenInterval(v.Position.Y, start.ToFloatPoint().Y, end.ToFloatPoint().Y, ToleranceParallel)
	} else {
		return false
	}
}

func (v AGV) IsAGVOnSubRoute(index int) bool {
	sr := v.Routes[index]

	switch sr.Type {
	case route.Straight:
		if v.IsAGVOnSegmentWithTolerance(sr.Start, sr.End) {
			return true
		}

	case route.QTurn:
		if v.IsAGVOnSegmentWithTolerance(sr.Start, sr.RefPoints[0]) || v.IsAGVOnSegmentWithTolerance(sr.RefPoints[1], sr.End) {
			return true
		}

	case route.Oblique:
		if sr.Start.X-sr.End.X == int(sr.RefParams[0]) {
			sr.RefPoints[0] = route.IntPoint{X: sr.Start.X, Y: sr.End.Y}
			sr.RefPoints[1] = route.IntPoint{X: sr.End.X, Y: sr.Start.Y}
		} else {
			sr.RefPoints[0] = route.IntPoint{X: sr.End.X, Y: sr.Start.Y}
			sr.RefPoints[1] = route.IntPoint{X: sr.Start.X, Y: sr.End.Y}
		}

		if v.IsAGVOnSegmentWithTolerance(sr.Start, sr.RefPoints[0]) || v.IsAGVOnSegmentWithTolerance(sr.End, sr.RefPoints[1]) {
			return true
		}

	case route.UTurn:
		if v.IsAGVOnSegmentWithTolerance(sr.Start, sr.RefPoints[0]) || v.IsAGVOnSegmentWithTolerance(sr.End, sr.RefPoints[1]) {
			return true
		}

	case route.STurn:
		if v.IsAGVOnSegmentWithTolerance(sr.Start, sr.RefPoints[0]) || v.IsAGVOnSegmentWithTolerance(sr.End, sr.RefPoints[1]) {
			return true
		}

	default:
		return false
	}

	return false
}

func (v AGV) UpdateCurrentMark() {
	if v.Current.Index == v.Claim.Index {
		return
	}

	newIdx := v.Current.Index
	switch v.MotionStatus {
	case MSStop:
		fallthrough
	case MSStraight:
		for i := newIdx; i < v.Claim.Index; i++ {
			if v.IsAGVOnSubRoute(i) {
				if v.IsArrivalWithTolerance(v.Routes[i].End) && i < v.Claim.Index {
					newIdx = i + 1
				} else {
					newIdx = i
				}
				break
			}
		}

	case MSQTurn:
		if v.Current.Index < v.CommandRM.Index+1 {
			sr := v.Routes[v.Current.Index+1]
			if sr.Type == route.QTurn &&
				route.FloatInCloseInterval(v.Position.X, sr.Start.ToFloatPoint().X, sr.End.ToFloatPoint().X, 0.3) &&
				route.FloatInCloseInterval(v.Position.Y, sr.Start.ToFloatPoint().Y, sr.End.ToFloatPoint().Y, 0.3) {
				newIdx = v.Current.Index + 1
			}
		}

	case MSOblique:
		if v.Current.Index < v.Claim.Index+1 {
			if v.Routes[v.Current.Index+1].Type == route.Oblique {
				newIdx = v.Current.Index + 1
			}
		}

	case MSUTurn:
		if v.Current.Index < v.Claim.Index+1 {
			if v.Routes[v.Current.Index+1].Type == route.UTurn {
				newIdx = v.Current.Index + 1
			}
		}

	case MSSTurn:
		if v.Current.Index < v.Claim.Index+1 {
			if v.Routes[v.Current.Index+1].Type == route.STurn {
				newIdx = v.Current.Index + 1
			}
		}

	default:
		return
	}

	v.Current.Index = newIdx
}

func (v AGV) Reset() {
	v.Current.Index = 0
	v.Claim.Index = 0
	v.Target.Index = 0
	v.CommandRM.Index = 0

	v.Claim.Position = v.Position.ToIntPoint()

	v.RunStatus = Idle
}

func (v AGV) ResetTrying() {
	v.Trying = v.Claim
	v.TryFailedCount = 0
}
