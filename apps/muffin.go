package main

import (
	"fmt"
	"github.com/rcgoodfellow/nmir"
)

func main() {

	a := nmir.NewNet()
	s0 := a.Node().Set(nmir.Props{"name": "s0"})
	s1 := a.Node().Set(nmir.Props{"name": "s1"})
	for i := 0; i < 5; i++ {
		n := a.Node().Set(nmir.Props{"name": fmt.Sprintf("n%d", i)})
		a.Link(
			[]*nmir.Endpoint{s0.Endpoint()},
			[]*nmir.Endpoint{n.Endpoint()},
		)
		a.Link(
			[]*nmir.Endpoint{s1.Endpoint()},
			[]*nmir.Endpoint{n.Endpoint()},
		)
	}

	nmir.NetSvg("muffin", a)

}
