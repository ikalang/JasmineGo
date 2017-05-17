package route

import (
	"log"
	"testing"
)

func TestSubRoute_InOutDirection_QTurn(t *testing.T) {
	sr := SubRoute{
		Type:      QTurn,
		MoveType:  10,
		Start:     Point{199050, 22600},
		End:       Point{200350, 24000},
		RefParams: [6]int32{400, 700, 1300, 0, 0, 0},
		RefPoints: [2]Point{{200350, 22600}, {0, 0}},
	}

	in, out := sr.InOutDirection()
	if in != XInc || out != YInc {
		t.Errorf("want in XInc, out YInc, result in %s, out %s", in, out)
	}
}

func TestSubRoute_InOutDirection_Oblique(t *testing.T) {
	sr := SubRoute{
		Type:      Oblique,
		MoveType:  26,
		Start:     Point{157450.000000, 22600.000000},
		End:       Point{154950.000000, 21400.000000},
		RefParams: [6]int32{1200, 0, 0, 0, 0, 0},
	}

	in, out := sr.InOutDirection()
	if in != XDec || out != YDec {
		t.Errorf("want in XDec, out YDec, result in %s, out %s", in, out)
	}
}

func TestSubRoute_String(t *testing.T) {
	sr := SubRoute{
		Type:      Oblique,
		MoveType:  26,
		Start:     Point{157450.000000, 22600.000000},
		End:       Point{154950.000000, 21400.000000},
		RefParams: [6]int32{1200, 0, 0, 0, 0, 0},
	}
	log.Print(sr)
}
