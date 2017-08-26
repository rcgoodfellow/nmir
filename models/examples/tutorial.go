package main

import (
	"fmt"
	"github.com/deter-project/testbed/core/nmir/models"
	"github.com/deter-project/testbed/core/nmir/viz"
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
