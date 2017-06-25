package nmir

import (
	//"github.com/ajstarks/svgo"
	"log"
	"math"
	"sort"
	//"os"
)

const (
	InitRadius = 100.0
	Chunk      = 10
	AggD       = 100
)

func QTree(net *Net) *Qnode {

	n := float64(len(net.Nodes))

	for i, x := range net.Nodes {
		theta := 2 * math.Pi / n * float64(i)

		x.Props["position"] = &Point{
			InitRadius * math.Cos(theta),
			InitRadius * math.Sin(theta),
		}
	}

	root := buildQ(net.Nodes, true)

	return root

}

func buildQ(nodes []*Node, odd bool) *Qnode {

	if len(nodes) == 0 {
		return nil
	}

	if odd {
		sort.Sort(Xnodes{nodes[:]})
	} else {
		sort.Sort(Ynodes{nodes[:]})
	}
	v := len(nodes) / 2

	root := &Qnode{Data: nodes[v]}

	root.Insert(buildQ(nodes[:v], !odd))

	if v < len(nodes)-1 {
		root.Insert(buildQ(nodes[v+1:], !odd))
	}

	return root

}

func QBounds(qtr *Qnode) Bounds {
	ul := QUpperLeft(qtr)
	lr := QLowerRight(qtr)
	return Bounds{
		lr.Y - ul.Y,
		lr.X - ul.X,
	}
}

func QUpperLeft(qtr *Qnode) *Point {
	return QBound(qtr, 0)
}

func QLowerRight(qtr *Qnode) *Point {
	return QBound(qtr, 3)
}

func QBound(qtr *Qnode, quad int) *Point {

	x := qtr.Quad[quad]
	if x != nil {
		return QBound(x, quad)
	}
	return x.Data.Position()

}

func QHeight(qtr *Qnode) int {

	h := 1
	hs := make([]int, 4)
	for i, x := range qtr.Quad {
		if x != nil {
			hs[i] = QHeight(x)
		}
	}
	max := 0
	for _, x := range hs {
		if x > max {
			max = x
		}
	}
	return h + max

}

func QSize(qtr *Qnode) int {

	s := 1
	for _, x := range qtr.Quad {
		if x != nil {
			s += QSize(x)
		}
	}
	return s
}

func QWeight(qtr *Qnode) float64 {

	qtr.Weight = qtr.Data.Weight()
	for _, x := range qtr.Quad {
		if x != nil {
			qtr.Weight += QWeight(x)
		}
	}
	return qtr.Weight
}

func QSpread(net *Net, qtr *Qnode) {

	for _, n := range net.Nodes {

		Qforce(n, qtr)

	}

}

func Qforce(node *Node, qtr *Qnode) {

	d := node_distance(node, qtr.Data.(*Node))
	if d >= AggD {
		qfab(qtr, node)
	}

}

func qlayout(net *Net) {

	//log.Printf("building qtree")
	qtr := QTree(net)

	log.Printf("qtree: size=%d height=%d weight=%f root=%s",
		QSize(qtr),
		QHeight(qtr),
		QWeight(qtr),
		qtr.Data.(*Node).Props["name"],
	)

	//log.Printf("spreading")
	QSpread(net, qtr)

}

func QNetSvg(name string, net *Net) error {
	log.Printf("net: size=%d", len(net.Nodes))

	for i := 0; i < 20; i++ {

		qlayout(net)

	}

	/*
		out, err := os.Create(name + ".svg")
		if err != nil {
			return err
		}
		defer out.Close()

		b := QBounds(qtr)
		cx, cy := int(b.Width/2.0), int(b.Height/2.0)

		canvas := svg.New(out)
		canvas.Start(int(b.Width), int(b.Height))

		defer canvas.End()

		doNetSvgLink(net, cx, cy, canvas)
		doNetSvgNode(net, cx, cy, canvas)
	*/

	return nil

}
