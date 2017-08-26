/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^'''''''''''''
 * vgen - visualization generator
 * ==============================
 *		This program renders an svg image of a nmir json file
 *
 *  Copyright The Deter Project 2017 - All Rights Reserved
 *  Apache 2.0 License
 *=================-----------------------------------------................```*/

package main

import (
	"github.com/rcgoodfellow/nmir"
	"github.com/rcgoodfellow/nmir/viz"
	"log"
	"os"
	"path"
	"strings"
)

func main() {

	log.SetFlags(0)

	if len(os.Args) < 2 {
		usage()
	}

	filename := os.Args[1]

	net, err := nmir.FromFile(filename)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = viz.NetSvg(
		path.Base(strings.Replace(filename, ".json", "", 1)), net)
	if err != nil {
		log.Fatalf("%v", err)
	}

}

func usage() {
	log.Fatal("usage: vgen nmir.json")
}
