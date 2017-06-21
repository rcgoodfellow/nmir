package nmir

import (
	"math"
)

func VTag(net Net) Net {

	net = initialSpread(net)
	return net

}

func initialSpread(net Net) Net {

	radius := 100.0
	increment := 2 * math.Pi / float64(len(net.Nodes))

	for i, n := range net.Nodes {

		angle := float64(i) * increment
		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		n.Props["position"] = Position{x, y}

	}

	return net

}

// types ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
