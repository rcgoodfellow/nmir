package main

import (
	"fmt"
	"github.com/rcgoodfellow/nmir/models"
	"github.com/rcgoodfellow/nmir/viz"
	"os"
)

func main() {

	tb := models.Tutorial_small()
	tb.ToFile("small.json")

	err := viz.NetSvg("small", tb)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}
