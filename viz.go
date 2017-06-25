package nmir

import (
	//	"fmt"
	"github.com/ajstarks/svgo"
	"log"
	"math"
	"os"
)

const (
	AK    = 1
	RK    = 1
	Iters = 1000
)

var Step = 1.0
var Max = 0.0
var CForce = 0
var RForce = 0

func VTag(net *Net) {

	initialSpread(net, &Vec2{0, 0})
	for i := 0; i < Iters; i++ {
		if !layout(net) {
			break
		}
		log.Printf(">>>- %d ------------->", i)
	}

}

type Bounds struct {
	Width, Height float64
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
		if math.Abs(p.X+cx) > b.Width {
			b.Width = math.Abs(p.X + cx)
		}
		if math.Abs(p.Y+cy) > b.Height {
			b.Height = math.Abs(p.Y + cy)
		}

	}
	b.Width *= 2
	b.Height *= 2

	b.Width += 50
	b.Height += 50

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

	cx, cy := int(b.Width/2.0), int(b.Height/2.0)

	canvas := svg.New(out)
	canvas.Start(int(b.Width), int(b.Height))

	defer canvas.End()

	doNetSvgLink(net, cx, cy, canvas)
	doNetSvgNode(net, cx, cy, canvas)

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

		color := "fill:#334f7c"
		if n.Valence() > 10 {
			color = "fill:#11a052"
		}

		pos = n.Props["position"].(*Vec2)
		canvas.Circle(
			cx+int(pos.X),
			cy+int(pos.Y),
			10,
			color,
		)
		canvas.Text(
			cx+int(pos.X),
			cy+int(pos.Y),
			//fmt.Sprintf("%d", n.Valence()),
			n.Props["name"].(string),
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
		n.Props["dp"] = &Vec2{0, 0}
		log.Printf("{%s} tagged", n.Props["name"])

		for _, e := range n.Endpoints {

			e.Props["position"] = &Vec2{0, 0}

		}

	}

}

func layout(net *Net) bool {

	//contract(net)
	//expand(net)
	//force(net)
	force(net)
	//adapt(net)
	//step(net)
	constrain(net)
	/*
		if adapt(net) {
			reset(net)
			return true
		}
		if Max < 1e-4 {
			return false
		}
	*/
	step(net)
	return true
	//labelNodes(net)

}

func adapt(net *Net) bool {
	result := false
	Max = 0
	for _, n := range net.Nodes {
		dp := n.Props["dp"].(*Vec2)
		x := math.Sqrt(dp.X*dp.X + dp.Y*dp.Y)
		if x > Max {
			Max = x
		}
	}
	//if Max < 1 {
	//Step *= 10
	//}
	if Max > 1000 {
		Step *= 0.1
		result = true
	}

	log.Printf("A[%f](%f)", Step, Max)
	return result
}

func step(net *Net) {
	for _, n := range net.Nodes {
		p := n.Props["position"].(*Vec2)
		dp := n.Props["dp"].(*Vec2)

		p.X += dp.X
		p.Y += dp.Y

		dp.X = 0
		dp.Y = 0
	}
}

func reset(net *Net) {
	for _, n := range net.Nodes {
		dp := n.Props["dp"].(*Vec2)
		dp.X = 0
		dp.Y = 0
	}
}

func force(net *Net) {

	for _, a := range net.Nodes {
		for _, b := range net.Nodes {
			if a == b {
				continue
			}
			fab(a, b)
		}
	}

}

func constrain(net *Net) {

	for _, a := range net.Nodes {
		for _, b := range a.Neighbors() {
			gab(a, b)
		}
	}

}

func fab(a, b *Node) {
	p := b.Props["dp"].(*Vec2)

	d := node_distance(a, b)
	if d == 0 {
		return
	}
	repulsive := (1 / (d)) * Step
	f := repulsive

	theta := node_angle(a, b)
	p.X += f * math.Cos(theta)
	p.Y += f * math.Sin(theta)

}

func qfab(a *Qnode, b *Node) {
	p := b.Props["dp"].(*Vec2)

	d := node_distance(a.Data.(*Node), b)
	if d == 0 {
		return
	}
	repulsive := (1 / (d)) * Step
	f := repulsive

	theta := node_angle(a.Data.(*Node), b)
	p.X += f * math.Cos(theta)
	p.Y += f * math.Sin(theta)

}

func gab(a, b *Node) {
	av := float64(a.Valence())
	p := b.Props["dp"].(*Vec2)

	d := node_distance(a, b)
	if d < av {
		return
	}
	attractive := 1 * Step
	f := attractive

	theta := node_angle(b, a)
	p.X += f * math.Cos(theta)
	p.Y += f * math.Sin(theta)

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
