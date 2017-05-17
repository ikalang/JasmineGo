package forbidden

import (
	"log"
	"traffic"
	"traffic/route"
)

const (
	NoInNoOut = iota
	IgnoreRelateAGV
	NoInOnlyOut
)

type ForbiddenArea struct {
	ID, Type, RelateID int
	Min, Max           route.IntPoint
	o                  traffic.OBB
}

var gFobiddenAreas map[int]ForbiddenArea
var gIsForbiddenAreaModified = false

func create(id, t, rid int, min, max route.IntPoint) ForbiddenArea {
	o := traffic.CreateOBB(route.FloatPoint{X: 0.5 * float64(min.X+max.X), Y: 0.5 * float64(min.Y+max.Y)},
		0.5*float64(max.X-min.X), 0.5*float64(max.Y-min.Y), 0)
	return ForbiddenArea{id, t, rid, min, max, o}
}

func Add(id, t, rid int, min, max route.IntPoint) bool {
	log.Printf("<forbidden.Add> id %d, type %d, AGV id %d,  min(%d, %d) max(%d, %d)\n",
		id, t, rid, min.X, min.Y, max.X, max.Y)
	defer log.Println("<forbidden.Add> exit")

	if _, ok := gFobiddenAreas[id]; ok {
		log.Printf("Forbidden Area %d already exist\n", id)
		return false
	}

	gFobiddenAreas[id] = create(id, t, rid, min, max)
	gIsForbiddenAreaModified = true
	return true
}

func Delete(id int) bool {
	log.Printf("<forbidden.Delete> id %d\n", id)
	defer log.Println("<forbidden.Delete> exit")

	if _, ok := gFobiddenAreas[id]; ok {
		log.Printf("can't find Forbidden Area %d \n", id)
		return false
	}

	delete(gFobiddenAreas, id)
	gIsForbiddenAreaModified = true
}

func Modify(id, t, rid int, min, max route.IntPoint) bool {
	log.Printf("<forbidden.Modify> id %d, type %d, AGV id %d,  min(%d, %d) max(%d, %d)\n",
		id, t, rid, min.X, min.Y, max.X, max.Y)
	defer log.Println("<forbidden.Modify> exit")

	if _, ok := gFobiddenAreas[id]; !ok {
		log.Printf("can't find Forbidden Area %d \n", id)
		return false
	}

	gFobiddenAreas[id] = create(id, t, rid, min, max)
	gIsForbiddenAreaModified = true
	return true
}

func IsAnyModified() bool {
	return gIsForbiddenAreaModified
}

func ResetModifiedFlag() {
	gIsForbiddenAreaModified = true
}
