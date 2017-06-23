package nmir

import (
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"os"
	//"sort"
)

var lps LayoutParameters

func init() {

	lps.K = 1
	lps.R = 1
	lps.Step = 1e-3

}

func VTag(net *Net) {

	initialSpread(net, &Vec2{0, 0})
	for i := 0; i < 100; i++ {
		layout(net)
		log.Println(">>>---------->")
	}

}

type Bounds struct {
	Xmin, Xmax, Ymin, Ymax float64
}

func bounds(net *Net, b *Bounds, cx, cy float64) {

	pos := net.Props["position"].(*Vec2)
	cx += pos.X
	cy += pos.Y

	for _, n := range net.Nets {
		bounds(n, b, cx, cy)
	}

	for _, n := range net.Nodes {

		p := n.Props["position"].(*Vec2)
		if p.X+cx > b.Xmax {
			b.Xmax = p.X + cx
		}
		if p.X+cx < b.Xmin {
			b.Xmin = p.X + cx
		}
		if p.Y+cy > b.Ymax {
			b.Ymax = p.Y + cy
		}
		if p.Y+cy < b.Ymin {
			b.Ymin = p.Y + cy
		}

	}

}

func NetSvg(name string, net *Net) error {

	VTag(net)
	out, err := os.Create(name + ".svg")
	if err != nil {
		return err
	}
	defer out.Close()

	var b Bounds
	bounds(net, &b, 0.0, 0.0)

	width := int(b.Xmax-b.Xmin) + 100
	height := int(b.Ymax-b.Ymin) + 100
	cx, cy := width/2, height/2
	p := net.Props["position"].(*Vec2)
	p.X = float64(cx)
	p.Y = float64(cy)

	canvas := svg.New(out)
	canvas.Start(width, height)

	defer canvas.End()

	doNetSvgLink(net, 0, 0, canvas)
	doNetSvgNode(net, 0, 0, canvas)

	return nil

}

func (n *Net) Global() Vec2 {

	p := n.Props["position"].(*Vec2)

	if n.Parent == nil {
		return *p
	}

	return p.Add(n.Parent.Global())

}

func (n *Node) Global() Vec2 {

	p := n.Props["position"].(*Vec2)

	return p.Add(n.Parent.Global())

}

func (n *Endpoint) Global() Vec2 {

	p := n.Props["position"].(*Vec2)

	return p.Add(n.Parent.Global())

}

func doNetSvgNode(net *Net, cx, cy int, canvas *svg.SVG) {
	pos := net.Props["position"].(*Vec2)
	cx += int(pos.X)
	cy += int(pos.Y)

	for _, n := range net.Nets {
		doNetSvgNode(n, cx, cy, canvas)
	}

	for _, n := range net.Nodes {

		pos = n.Props["position"].(*Vec2)
		canvas.Circle(
			cx+int(pos.X),
			cy+int(pos.Y),
			10,
			"fill:#334f7c",
		)

	}

}

func doNetSvgLink(net *Net, cx, cy int, canvas *svg.SVG) {

	pos := net.Props["position"].(*Vec2)
	cx += int(pos.X)
	cy += int(pos.Y)

	for _, n := range net.Nets {
		doNetSvgLink(n, cx, cy, canvas)
	}

	for _, l := range net.Links {

		for _, p_ := range l.Endpoints[0] {
			for _, q_ := range l.Endpoints[1] {

				if l.IsLocal() {
					var p, q *Vec2
					p = p_.Parent.Props["position"].(*Vec2)
					q = q_.Parent.Props["position"].(*Vec2)

					canvas.Line(
						cx+int(p.X), cy+int(p.Y),
						cx+int(q.X), cy+int(q.Y),
						"stroke:#182a47",
					)
				} else {
					p := p_.Global()
					q := q_.Global()
					canvas.Line(
						int(p.X), int(p.Y),
						int(q.X), int(q.Y),
						"stroke:#182a47",
					)
				}

			}
		}

	}

}

func initialSpread(net *Net, centroid *Vec2) {

	radius := 10.0
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

		for _, e := range n.Endpoints {

			e.Props["position"] = &Vec2{0, 0}

		}

	}

}

func layout(net *Net) {

	//contract(net)
	//expand(net)
	force(net)
	labelNodes(net)

}

