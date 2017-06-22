package nmir

import (
	"log"
	"math"
	//"sort"
)

var lps LayoutParameters

func init() {

	lps.K = 10
	lps.R = 100

}

func VTag(net Net) Net {

	net = initialSpread(net)
	for i := 0; i < 20; i++ {
		log.Println(">>>---------->")
		layout(&net)
	}
	return net

}

func initialSpread(net Net) Net {

	radius := 100.0
	increment := 2 * math.Pi / float64(len(net.Nodes))

	for i, n := range net.Nodes {

		angle := float64(i) * increment
		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		n.Props["position"] = &Vec2{x, y}

	}

	return net

}

func layout(net *Net) {

	contract(net)
	expand(net)
	labelNodes(net)

}

func contract(net *Net) {

	for _, l := range net.Links {

		a, b := endpointSets(net, *l)

		for _, p := range a {
			for _, q := range b {
				theta := angle(p, q)
				p_pos := p.Props["position"].(*Vec2)
				q_pos := q.Props["position"].(*Vec2)

				p_pos.X -= lps.K * math.Cos(theta)
				p_pos.Y -= lps.K * math.Sin(theta)

				theta -= math.Pi

				q_pos.X -= lps.K * math.Cos(theta)
				q_pos.Y -= lps.K * math.Sin(theta)
			}
		}

	}

}

func expand(net *Net) {

	for _, a := range net.Nodes {
		for _, b := range net.Nodes {

			if a == b {
				continue
			}

			theta := angle(a, b)
			dist := distance(a, b)
			a_pos := a.Props["position"].(*Vec2)
			b_pos := b.Props["position"].(*Vec2)
			r := lps.R / dist

			a_pos.X += r * math.Cos(theta)
			a_pos.Y += r * math.Sin(theta)

			theta -= math.Pi

			b_pos.X += r * math.Cos(theta)
			b_pos.Y += r * math.Sin(theta)

		}
	}

}

func labelNodes(net *Net) {

	/*
		//calculate the angles of the outgoing links for each node and place a label
		//in the most vacant (wrt. lines) spot outside the node

		//calculating the angles
		for _, n := range net.Nodes {
			var angles sort.Float64Slice
			for _, e := range n.Endpoints {
				for _, nbr := range e.Neighbors {

					angles = append(angles, angle(n, nbr.Endpoint.Parent))

				}
			}

			//determine the most open position
			pos := 0
			label_angle := 0.0
			if len(angles) > 0 {
				angles.Sort()

				sep := 0.0
				for i, _ := range angles[:len(angles)-1] {

					delta := angles[i] - angles[i+1]
					log.Printf("[%s] %d-%d %f %f", n.Props["name"], i, i+1, angles[i], delta)
					if math.Abs(delta) > sep {
						sep = math.Abs(delta)
						pos = i
						label_angle = angles[pos] + delta/2.0
					}
				}

				delta := angles[len(angles)-1] - (2*math.Pi + angles[0])
				log.Printf("[%s] %d-%d %f %f",
					n.Props["name"], len(angles), 0, angles[len(angles)-1], delta)
				if math.Abs(delta) > sep {
					sep = math.Abs(delta)
					pos = len(angles) - 1
					label_angle = angles[pos] + delta/2.0
				}

			}
			//log.Printf("{%s} %f", n.Props["name"], label_angle)
			n.Props["label_angle"] = label_angle
		}
	*/

	for _, n := range net.Nodes {
		n.Props["label_angle"] = math.Pi / 4.0
	}

}

func angle(a, b *Node) float64 {

	a_pos := a.Props["position"].(*Vec2)
	b_pos := b.Props["position"].(*Vec2)

	dx := a_pos.X - b_pos.X
	dy := a_pos.Y - b_pos.Y

	theta := math.Atan2(dy, dx)
	if theta < 0 {
		theta += 2 * math.Pi
	}

	return theta

}

func distance(a, b *Node) float64 {

	a_pos := a.Props["position"].(*Vec2)
	b_pos := b.Props["position"].(*Vec2)

	dx := a_pos.X - b_pos.X
	dy := a_pos.Y - b_pos.Y

	return math.Sqrt(dx*dx + dy*dy)

}

func endpointSets(net *Net, l Link) ([]*Node, []*Node) {

	a := []*Node{}
	b := []*Node{}

	for _, e := range l.Endpoints[0] {
		a = append(a, net.GetNode(e.Id))
	}
	for _, e := range l.Endpoints[1] {
		b = append(b, net.GetNode(e.Id))
	}

	return a, b

}

// types ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type LayoutParameters struct {
	K, R float64
}
