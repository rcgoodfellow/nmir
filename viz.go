/*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 * Deter 2.0 nmir
 * ==================
 *   A bit of glue to integrate the viz submodule.
 *
 *	Copyright the Deter Project 2017 - All Rights Reserved
 *	License: Apache 2.0
 *~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

package nmir

type Point struct {
	X, Y float64
}

// Node Positional interface implementation ~~~

func (n *Node) Position() *Point {
	return n.Props["position"].(*Point)
}
func (n *Node) Velocity() *Point {
	return n.Props["dp"].(*Point)
}
func (n *Node) Weight() float64 {
	v := float64(n.Valence())
	if v < 15.0 {
		v = 15.0
	}
	return 10 * v
}

// Net Positional interface implementation ~~~

func (n *Net) Position() *Point {
	return n.Props["position"].(*Point)
}
func (n *Net) Velocity() *Point {
	return n.Props["dp"].(*Point)
}

// Endpoint Positional interface implementation ~~~

func (n *Endpoint) Position() *Point {
	return n.Props["position"].(*Point)
}
func (n *Endpoint) Velocity() *Point {
	return n.Props["dp"].(*Point)
}

// local-global point translations

func (n *Net) Global() Point {

	p := *n.Position()

	if n.Parent == nil {
		return p
	}

	return Add(p, n.Parent.Global())

}

func (n *Node) Global() Point {

	p := *n.Position()

	return Add(p, n.Parent.Global())

}

func (n *Endpoint) Global() Point {

	p := *n.Position()

	return Add(p, n.Parent.Global())

}

func Add(v, x Point) Point {

	return Point{v.X + x.X, v.Y + x.Y}

}
