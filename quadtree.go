package nmir

type Point struct {
	X, Y float64
}

type Positional interface {
	Position() *Point
	Velocity() *Point
	Weight() float64
}

type Quad [4]*Qnode

type Qnode struct {
	Quad     Quad
	Data     Positional
	Centroid Point
	Weight   float64
}

func (q *Qnode) Insert(x *Qnode) {
	if x == nil {
		return
	}
	q.Weight += x.Data.Weight()

	i := q.Select(x.Data)
	p := q.Quad[i]

	if p != nil {
		p.Insert(x)
		return
	} else {
		q.Quad[i] = x
	}
}

func (q *Qnode) Select(x Positional) int {

	if x.Position().Y <= q.Data.Position().Y {
		if x.Position().X <= q.Data.Position().X {
			return 0
		} else {
			return 1
		}
	} else {
		if x.Position().X > q.Data.Position().X {
			return 2
		} else {
			return 3
		}
	}

}

// --------- Sorting

type Xnodes struct {
	N []*Node
}
type Ynodes struct {
	N []*Node
}

func (x Xnodes) Len() int      { return len(x.N) }
func (x Xnodes) Swap(i, j int) { x.N[i], x.N[j] = x.N[j], x.N[i] }
func (x Xnodes) Less(i, j int) bool {
	return x.N[i].Position().X < x.N[j].Position().X
}

func (x Ynodes) Len() int      { return len(x.N) }
func (x Ynodes) Swap(i, j int) { x.N[i], x.N[j] = x.N[j], x.N[i] }
func (x Ynodes) Less(i, j int) bool {
	return x.N[i].Position().Y < x.N[j].Position().Y
}

// --------- Connections to nmir

func (n *Node) Position() *Point {
	return n.Props["position"].(*Point)
}
func (n *Node) Velocity() *Point {
	return n.Props["dp"].(*Point)
}

func (n *Net) Position() *Point {
	return n.Props["position"].(*Point)
}
func (n *Net) Velocity() *Point {
	return n.Props["dp"].(*Point)
}

func (n *Endpoint) Position() *Point {
	return n.Props["position"].(*Point)
}
func (n *Endpoint) Velocity() *Point {
	return n.Props["dp"].(*Point)
}

func (n *Node) Weight() float64 {
	return 5 * float64(n.Valence())
	//return 1.0
}
