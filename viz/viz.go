/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 * Deter 2.0 nmir-viz
 * ==================
 *   This file provides visualization functionality for nmir networks. The
 *   following features are provided at this time:
 *
 *			- force layout computation
 *			- svg image generation
 *
 *	Copyright the Deter Project 2017 - All Rights Reserved
 *	License: Apache 2.0
 *~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

package viz

import (
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/rcgoodfellow/nmir"
	"math"
	"os"
)

const (
	// number of iterations to use for force layout
	Iters = 100
)

// The positional interface must be implemented by data that is laied out and
// presented by this code
type Positional interface {
	Position() *nmir.Point
	Velocity() *nmir.Point
	Weight() float64
}

// Simple data structure to keep track of the size of things
type Bounds struct {
	Width, Height float64
}

// Tag an nmir network with position and velocity values in pereparation for a
// layout computation
func Vtag(net *nmir.Net) {
	n := float64(len(net.Nodes))
	net.Props["position"] = &nmir.Point{0, 0}

	for i, x := range net.Nodes {
		theta := 2 * math.Pi / n * float64(i)

		InitRadius := 2.0 * float64(len(net.Nodes))

		x.Props["position"] = &nmir.Point{
			InitRadius * math.Cos(theta),
			InitRadius * math.Sin(theta),
		}
		x.Props["dp"] = &nmir.Point{0, 0}

		for _, e := range x.Endpoints {
			e.Props["position"] = &nmir.Point{0, 0}
		}
	}
}

// Compute the bounds of a network relative to the central point (cx, cy)
func bounds(net *nmir.Net, b *Bounds, cx, cy float64) {

	pos := net.Props["position"].(*nmir.Point)
	cx += pos.X
	cy += pos.Y

	for _, n := range net.Nets {
		bounds(n, b, cx, cy)
	}

	for _, n := range net.Nodes {

		p := n.Props["position"].(*nmir.Point)
		if math.Abs(p.X+cx) > b.Width {
			b.Width = math.Abs(p.X + cx)
		}
		if math.Abs(p.Y+cy) > b.Height {
			b.Height = math.Abs(p.Y + cy)
		}

	}
	b.Width *= 2
	b.Height *= 2

	b.Width += 100
	b.Height += 100

}

func calcLabelPos(n *nmir.Node, lbl string) *nmir.Point {

	p := n.Position()
	dp := n.Velocity()
	arg := math.Atan2(dp.Y, dp.X)
	if arg < 0 {
		arg = 2*math.Pi + arg
	}

	flip := 0.0
	if arg > math.Pi/2 && arg < 3*math.Pi/2 {
		flip = -(5.0 * float64(len(lbl)))
	}

	return &nmir.Point{p.X + math.Cos(arg)*10 + flip, p.Y + math.Sin(arg)*10}

}

// Draw the nodes of an nmir network on the provided svg canvas relative
// to the central point (cx, cy)
func SvgDrawNodes(net *nmir.Net, cx, cy int, canvas *svg.SVG) {
	pos := net.Position()
	cx += int(pos.X)
	cy += int(pos.Y)

	for _, n := range net.Nets {
		SvgDrawNodes(n, cx, cy, canvas)
	}

	for _, n := range net.Nodes {

		color := "fill:#9eb4d8"
		if n.Valence() > 10 {
			color = "fill:#11a052"
		}

		//FIXME remove debug reference node
		if n.Props["name"] == "bpc223" {
			color = "fill:#ff0000;font-size:11px;"
		}

		pos = n.Position()
		canvas.Group(
			fmt.Sprintf("id='%s'", n.Id),
			"class='nmir-node'",
		)
		canvas.Circle(
			cx+int(pos.X),
			cy+int(pos.Y),
			10,
			color,
		)
		lbl := n.Props["name"].(string)
		pos = calcLabelPos(n, lbl)
		canvas.Text(
			cx+int(pos.X),
			cy+int(pos.Y),
			n.Props["name"].(string),
		)
		canvas.Gend()

	}

}

