package main

import (
	"fmt"
	"github.com/rcgoodfellow/nmir"
	"github.com/rcgoodfellow/nmir/viz"
)

func main() {

	a := nmir.NewNet()
	s0 := a.Node().Set(nmir.Props{"name": "s0"})
	s1 := a.Node().Set(nmir.Props{"name": "s1"})
	for i := 0; i < 5; i++ {
		n := a.Node().Set(nmir.Props{"name": fmt.Sprintf("n%d", i)})
		a.Link(s0.Endpoint(), n.Endpoint())
		a.Link(s1.Endpoint(), n.Endpoint())
	}

	viz.NetSvg("muffin", a)

}
