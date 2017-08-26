package main

import (
	"fmt"
	"github.com/deter-project/testbed/core/nmir/models"
	"github.com/deter-project/testbed/core/nmir/viz"
	"os"
)

func main() {

	tb := models.CEF_3bed()
	tb.ToFile("cef-3host.json")

	err := viz.NetSvg("cef-3host", tb)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}
