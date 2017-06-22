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

func VTag(net *Net) {

	initialSpread(net, &Vec2{0, 0})
	for i := 0; i < 10; i++ {
		layout(net)
		log.Println(">>>---------->")
	}

}

func initialSpread(net *Net, centroid *Vec2) {

	radius := 100.0
	phase := math.Pi / -8.0
	net.Props["position"] = centroid
	increment := 2 * math.Pi / float64(len(net.Nets))

	for i, n := range net.Nets {

		angle := float64(i)*increment + phase
		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		v := &Vec2{x, y}
		n.Props["position"] = &Vec2{x, y}

		initialSpread(n, v)

	}

	increment = 2 * math.Pi / float64(len(net.Nodes))
	for i, n := range net.Nodes {

		angle := float64(i) * increment
		x := radius * math.Cos(angle)
		y := radius * math.Sin(angle)
		n.Props["position"] = &Vec2{x, y}
		log.Printf("{%s} tagged", n.Props["name"])

	}

}

func layout(net *Net) {

	contract(net)
	expand(net)
	labelNodes(net)

}

func contract(net *Net) {

	log.Println("contract start")

	for _, n := range net.Nets {
		contract(n)
	}

	for _, l := range net.Links {

		for _, p_ := range l.Endpoints[0] {
			for _, q_ := range l.Endpoints[1] {
				var p, q *Vec2

				if l.IsLocal() {
					p = p_.Parent.Props["position"].(*Vec2)
					q = q_.Parent.Props["position"].(*Vec2)
				} else {
					p = p_.Parent.Parent.Props["position"].(*Vec2)
					q = q_.Parent.Parent.Props["position"].(*Vec2)
				}

				theta := angle(p, q)

				p.X -= lps.K * math.Cos(theta)
				p.Y -= lps.K * math.Sin(theta)

				theta -= math.Pi

				q.X -= lps.K * math.Cos(theta)
				q.Y -= lps.K * math.Sin(theta)
			}
		}

	}

	log.Println("contract finish")

}

func expand(net *Net) {

	log.Println("expand start")

	for _, a := range net.Nets {
		for _, b := range net.Nets {
			a_pos := a.Props["position"].(*Vec2)
			b_pos := b.Props["position"].(*Vec2)

			do_expand(a_pos, b_pos, lps.R*10)
		}
	}

	for _, n := range net.Nets {
		expand(n)
	}

	for _, x := range net.Nodes {

		o := x.Props["position"].(*Vec2)

		for _, e := range x.Endpoints {

			for _, n := range e.Neighbors {

				p := n.Endpoint.Parent.Props["position"].(*Vec2)
				do_expand(o, p, lps.R)

				for _, m := range e.Neighbors {
					if n == m {
						continue
					}
					q := m.Endpoint.Parent.Props["position"].(*Vec2)
					do_expand(p, q, lps.R)
				}

			}

		}
	}

	log.Println("expand finish")

}

func do_expand(a, b *Vec2, repel float64) {
	if a == b {
		return
	}
	theta := angle(a, b)
	dist := distance(a, b)
	r := repel / dist

	a.X += r * math.Cos(theta)
	a.Y += r * math.Sin(theta)

	theta -= math.Pi

	b.X += r * math.Cos(theta)
	b.Y += r * math.Sin(theta)
}

func labelNodes(net *Net) {

	for _, n := range net.Nets {
		labelNodes(n)
	}

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

func angle(a, b *Vec2) float64 {

	dx := a.X - b.X
	dy := a.Y - b.Y

	theta := math.Atan2(dy, dx)
	if theta < 0 {
		theta += 2 * math.Pi
	}

	return theta

}

func distance(a, b *Vec2) float64 {

	dx := a.X - b.X
	dy := a.Y - b.Y

	return math.Sqrt(dx*dx + dy*dy)

}

// types ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type LayoutParameters struct {
	K, R float64
}
