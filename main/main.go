package main

import (
	"log"
	"traffic/legume"
	"traffic/util"
)

func main() {
	l := legume.BaseLegume(10, util.XInc, util.YInc)
	log.Print(l)
}