// Draw the links of an nmir network on the provided svg canvas relative
// to the central point (cx, cy)
func SvgDrawLinks(net *nmir.Net, cx, cy int, canvas *svg.SVG) {

	pos := net.Position()
	cx += int(pos.X)
	cy += int(pos.Y)

	for _, n := range net.Nets {
		SvgDrawLinks(n, cx, cy, canvas)
	}

	for _, l := range net.Links {

		for _, p_ := range l.Endpoints {
			for _, q_ := range l.Endpoints {

				canvas.Group(
					fmt.Sprintf("id='%s'", l.Id),
					"class='nmir-link'",
				)
				if l.IsLocal() {
					var p, q *nmir.Point
					p = p_.Parent.Position()
					q = q_.Parent.Position()

					canvas.Line(
						cx+int(p.X), cy+int(p.Y),
						cx+int(q.X), cy+int(q.Y),
						"stroke:#a4a9b2; stroke-width: 2;",
					)
				} else {
					p := p_.Global()
					q := q_.Global()
					canvas.Line(
						int(p.X), int(p.Y),
						int(q.X), int(q.Y),
						"stroke:#182a47; stroke-width: 2;",
					)
				}
				canvas.Gend()

			}
		}

	}

}

// Iterate through all the nodes and apply the current velocities to
// their current positions
func step(net *nmir.Net) {
	for _, n := range net.Nodes {
		p := n.Props["position"].(*nmir.Point)
		dp := n.Props["dp"].(*nmir.Point)

		p.X += dp.X
		p.Y += dp.Y

		dp.X = 0
		dp.Y = 0
	}
}

// Iterate through all the links and compute the force imposed by the link
// on the connected nodes
func constrain(net *nmir.Net) {

	for _, l := range net.Links {
		for _, ea := range l.Endpoints {
			for _, eb := range l.Endpoints {

				gab(ea.Parent, eb.Parent)

			}
		}
	}

}

// Compute the force of a-on-b due to node repulsion
func fab(a, b Positional) {

	p := b.Velocity()

	d := distance(a, b)
	if d == 0 {
		return
	}
	m := a.Weight()
	repulsive := (m / (d))
	f := repulsive

	theta := angle(a, b)
	p.X += f * math.Cos(theta)
	p.Y += f * math.Sin(theta)

}

// Compute the mutual force of a-on-b and b-on-a due to link constriction
func gab(a, b *nmir.Node) {

	av := float64(a.Valence())
	bv := float64(b.Valence())
	da := a.Props["dp"].(*nmir.Point)
	db := b.Props["dp"].(*nmir.Point)

	d := distance(a, b)
	if d < 10.0 {
		return
	}

	theta := angle(b, a)
	db.X += (d / 10.0 / bv) * math.Cos(theta)
	db.Y += (d / 10.0 / bv) * math.Sin(theta)

	theta = angle(a, b)
	da.X += (d / 10.0 / av) * math.Cos(theta)
	da.Y += (d / 10.0 / av) * math.Sin(theta)

}

// Calculate the angle between two positional elements
func angle(a, b Positional) float64 {

	dx := b.Position().X - a.Position().X
	dy := b.Position().Y - a.Position().Y

	theta := math.Atan2(dy, dx)
	if theta < 0 {
		theta += 2 * math.Pi
	}

	return theta

}

// Calculate the distance between two positional elements
func distance(a, b Positional) float64 {

	dx := a.Position().X - b.Position().X
	dy := a.Position().Y - b.Position().Y

	return math.Sqrt(dx*dx + dy*dy)

}

// Create an SVG image
func NetSvg(name string, net *nmir.Net) error {
	//log.Printf("net: size=%d", len(net.Nodes))

	// prepate the network
	Vtag(net)

	// run the layout engine
	for i := 0; i < Iters; i++ {

		//log.Printf(">>==={%d}-------->", i)
		layout(net)

	}

	// open a file to write the svg out to
	out, err := os.Create(name + ".svg")
	if err != nil {
		return err
	}
	defer out.Close()

	// determine the size of the svg canvas
	var b Bounds
	bounds(net, &b, 0, 0)
	cx, cy := int(b.Width/2.0), int(b.Height/2.0)

	// create the svg
	canvas := svg.New(out)
	canvas.Start(int(b.Width), int(b.Height))

	defer canvas.End()

	SvgDrawLinks(net, cx, cy, canvas)
	SvgDrawNodes(net, cx, cy, canvas)

	js, err := net.Json()
	if err != nil {
		return err
	}

	canvas.Script("text/json", string(js))

	return nil

}

// Compute the layout of a nmir network, returning the root of the quadtree
// that was used to calcuate it.
func layout(net *nmir.Net) *Pinode {

	ptr := PTree(net)

	ptr.Forces(net)
	step(net)
	constrain(net)
	step(net)
	ptr.Forces(net)

	return ptr

}
