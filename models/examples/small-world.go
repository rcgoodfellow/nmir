package main

import (
	"fmt"
	"github.com/rcgoodfellow/nmir/models"
	"github.com/rcgoodfellow/nmir/viz"
	"os"
)

func main() {

	world := models.CEF_SmallWorld()
	world.ToFile("small-world.json")

	err := viz.NetSvg("small-world", world)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	fmt.Println(world.String())

}
