package legume

import (
	"log"
	"testing"
	"traffic/util"
)

func TestBaseLegume(t *testing.T) {
	l := baseLegume(10, util.XInc, util.YInc)
	log.Print(l)
}