const (
	AK   = 0.25
	RK   = 0.75
	FK   = 0.0023
	Step = 1e-2
)

func force(net *Net) {

	for _, a := range net.Nodes {

		for _, b := range a.Neighbors() {

			fab(a, b)

		}

	}

}

func fab(a, b *Node) {
	av := float64(a.Valence())
	bv := float64(b.Valence())
	p := a.Props["position"].(*Vec2)

	d := node_distance(a, b)
	friction := av * FK
	attractive := -(bv * AK) * Step
	repulsive := (bv * av * RK / d * d) * Step
	f := (repulsive + attractive) / friction

	theta := node_angle(a, b)
	p.X += f * math.Cos(theta)
	p.Y += f * math.Sin(theta)

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

				//k := 1.0
				pw := 1.0
				qw := 1.0

				if l.IsLocal() {
					pw = float64(p_.Parent.Valence())
					qw = float64(p_.Parent.Valence())
					p = p_.Parent.Props["position"].(*Vec2)
					q = q_.Parent.Props["position"].(*Vec2)
					//d := distance(p, q)
					//k = float64((pw + qw)) * (1.0 / d * d)
				} else {
					p = p_.Parent.Parent.Props["position"].(*Vec2)
					q = q_.Parent.Parent.Props["position"].(*Vec2)
				}

				theta := angle(p, q)
				p.X += ((qw * math.Cos(theta)) / pw * 0.1) * lps.Step
				p.Y += ((qw * math.Sin(theta)) / pw * 0.1) * lps.Step

				//theta -= math.Pi
				theta = angle(q, p)
				q.X += ((pw * math.Cos(theta)) / qw * 0.1) * lps.Step
				q.Y += ((pw * math.Sin(theta)) / qw * 0.1) * lps.Step
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

			do_expand(a_pos, b_pos, 1, 1, 1, 1, lps.R*10)
		}
	}

	for _, n := range net.Nets {
		expand(n)
	}

	for _, x := range net.Nodes {

		o := x.Props["position"].(*Vec2)
		ow := float64(x.Valence())

		for _, e := range x.Endpoints {
			for _, n := range e.Neighbors {

				p := n.Endpoint.Parent.Props["position"].(*Vec2)
				pw := float64(n.Endpoint.Parent.Valence())

				d := distance(o, p)

				of := ow / d * d
				pf := pw / d * d

				//k := float64(ow+pw) * (1.0 / d * d)
				do_expand(o, p, of, pf, ow, pw, lps.R)

				for _, m := range e.Neighbors {
					if n == m {
						continue
					}

					q := m.Endpoint.Parent.Props["position"].(*Vec2)
					qw := float64(m.Endpoint.Parent.Valence())

					d := distance(q, p)
					qf := qw / d * d

					//k := float64(qw+pw) * (1.0 / d * d)
					do_expand(p, q, pf, qf, pw, qw, lps.R)
				}

			}

		}
	}

	log.Println("expand finish")

}

func do_expand(a, b *Vec2, af, bf, aw, bw, repel float64) {
	if a == b {
		return
	}
	theta := angle(b, a)
	//dist := distance(a, b)
	//r := repel // / dist

	a.X += (bf * aw * 0.1 * math.Cos(theta)) * lps.Step
	a.Y += (bf * aw * 0.1 * math.Sin(theta)) * lps.Step

	//theta -= math.Pi
	theta = angle(a, b)

	b.X += (af * bw * 0.1 * math.Cos(theta)) * lps.Step
	b.Y += (af * bw * 0.1 * math.Sin(theta)) * lps.Step
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

func node_angle(a, b *Node) float64 {
	return angle(
		a.Props["position"].(*Vec2),
		b.Props["position"].(*Vec2),
	)
}

func node_distance(a, b *Node) float64 {
	return distance(
		a.Props["position"].(*Vec2),
		b.Props["position"].(*Vec2),
	)
}

func angle(a, b *Vec2) float64 {

	dx := b.X - a.X
	dy := b.Y - a.Y

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

func (v *Vec2) Add(x Vec2) Vec2 {

	return Vec2{v.X + x.X, v.Y + x.Y}

}

// types ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

type Vec2 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type LayoutParameters struct {
	K, R, Step float64
}
